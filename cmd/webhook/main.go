package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/Gearbox-protocol/sdk-go/log"
	"github.com/Gearbox-protocol/sdk-go/utils"
	"github.com/joho/godotenv"
)

type cmdObj struct {
	cmd      []string
	dontFail bool
}

func getCmds(remoteDB string, gpointbotDB string, forkUrl string) []cmdObj {
	var cmds = []cmdObj{
		// {"bash", "-x", "/Users/harshjain/BACKUP/gearbox/third-eye/db_scripts/local_testing/local_test.sh", "139.177.179.137", "172.232.121.133", "harshjain"},
		{cmd: []string{"sudo systemctl stop gpointbot"}},
		{cmd: []string{"psql", gpointbotDB, "-c", "drop table last_snaps"}, dontFail: true},
		{cmd: []string{"psql", gpointbotDB, "-c", "drop table user_points"}, dontFail: true},
		{cmd: []string{"psql", gpointbotDB, "-c", "drop table events"}, dontFail: true},
		{cmd: []string{"sudo systemctl restart gpointbot"}},
		{cmd: []string{"sudo systemctl restart trading_price"}},
		{cmd: []string{"sudo systemctl restart gearbox-ws"}},
		{cmd: []string{"sudo systemctl stop anvil-third-eye"}},
		{cmd: []string{"sudo systemctl stop charts_server"}},
		{cmd: []string{"sudo systemctl stop gearbox-ws"}},
		{cmd: []string{"bash", "-x", "/home/debian/anvil-third-eye/db_scripts/local_testing/anvil_test.sh", remoteDB, "debian", forkUrl}},
		{cmd: []string{"sudo systemctl restart gearbox-ws"}},
		{cmd: []string{"sudo systemctl restart anvil-third-eye"}},
		{cmd: []string{"sudo systemctl restart charts_server"}},
	}
	return cmds
}

type Config struct {
	log.CommonEnvs
	Port        string `env:"PORT" default:"9090"`
	RemoteDB    string `env:"REMOTE_GEARBOX_DB" default:""`
	GpointbotDB string `env:"LOCAL_GPOINTBOT_DB" default:""`
	ChainId     int    `env:"CHAIN_ID" default:"7878"`
	Forkurl     string `env:"FORK_URL" default:""`
}

func getConfig() *Config {
	godotenv.Load(".env")
	cfg := &Config{}
	utils.ReadFromEnv(&cfg.CommonEnvs)
	utils.ReadFromEnv(cfg)
	return cfg
}

func runCmdOld(cmdStr []string) {
	cmd := exec.Command(cmdStr[0], cmdStr[1:]...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Fatal(err,
			"stdIn ", stdout.String(),
			"stdOut", stderr.String(),
			"for cmd", cmdStr, len(cmdStr),
		)
	}
}

func copyAndCapture(w io.Writer, r io.Reader) ([]byte, error) {
	var out []byte
	buf := make([]byte, 1024, 1024)
	for {
		n, err := r.Read(buf[:])
		if n > 0 {
			d := buf[:n]
			out = append(out, d...)
			_, err := w.Write(d)
			if err != nil {
				return out, err
			}
		}
		if err != nil {
			// Read returns io.EOF at the end of file, which is not an error for us
			if err == io.EOF {
				err = nil
			}
			return out, err
		}
	}
}
func runCmdNew(cmdStr []string) (string, string, error) {
	cmd := exec.Command(cmdStr[0], cmdStr[1:]...)

	var stdout, stderr []byte
	var errStdout, errStderr error
	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()
	err := cmd.Start()
	if err != nil {
		log.Fatalf("cmd.Start() failed with '%s'\n", err)
	}

	// cmd.Wait() should be called only after we finish reading
	// from stdoutIn and stderrIn.
	// wg ensures that we finish
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		stdout, errStdout = copyAndCapture(os.Stdout, stdoutIn)
		wg.Done()
	}()

	stderr, errStderr = copyAndCapture(os.Stderr, stderrIn)

	wg.Wait()

	err = cmd.Wait()
	if err != nil {
		return "", "", fmt.Errorf("cmd.Run() failed with %s", err)
	}
	if errStdout != nil || errStderr != nil {
		return "", "", fmt.Errorf("failed to capture stdout or stderr")
	}
	outStr, errStr := string(stdout), string(stderr)
	return outStr, errStr, nil
}
func (m *runCmdsObj) runCmds() {
	m.mu.Lock()
	defer m.mu.Unlock()
	log.AMQPMsg("Anvil Webhook received")
	for _, cmd := range getCmds(m.remoteGBDB, m.localGpointbotDB, m.forkUrl) {
		cmdStr := cmd.cmd
		if len(cmdStr) == 1 {
			cmdStr = strings.Split(cmdStr[0], " ")
		}
		// runCmdOld(cmdStr)
		stdout, stderr, err := runCmdNew(cmdStr)
		if err != nil {
			log.Info(stdout)
			log.Info(stderr)
			if cmd.dontFail {
				log.Warn("Ignoring error", err, "for cmd", cmdStr)
			} else {
				log.Fatal(err)
			}
		}
	}
}

type runCmdsObj struct {
	mu               sync.Mutex
	remoteGBDB       string
	localGpointbotDB string
	forkUrl          string
}

func (m *runCmdsObj) ServeHTTP(hw http.ResponseWriter, hr *http.Request) {
	if hr.Method == "POST" {
		go m.runCmds()
		fmt.Fprint(hw, "OK")
	} else {
		fmt.Fprint(hw, "Only POST allowed")
	}
}

func server() {
	cfg := getConfig()
	log.NewAMQPService(
		cfg.AMQPEnable,
		cfg.AMQPUrl,
		log.LoggingConfig{
			Exchange:     "TelegramBot",
			ChainId:      int64(cfg.ChainId),
			RiskEndpoint: cfg.RiskEndpoint,
			RiskSecret:   cfg.RiskSecret,
		},
		cfg.AppName,
	)
	//
	mux := http.NewServeMux()
	mux.Handle("/anvil_fork_reset", &runCmdsObj{remoteGBDB: cfg.RemoteDB, localGpointbotDB: cfg.GpointbotDB, forkUrl: cfg.Forkurl})
	mux.Handle("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "OK")
	}))

	log.AMQPMsg("Anvil Webhook started")
	srv := http.Server{
		Addr:    utils.GetHost(cfg.Port),
		Handler: mux,
	}
	srv.ListenAndServe()
}
func main() {
	server()
}
