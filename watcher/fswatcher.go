package watcher

import (
	log "github.com/sirupsen/logrus"
	"regexp"
	"sysafari.com/softpak/rattler/softpak"

	"github.com/fsnotify/fsnotify"
)

// RE_XML_FILE XML file
const (
	RE_XML_FILE = ".*\\.xml"
)

// Watch Listen the directory for declare country
func Watch(dir string, dc string) {
	log.Debug("Watch dir: ", dir)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Error(err)
	}
	defer func(watcher *fsnotify.Watcher) {
		err := watcher.Close()
		if err != nil {
			log.Panic("error closing watcher:", err)
		}
	}(watcher)

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:

				if !ok {
					return
				}
				if event.Op&fsnotify.Create == fsnotify.Create {
					compileXml := regexp.MustCompile(RE_XML_FILE)
					filename := event.Name

					if compileXml.MatchString(filename) {
						softpak.SendExportXml(filename, dc)
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Error("error:", err)
			}
		}
	}()

	err = watcher.Add(dir)

	if err != nil {
		log.Fatal(err)
	}
	<-done
}
