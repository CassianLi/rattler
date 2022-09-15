package softpak

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"sysafari.com/softpak/rattler/util"
)

type ImportDocument struct {
	Filename string `json:"filename"`
	Document string `json:"document"`
}

// SaveImportDocument saves the import xml document
func SaveImportDocument(message string) {
	// 去除转义符
	msg, err := strconv.Unquote(message)
	doc := ImportDocument{}
	if err != nil {
		err = json.Unmarshal([]byte(message), &doc)
	} else {
		err = json.Unmarshal([]byte(msg), &doc)
	}

	if err != nil {
		log.Errorf("Unmarshal queue message, err: %v", err)
		fmt.Println("Unmarshal queue message, err: ", err)
		return
	}

	filename := doc.Filename
	document := doc.Document
	importDir := viper.GetString("import.xml-dir")

	canSave := util.IsDir(importDir) || util.CreateDir(importDir)
	if !canSave {
		log.Errorf("Import directory %s not exists, dont save import xml document", importDir)
		return
	}

	fp := filepath.Join(importDir, filename)
	err = ioutil.WriteFile(fp, []byte(document), os.ModePerm)
	if err != nil {
		log.Errorf("Write file %s error: %v", fp, err)
	} else {
		log.Infof("Write file %s ok", fp)
	}

}
