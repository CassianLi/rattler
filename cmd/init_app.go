package cmd

import (
	"fmt"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"strings"
	"sysafari.com/softpak/rattler/rabbit"
	"sysafari.com/softpak/rattler/softpak"
	"sysafari.com/softpak/rattler/util"
	"sysafari.com/softpak/rattler/watcher"
)

// ListenAmqpForImportXml Listen to the message queue and
// save the Import xml to the specified path
func ListenAmqpForImportXml() {
	rbmq := &rabbit.Rabbit{
		Url:          viper.GetString("rabbitmq.url"),
		Exchange:     viper.GetString("rabbitmq.import.exchange"),
		ExchangeType: viper.GetString("rabbitmq.import.exchange-type"),
		Queue:        viper.GetString("rabbitmq.import.queue"),
	}

	log.Infof("Starting ... RabbitMQ consumer: %v ", rbmq)
	rabbit.Consume(rbmq, softpak.SaveImportDocument)
}

// EchoRoutes Set routes to echo
func EchoRoutes() {
	e := echo.New()

	e.GET("/download/pdf/:origin/:target", softpak.DownloadTaxPdf)
	e.GET("/download/xml/:dc/:filename", softpak.DownloadExportXml)

	port := viper.GetString("port")
	if port == "" {
		port = "1324"
	}

	log.Errorf("Rattler server started: %v", e.Start(":"+port))
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
	watcher.Watch(watchDir, declareCountry)
}
