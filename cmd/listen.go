/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"sysafari.com/softpak/rattler/rabbit"
	"sysafari.com/softpak/rattler/softpak"

	"github.com/spf13/cobra"
)

// listenCmd represents the listen command
var listenCmd = &cobra.Command{
	Use:   "listen",
	Short: "Monitor the messages in the message queue and save the corresponding messages as XML files to the specified path",
	Long:  `Listen for new messages in the message queue, and save the new customs declaration request in the form of an XML file in the disk path. For example:`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("listen called")
		listenAmqpForImportXml()
	},
}

// Listen to the message queue and
// save the Import xml to the specified path
func listenAmqpForImportXml() {
	rbmq := &rabbit.Rabbit{
		Url:      viper.GetString("rabbitmq.url"),
		Exchange: viper.GetString("rabbitmq.exchange"),
		Queue:    viper.GetString("rabbitmq.queue"),
	}

	log.Infof("Starting ... RabbitMQ consumer: %v ", rbmq)
	rabbit.Consume(rbmq, "topic", softpak.SaveImportDocument)
}

func init() {
	rootCmd.AddCommand(listenCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listenCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listenCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
