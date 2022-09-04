package softpak

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"strconv"
	"sysafari.com/softpak/rattler/util"
)

type ImportDocument struct {
	Filename string `json:"filename"`
	Document string `json:"document"`
}

// 去除json中的转义字符
func disableEscapeHtml(data interface{}) (string, error) {
	bf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(bf)
	jsonEncoder.SetEscapeHTML(true)
	if err := jsonEncoder.Encode(data); err != nil {
		return "", err
	}
	return bf.String(), nil
}

// SaveImportDocument saves the import xml document
func SaveImportDocument(message string) {
	// 去除转义符
	msg, err := strconv.Unquote(message)
	if err != nil {
		log.Errorf("Un quote failde %s", message)
		return
	}
	doc := ImportDocument{}
	err = json.Unmarshal([]byte(msg), &doc)
	if err != nil {
		fmt.Println(err)
		return
	}
	filename := doc.Filename
	document := doc.Document
	importDir := viper.GetString("import-dir")

	canSave := util.IsDir(importDir) || util.CreateDir(importDir)
	if !canSave {
		log.Errorf("Import directory %s not exists, dont save import xml document", importDir)
	} else {
		filepath := importDir + "/" + filename
		err = ioutil.WriteFile(filepath, []byte(document), os.ModePerm)
		if err != nil {
			log.Errorf("Write file %s error: %v", filepath, err)
		} else {
			log.Infof("Write file %s ok", filepath)
		}
	}

}
