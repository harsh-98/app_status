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
		cfg.AppName,
	)
	log.AMQPMsg("App checker started")
}

type StatusConfig map[string]ApplicationsUrl

type ApplicationsUrl map[string][]string

func main() {
	startLogging()
	//
	statusCfg := loadStatusConfig()
	for {
		for netName, apps := range statusCfg {
			checkStatus(netName, apps)

		}
		time.Sleep(5 * time.Minute)
	}
}
func checkStatus(netName string, statusCfg ApplicationsUrl) {
	for app, urls := range statusCfg {
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
func loadStatusConfig() StatusConfig {
	cfg := map[string]ApplicationsUrl{}
	dataStr, err := core.GetJsonnetFile("config.jsonnet", core.JsonnetImports{})
	log.CheckFatal(err)
	reader := bytes.NewBuffer([]byte(dataStr))
	utils.ReadJsonReaderAndSetInterface(reader, &cfg)
	return cfg
}
