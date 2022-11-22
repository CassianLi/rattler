package softpak

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"sysafari.com/softpak/rattler/util"
)

// ExportListenDicFiles 获取申报国家Export 监听路径下的文件列表
func ExportListenDicFiles(dc string) (files []ExportFileListDTO, err error) {
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
			ef := ExportFileListDTO{
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
