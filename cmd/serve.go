/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"sysafari.com/softpak/rattler/softpak"

	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start a web server",
	Long:  `Start a web server, access files. default port: 1324 .`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("serve called")

		echoRoute()
	},
}

func echoRoute() {
	e := echo.New()

	e.GET("/download/pdf/:origin/:target", softpak.DownloadExportPdf)
	e.GET("/download/xml/:declareCountry/:filename", softpak.DownloadExportXml)

	port := viper.GetString("port")
	if port == "" {
		port = "1324"
	}

	log.Errorf("Rattler server started: %v", e.Start(":"+port))
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	//serveCmd.PersistentFlags().StringVar(&port, "port", "1324", "The port of the server to soft pak files,")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
