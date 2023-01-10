package cmd

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"sysafari.com/softpak/rattler/componet"
	"sysafari.com/softpak/rattler/config"
	rabbit2 "sysafari.com/softpak/rattler/dao/rabbit"
	"sysafari.com/softpak/rattler/service"

	"strings"
	_ "sysafari.com/softpak/rattler/docs" // docs is generated by Swag CLI, you have to import it.
	"sysafari.com/softpak/rattler/util"
)

// ListenAmqpForImportXml Listen to the message queue and
// save the Import xml to the specified path
func ListenAmqpForImportXml() {
	rbmq := &rabbit2.Rabbit{
		Url:          viper.GetString("rabbitmq.url"),
		Exchange:     viper.GetString("rabbitmq.import.exchange"),
		ExchangeType: viper.GetString("rabbitmq.import.exchange-type"),
		Queue:        viper.GetString("rabbitmq.import.queue"),
	}

	log.Infof("Starting ... RabbitMQ consumer: %v ", rbmq)
	rabbit2.Consume(rbmq, service.SaveImportDocument)
}

// ListenExportXML Export xml listen config
func ListenExportXML(declareCountry string) {
	if "NL" != declareCountry && "BE" != declareCountry {
		log.Panicf("%s is not a valid declare country(NL | BE)", declareCountry)
	}
	watchDir := viper.GetString(fmt.Sprintf("watcher.%s.watch-dir", strings.ToLower(declareCountry)))

	if !util.IsDir(watchDir) {
		log.Panicf("Watch directory: %s is not empty，Error !!", watchDir)
	}
	service.Watch(watchDir, declareCountry)
}

// RemoverWork Start the remover process
func RemoverWork() {
	var err error

	config.NlRemover = componet.InitRemoveQueue(200)
	if err != nil {
		log.Panicf("Create remover for NL failed: %s", err)
	}
	// 启动移动文件进程 NL
	go config.NlRemover.Run()

	config.BeRemover = componet.InitRemoveQueue(200)
	if err != nil {
		log.Panicf("Create remover for BE failed: %s", err)
	}
	// 启动移动文件进程 BE
	go config.BeRemover.Run()
}
