package service

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"sysafari.com/softpak/rattler/componet"
	"sysafari.com/softpak/rattler/config"
	"sysafari.com/softpak/rattler/model"
	"sysafari.com/softpak/rattler/util"
	"time"
)

// 重新发送export 文件服务

func moveOneExportFile(file string, dc string, inListenDir bool) (err error) {
	var exportWatchDir string
	var remover *componet.RemoveQueue
	if dc == "NL" {
		exportWatchDir = viper.GetString("watcher.nl.watch-dir")
		remover = config.NlRemover
	}

	if dc == "BE" {
		exportWatchDir = viper.GetString("watcher.be.watch-dir")
		remover = config.BeRemover
	}

	if exportWatchDir == "" {
		log.Panicf("Watch directory (DC:%s = %s) is empty，Error!!", dc, exportWatchDir)
	}

	tmpDir := viper.GetString("tmp-dir")

	if !util.IsExists(file) {
		return errors.New(fmt.Sprintf("The file:%s not exist", file))
	}

	fn := filepath.Base(file)

	if inListenDir {
		fileTmp := filepath.Join(tmpDir, fn)
		err = os.Rename(file, fileTmp)
		if err != nil {
			return err
		}
		time.Sleep(1 * time.Second)
		file = fileTmp
	}
	// 加入移动文件路径
	remover.Add(componet.RemoveParam{
		SourceFile: file,
		MoveTo:     exportWatchDir,
	})

	return nil
}

// ResendExport 重新发送Export 文件
func ResendExport(dc string, params *model.FileResendRequest) (errs []string) {
	for _, path := range params.FilePaths {
		err := moveOneExportFile(path, dc, params.InListeningPath)
		if err != nil {
			errs = append(errs, err.Error())
		}
	}
	return errs
}
