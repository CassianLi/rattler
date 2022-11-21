package softpak

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"path/filepath"
	"strings"
	"sysafari.com/softpak/rattler/util"
)

// SearchFile Search softpak file
type SearchFile struct {
	DeclareCountry string `json:"declare_country"`
	// Type TAX_BILL, EXPORT_XML
	Type         string             `json:"type"`
	Directory    string             `json:"directory"`
	Filenames    []string           `json:"filenames"`
	SearchResult []SearchFileResult `json:"searchResult"`
	Errors       []string           `json:"errors"`
}

// ready Ready for search
func (sf SearchFile) ready() {
	if "NL" == sf.DeclareCountry {
		if "TAX_BILL" == sf.Type {
			sf.Directory = viper.GetString("ser-dir.nl.tax-bill")
		}
		if "EXPORT_XML" == sf.Type {
			sf.Directory = viper.GetString("watcher.nl.backup-dir")
		}
	}

	if "BE" == sf.DeclareCountry {
		if "TAX_BILL" == sf.Type {
			sf.Directory = viper.GetString("ser-dir.be.tax-bill")
		}
		if "EXPORT_XML" == sf.Type {
			sf.Directory = viper.GetString("watcher.be.backup-dir")
		}
	}

	if !util.IsDir(sf.Directory) || !util.IsExists(sf.Directory) {
		sf.Errors = append(sf.Errors, fmt.Sprintf("The file directory %s not exists", sf.Directory))
	}
}

// search Start to search file
func (sf *SearchFile) search() {
	if len(sf.Errors) > 0 {
		return
	}

	var files []string
	err := filepath.Walk(sf.Directory, util.Visit(&files))
	if err != nil {
		sf.Errors = append(sf.Errors, fmt.Sprintf("Failed to get all file names through the directory:%s", sf.Directory))
		return
	}
	log.Infof("The directory: %s contains %d files.", sf.Directory, len(files))

	for _, file := range files {
		filename := filepath.Base(file)
		for _, s := range sf.Filenames {
			if strings.Contains(filename, s) {
				sfr := SearchFileResult{
					Type:       sf.Type,
					SearchText: s,
					Filename:   filename,
					Filepath:   file,
				}
				sf.SearchResult = append(sf.SearchResult, sfr)
				break
			}
		}
	}
}

// GetSearchResult Begin to search
func (sf *SearchFile) GetSearchResult() ([]SearchFileResult, []string) {
	sf.ready()
	sf.search()

	return sf.SearchResult, sf.Errors
}
