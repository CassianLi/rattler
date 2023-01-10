package cmd

import (
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	echoSwagger "github.com/swaggo/echo-swagger"
	"sysafari.com/softpak/rattler/model"
	"sysafari.com/softpak/rattler/route"
)

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
	e.Validator = &model.CustomValidator{
		Validator: validator.New(),
	}

	// swagger
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.GET("/download/pdf/:origin/:target", route.DownloadTaxPdf)
	e.GET("/download/xml/:dc/:filename", route.DownloadExportXml)

	e.POST("/search/file", route.SearchFile)

	// Export 监听路径下的文件列表
	e.GET("/export/list/:dc", route.ExportListenFiles)

	// http://domain.com/export/resend
	e.POST("/export/remover/:dc", route.ExportFileResend)

	port := viper.GetString("port")
	if port == "" {
		port = "1324"
	}

	log.Errorf("Rattler server started: %v", e.Start(":"+port))
}
