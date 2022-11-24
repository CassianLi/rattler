package cmd

import (
	"fmt"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	echoSwagger "github.com/swaggo/echo-swagger"

	"strings"
	_ "sysafari.com/softpak/rattler/docs" // docs is generated by Swag CLI, you have to import it.
	"sysafari.com/softpak/rattler/rabbit"
	"sysafari.com/softpak/rattler/softpak"
	"sysafari.com/softpak/rattler/util"
	"sysafari.com/softpak/rattler/watcher"
	"sysafari.com/softpak/rattler/web"
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
// @title Rattler API
// @version 1.0
// @description This is a server for soft-pak, can download tax-bill and export xml files
// @termsOfService http://swagger.io/terms/

// @contact.name Joker
// @contact.email ljr@y-clouds.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:7003
func EchoRoutes() {
	e := echo.New()
	e.Validator = &web.CustomValidator{Validator: validator.New()}
	// swagger
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.GET("/download/pdf/:origin/:target", web.DownloadTaxPdf)
	e.GET("/download/xml/:dc/:filename", web.DownloadExportXml)

	e.POST("/search/file", web.SearchFile)

	// Export 监听路径下的文件列表
	e.GET("/export/list/:dc", web.ExportListenFiles)

	// http://domain.com/export/resend
	e.POST("/export/resend/:dc", web.DownloadExportXml)

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
