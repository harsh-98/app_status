package main

import (
	"bytes"
	"net/http"
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
			ChainId:  1,
		},
		"TEST_app_status_checker",
	)
	log.AMQPMsg("App checker started")
}

type StatusConfig struct {
	Mainnet ApplicationsUrl `json:"mainnet"`
}

type ApplicationsUrl map[string][]string

func main() {
	startLogging()
	//
	statusCfg := loadStatusConfig()
	for {
		checkStatus(statusCfg.Mainnet)
		time.Sleep(5 * time.Minute)
	}
}
func checkStatus(statusCfg ApplicationsUrl) {
	for app, urls := range statusCfg {
		for _, url := range urls {
			resp, err := http.Get(url)
			if err != nil {
				log.Errorf("Error(%s): %s[%s] is down.", err.Error(), app, url)
			} else if resp.StatusCode/100 != 2 {
				{
					log.Errorf("Error(%s): %s[%s] is down.", resp.Status, app, url)
				}
			}
		}
	}
	log.Infof("Checked %d applications", len(statusCfg))
}
func loadStatusConfig() *StatusConfig {
	cfg := &StatusConfig{}
	dataStr, err := core.GetJsonnetFile("config.jsonnet", core.JsonnetImports{})
	log.CheckFatal(err)
	reader := bytes.NewBuffer([]byte(dataStr))
	utils.ReadJsonReaderAndSetInterface(reader, cfg)
	return cfg
}
