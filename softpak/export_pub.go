package softpak

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
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
	log.Infof("Declare country: %s export xml: %s reading ", declareCountry, filename)

	content, err := os.ReadFile(filename)
	if err != nil {
		log.Error("Read XML file error:", err)
	}
	contentStr := string(content)

	log.Debugf("Min size xml content:  %s ", contentStr)

	xmlContent := ExportXmlInfo{
		FileName:       filepath.Base(filename),
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
		moveFileToBackup(filename, declareCountry)

		//jobNumber, _ := getJobNumber(filename)
		// Send xml info to MQ
		publishMessageToMQ(bf.String(), declareCountry)
	}

}

// publishMessageToMQ publishes the message to MQ
func publishMessageToMQ(message string, declareCountry string) {
	qPrefix := viper.GetString("rabbitmq.export.queue")

	//seq := strconv.Itoa(jobNumber % queueCount)
	//fmt.Println(seq)
	var queueName = strings.ToLower(qPrefix + "." + declareCountry)

	fmt.Println(queueName)
	rbmq := &rabbit.Rabbit{
		Url:          viper.GetString("rabbitmq.url"),
		Exchange:     viper.GetString("rabbitmq.export.exchange"),
		ExchangeType: viper.GetString("rabbitmq.export.exchange-type"),
		Queue:        queueName,
	}

	rabbit.Publish(rbmq, message)
}

// moveFileToBackup Move file to back up location
func moveFileToBackup(fp string, dc string) {
	backupDir := viper.GetString(fmt.Sprintf("watcher.%s.backup-dir", strings.ToLower(dc)))
	if len(backupDir) == 0 {
		log.Errorf("Backup dir is empty, cannot move file %s", fp)
	} else {
		// backup directory not exists create it
		canMove := util.IsDir(backupDir) || util.CreateDir(backupDir)
		if !canMove {
			log.Errorf("Cannot create backup dir %s , dont move file %s", backupDir, fp)
		} else {
			filename := filepath.Base(fp)
			targetFilename := filepath.Join(backupDir, filename)

			err := os.Rename(fp, targetFilename)
			if err != nil {
				log.Errorf("Backup export file %s failed, error: %v", filename, err)
			} else {
				log.Infof("Backup file %s moved to %s", fp, targetFilename)
			}
		}
	}
}
