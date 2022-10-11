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
		log.Error("New file watcher error:", err)
	}
	defer func(watcher *fsnotify.Watcher) {
		err := watcher.Close()
		if err != nil {
			log.Panic("Close fsnotify watcher, error:", err)
		}
	}(watcher)

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:

				// Channel was closed (i.e. Watcher.Close() was called).
				if !ok {
					log.Error("Event Channel was closed (i.e. Watcher.Close() was called):", err)
					return
				}
				log.Printf("File:%s is  %s", event.Name, event.Op)
				if (event.Op&fsnotify.Create == fsnotify.Create) || (event.Op&fsnotify.Write == fsnotify.Write) {
					compileXml := regexp.MustCompile(RE_XML_FILE)
					filename := event.Name

					if compileXml.MatchString(filename) {
						softpak.SendExportXml(filename, dc)
					}
				}
			case err, ok := <-watcher.Errors:
				// Channel was closed (i.e. Watcher.Close() was called).
				if !ok {
					log.Error("File watcher Error Channel was closed (i.e. Watcher.Close() was called):", err)
					return
				}
				log.Error("File read error:", err)
			}
		}
	}()

	err = watcher.Add(dir)

	if err != nil {
		log.Fatal(err)
	}
	<-done
}
