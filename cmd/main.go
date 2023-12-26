package main

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Gearbox-protocol/sdk-go/core"
	"github.com/Gearbox-protocol/sdk-go/log"
	"github.com/Gearbox-protocol/sdk-go/utils"
	"github.com/joho/godotenv"
)

func startLogging() {
	godotenv.Load(".env") // if file is not found, ignore
	cfg := log.CommonEnvs{}

	utils.ReadFromEnv(&cfg)

	log.NewAMQPService(
		cfg.AMQPEnable,
		cfg.AMQPUrl,
		log.LoggingConfig{
			Exchange: "TelegramBot",
			ChainId:  7878,
		},
		cfg.AppName,
	)

}

type StatusConfig map[string]ApplicationsUrl

type ApplicationsUrl map[string][]string

func main() {
	startLogging()
	port, err := strconv.ParseInt(utils.GetEnvOrDefault("PORT", "8080"), 10, 64)
	log.CheckFatal(err)
	//
	mgr := newStatusManager(port)
	mgr.Start()
	log.AMQPMsg("App checker started")
}

type StatusManager struct {
	statusCfg StatusConfig
	port      int64
	dontCheck map[string]bool
}

func (mgr StatusManager) Start() {
	mgr.server()
	mgr.loop()
}

func filter(netName, app string) string {
	return fmt.Sprintf("%s_%s", netName, app)
}
func (mgr StatusManager) server() {
	mux := http.NewServeMux()
	mux.HandleFunc("/dontCheck/get", func(resp http.ResponseWriter, req *http.Request) {
		WriteSuccess(resp, mgr.dontCheck)
	})
	mux.HandleFunc("/dontCheck/update", func(resp http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		netName := query.Get("network")
		application := query.Get("application")
		op := query.Get("operation")
		if mgr.statusCfg[netName] == nil || mgr.statusCfg[netName][application] == nil {
			WriteErr(resp, 400, fmt.Errorf("Unknown application %s_%s", netName, application))
			return
		}
		f := filter(netName, application)
		switch op {
		case "add":
			mgr.dontCheck[f] = true
		case "remove":
			delete(mgr.dontCheck, f)
		default:
			log.Warn("Unknown operation", op)
		}
		WriteSuccess(resp, mgr.dontCheck)
	})
	//
	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", mgr.port),
		Handler: mux,
	}
	go srv.ListenAndServe()
}
func (mgr StatusManager) loop() {
	for {
		mgr.checkStatus()
		time.Sleep(5 * time.Minute)
	}
}

func (mgr StatusManager) checkStatus() {
	for netName, statusCfg := range mgr.statusCfg {
		for app, urls := range statusCfg {
			if mgr.dontCheck[filter(netName, app)] {
				continue
			}
			for _, url := range urls {
				resp, err := http.Get(url)
				if err != nil {
					log.Errorf("Error(%s): %s_%s[%s] is down.", err.Error(), netName, app, url)
				} else if resp.StatusCode/100 != 2 {
					{
						log.Errorf("Error(%s): %s_%s[%s] is down.", resp.Status, netName, app, url)
					}
				}
			}
		}
		log.Infof("Checked %d applications for %s", len(statusCfg), netName)
	}
}
func newStatusManager(port int64) *StatusManager {
	cfg := map[string]ApplicationsUrl{}
	dataStr, err := core.GetJsonnetFile("config.jsonnet", core.JsonnetImports{})
	log.CheckFatal(err)
	reader := bytes.NewBuffer([]byte(dataStr))
	utils.ReadJsonReaderAndSetInterface(reader, &cfg)
	return &StatusManager{
		statusCfg: (StatusConfig)(cfg),
		port:      port,
		dontCheck: map[string]bool{},
	}
}

// utils
func WriteSuccess(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(utils.ToJsonBytes(map[string]interface{}{"data": data}))
}

func WriteErr(w http.ResponseWriter, code int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(utils.ToJsonBytes(map[string]string{"message": err.Error()}))
}
