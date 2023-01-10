/*
Copyright © 2022 Joker

*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"sysafari.com/softpak/rattler/config"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "rattler",
	Short: "SoftPak Client",
	Long: `Rattler will simultaneously start the file server to access 
the Export XML file and the tax bill file, and simultaneously start the 
Import XML listener and the Export XML (NL|BE) file creation listener asynchronously. 
For example:`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		// async start listen import xml queue
		go ListenAmqpForImportXml()

		// async start listen export xml directory(NL | BE)
		go ListenExportXML("NL")
		go ListenExportXML("BE")

		// start remover to resend export xml
		go RemoverWork()

		// start file server
		EchoRoutes()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}

}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", ".rattler.yaml", "config file (default is .rattler.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".rattler" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".rattler")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	initLog()

}

func initLog() {
	path, _ := os.Executable()
	_, exec := filepath.Split(path)
	logFile := exec + ".log"
	logFilePath := filepath.Join(viper.GetString("log.directory"), logFile)

	config.InitLog(logFilePath, viper.GetString("log.level"))
}
