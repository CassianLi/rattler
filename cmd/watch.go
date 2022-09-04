/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
	"sysafari.com/softpak/rattler/watcher"
)

const (
	NL = "NL"
	BE = "BE"
)

// watchCmd represents the watch command
var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Monitor changes to files in the specified file path",
	Long: `Monitor the path of the customs declaration result file. 
	Once a new customs declaration result file is generated, 
	the content of the file will be read immediately and sent to the message 
	queue in the form of JSON. For example:`,

	Run: func(cmd *cobra.Command, args []string) {
		log.Info("watch called")
		initWatchConfig()
	},
}

func initWatchConfig() {
	declareCountry := strings.ToUpper(viper.GetString("declare-country"))
	if NL != declareCountry && BE != declareCountry {
		log.Panicf("%s is not a valid declare country(NL | BE)", declareCountry)
	}
	watchDir := viper.GetString("watcher.watch-dir")

	if len(watchDir) == 0 {
		log.Panic("Watch directory is empty，Error !!")
	}
	watcher.Watch(watchDir, declareCountry)
}

func init() {
	rootCmd.AddCommand(watchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	//watchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// watchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
