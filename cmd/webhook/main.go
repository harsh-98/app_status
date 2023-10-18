package main

import (
	"fmt"
	"net/http"
	"os/exec"
	"sync"

	"github.com/Gearbox-protocol/sdk-go/log"
	"github.com/Gearbox-protocol/sdk-go/utils"
	"github.com/joho/godotenv"
)

var cmds = []string{
	"sudo systemctl stop gpointbot",
	"cd /home/debian/gpointbot; sqlite3 local.db  'drop table last_snaps ; drop table user_points;'",
	"sudo systemctl restart gpointbot",
	"sudo systemctl restart trading_price",
	"sudo systemctl restart gearbox-ws",
	"sudo systemctl stop third-eye",
	"sudo systemctl stop charts_server",
	"cd /home/debian/third-eye; bash -x ./db_scripts/local_testing/local_test.sh '139.177.179.137' '' debian",
	"sudo systemctl restart third-eye",
	"sudo systemctl restart charts_server",
}

type Config struct {
	log.CommonEnvs
	Port int64 `env:"PORT" default:"9090"`
}

func getConfig() *Config {
	godotenv.Load(".env")
	cfg := &Config{}
	utils.ReadFromEnv(&cfg.CommonEnvs)
	utils.ReadFromEnv(cfg)
	return cfg
}

func (m *runCmdsObj) runCmds() {
	m.mu.Lock()
	defer m.mu.Unlock()
	log.AMQPMsg("Anvil Webhook received")
	for _, cmdStr := range cmds {
		cmd := exec.Command(cmdStr)
		err := cmd.Run()
		log.CheckFatal(err)
	}
}

type runCmdsObj struct {
	mu sync.Mutex
}

func (m *runCmdsObj) ServeHTTP(hw http.ResponseWriter, hr *http.Request) {
	go m.runCmds()
	fmt.Fprint(hw, "OK")
}

func server() {
	cfg := getConfig()
	log.NewAMQPService(
		cfg.AMQPEnable,
		cfg.AMQPUrl,
		log.LoggingConfig{
			Exchange:     "TelegramBot",
			ChainId:      7878,
			RiskEndpoint: cfg.RiskEndpoint,
			RiskSecret:   cfg.RiskSecret,
		},
		cfg.AppName,
	)
	//
	mux := http.NewServeMux()
	mux.Handle("/anvil_fork_reset", &runCmdsObj{})
	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: mux,
	}
	mux.Handle("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "OK")
	}))

	log.AMQPMsg("Anvil Webhook started")
	srv.ListenAndServe()
}
func main() {
	server()
}
