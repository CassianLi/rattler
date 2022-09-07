package softpak

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"strconv"
	"strings"
	"sysafari.com/softpak/rattler/rabbit"
	"sysafari.com/softpak/rattler/util"
)

type WatchConfig struct {
	Watch     bool
	WatchDir  string
	BackupDir string
}

// Dc Declare country
type Dc uint32

// Declare country enum
const (
	NL Dc = 1 << iota
	BE
)

// ExportXmlInfo Export XML file information
type ExportXmlInfo struct {
	FileName       string `json:"fileName"`
	DeclareCountry string `json:"declareCountry"`
	Content        string `json:"content"`
}

// SendExportXml sends export Xml file to the MQ
// Compress the content of the XML file before sending,
// and then create a json object and send it to the message queue
func SendExportXml(filename string, declareCountry string) {
	//backupDir := viper.GetString("watcher.nl.backup-dir")
	log.Infof("Declare country: %s export xml: %s reading ", declareCountry, filename)

	content, err := os.ReadFile(filename)
	if err != nil {
		log.Error("Read XML file error:", err)
	}

	contentStr := string(content)

	log.Debugf("Min size xml content:  %s ", contentStr)
	fps := strings.Split(filename, "/")

	xmlContent := ExportXmlInfo{
		FileName:       fps[len(fps)-1],
		DeclareCountry: declareCountry,
		Content:        contentStr,
	}
	// Serialize to JSON
	bf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(bf)
	jsonEncoder.SetEscapeHTML(false)
	err = jsonEncoder.Encode(xmlContent)

	if err != nil {
		log.Error("Serialize Export xml file to JSON failed, dont publish. ", err)
	} else {
		// backup export xml
		moveFileToBackup(filename)

		//jobNumber, _ := getJobNumber(filename)
		// Send xml info to MQ
		publishMessageToMQ(bf.String(), declareCountry)
	}

}

// getJobNumber returns the job number of Export xml
func getJobNumber(filename string) (int, error) {
	fns := strings.Split(filename, "/")
	fn := fns[len(fns)-1]

	jbn := strings.Split(strings.Split(fn, "_")[1], "-")[2]
	log.Infof("The file: %s's job number is: %s", filename, jbn)

	return strconv.Atoi(jbn)
}

// publishMessageToMQ publishes the message to MQ
func publishMessageToMQ(message string, declareCountry string) {
	qPrefix := viper.GetString("rabbitmq.queue")

	//seq := strconv.Itoa(jobNumber % queueCount)
	//fmt.Println(seq)
	var queueName = strings.ToLower(qPrefix + "." + declareCountry)

	fmt.Println(queueName)
	rbmq := &rabbit.Rabbit{
		Url:          viper.GetString("rabbitmq.url"),
		Exchange:     viper.GetString("rabbitmq.exchange"),
		ExchangeType: viper.GetString("rabbitmq.exchange-type"),
		Queue:        queueName,
	}

	rabbit.Publish(rbmq, message)
}

// moveFileToBackup Move file to back up location
func moveFileToBackup(filename string) {
	backupDir := viper.GetString("watcher.backup-dir")
	if len(backupDir) == 0 {
		log.Errorf("Backup dir is empty, cannot move file %s", filename)
	} else {
		// backup directory not exists create it
		canMove := util.IsDir(backupDir) || util.CreateDir(backupDir)
		if !canMove {
			log.Errorf("Cannot create backup dir %s , dont move file %s", backupDir, filename)
		} else {
			fnpaths := strings.Split(filename, "/")
			fn := fnpaths[len(fnpaths)-1]
			targetFilename := backupDir + "/" + fn

			err := os.Rename(filename, targetFilename)
			if err != nil {
				log.Errorf("Backup export file %s failed, error: %v", filename, err)
			} else {
				log.Infof("Backup file %s moved to %s", filename, targetFilename)
			}
		}
	}
}
