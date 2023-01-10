package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
	rabbit2 "sysafari.com/softpak/rattler/dao/rabbit"
	"sysafari.com/softpak/rattler/model"
	"sysafari.com/softpak/rattler/util"
	"time"
)

// 报关结果放行文件服务类

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

	// backup export xml
	fn, err := moveFileToBackup(filename, declareCountry)
	if err != nil {
		// Backup failed send original file name
		fn = filepath.Base(filename)
	}

	xmlContent := ExportXmlInfo{
		FileName:       fn,
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

	rbmq := &rabbit2.Rabbit{
		Url:          viper.GetString("rabbitmq.url"),
		Exchange:     viper.GetString("rabbitmq.export.exchange"),
		ExchangeType: viper.GetString("rabbitmq.export.exchange-type"),
		Queue:        queueName,
	}

	rabbit2.Publish(rbmq, message)
}

// moveFileToBackup Move file to back up location
func moveFileToBackup(fp string, dc string) (string, error) {
	fn := filepath.Base(fp)

	firstPt := strings.Split(fn, "_")[0]
	parse, err := time.Parse("200601", firstPt)
	var year, month, newFileName string
	if err != nil {
		year = time.Now().Format("2006")
		month = time.Now().Format("01")
		newFileName = fmt.Sprintf("%s%s_%s", year, month, fn)
	} else {
		log.Warnf("The file:%s within date ,backup is origin filename.", fn)
		year = parse.Format("2006")
		month = parse.Format("01")
		newFileName = fn
	}

	backupDir := viper.GetString(fmt.Sprintf("watcher.%s.backup-dir", strings.ToLower(dc)))
	bacdir := filepath.Join(backupDir, year, month)
	// backup directory not exists create it
	canMove := util.IsDir(bacdir) || util.CreateDir(bacdir)
	if !canMove {
		log.Errorf("Cannot create backup dir %s , dont move file %s", bacdir, fp)
		return "", errors.New(fmt.Sprintf("Cannot create backup dir %s , dont move file %s", bacdir, fp))
	}
	filename := filepath.Base(fp)
	targetFilename := filepath.Join(bacdir, newFileName)

	err = os.Rename(fp, targetFilename)
	if err != nil {
		log.Errorf("Backup export file %s failed, error: %v", filename, err)
		return "", err
	}
	return newFileName, nil
}

// ExportListenDicFiles 获取申报国家Export 监听路径下的文件列表
func ExportListenDicFiles(dc string) (files []model.ExportFileListDTO, err error) {
	var listenDir string
	if "NL" == dc {
		listenDir = viper.GetString("watcher.nl.watch-dir")
	}
	if "BE" == dc {
		listenDir = viper.GetString("watcher.be.watch-dir")
	}
	fmt.Println("listenDir", listenDir)
	if !util.IsDir(listenDir) || !util.IsExists(listenDir) {
		return nil, errors.New("the monitoring path is wrong. Check whether the declared country exists")
	}

	// 获取文件列表
	var fs []string
	err = filepath.Walk(listenDir, util.Visit(&fs))
	if err != nil {
		return nil, err
	}
	fmt.Println("fs", fs)

	for _, f := range fs {
		info, err := os.Stat(f)
		if err == nil {
			ef := model.ExportFileListDTO{
				Filename: filepath.Base(f),
				Filepath: f,
				Size:     info.Size(),
				ModTime:  info.ModTime().Format("2006-01-02 15:04:05"),
			}
			files = append(files, ef)
		} else {
			log.Errorf("File: %s get stat failed, error: %v", f, err)
		}
	}

	return files, err
}
