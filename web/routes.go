package web

import (
	"fmt"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"path/filepath"
	"strings"
	"sysafari.com/softpak/rattler/softpak"
	"sysafari.com/softpak/rattler/util"
)

// DownloadTaxPdf Download the export PDF file
// Download the export PDF file ,Specify the file name for the download
func DownloadTaxPdf(c echo.Context) error {
	nlTaxDir := viper.GetString("ser-dir.nl.tax-bill")
	beTaxDir := viper.GetString("ser-dir.be.tax-bill")

	origin := c.Param("origin") + ".pdf"
	target := c.Param("target")
	if !strings.Contains(target, ".pdf") {
		target = target + ".pdf"
	}

	dc := strings.ToUpper(c.QueryParam("dc"))

	var filePath string
	// dc 为空则为nl
	if "BE" == dc {
		filePath = filepath.Join(beTaxDir, origin)
	} else {
		filePath = filepath.Join(nlTaxDir, origin)
	}

	if util.IsExists(filePath) {
		return c.Attachment(filePath, target)
	}

	log.Errorf("Download tax-bill pdf failed,%s is not found", filePath)

	return c.String(http.StatusNotFound,
		fmt.Sprintf("Download tax-bill pdf failed,%s is not found.", origin))
}

// DownloadExportXml Download the export XML file
// Only asl user can access this api
func DownloadExportXml(c echo.Context) error {
	nlExportDir := viper.GetString("ser-dir.nl.export")
	beExportDir := viper.GetString("ser-dir.be.export")

	dc := strings.ToUpper(c.Param("dc"))
	filename := c.Param("filename")

	year := filename[0:4]
	month := filename[4:6]

	needDownload := c.QueryParam("download")

	var filePath string
	if "NL" == dc {
		filePath = filepath.Join(nlExportDir, year, month, filename)
	} else if "BE" == dc {
		filePath = filepath.Join(beExportDir, filename)
	} else {
		return c.String(http.StatusNotFound, fmt.Sprintf("%s is not a valid declare country", dc))
	}

	fmt.Println("download filePath:", filePath)
	if util.IsExists(filePath) {
		if "1" == needDownload {
			return c.Attachment(filePath, filename)
		}
		return c.File(filePath)
	}

	log.Errorf("Download export xl failed,%s is not found.", filePath)

	return c.String(http.StatusNoContent, fmt.Sprintf("The file %s is not found", filename))
}

// SearchFile Search for tax bill files and Export declaration XML files
func SearchFile(c echo.Context) (err error) {
	var errs []string
	sfd := new(SearchFileRequest)
	if err = c.Bind(sfd); err != nil {
		errs = append(errs, err.Error())
	}
	if err = c.Validate(sfd); err != nil {
		errs = append(errs, err.Error())
	}

	if len(errs) > 0 {
		return c.JSON(http.StatusBadRequest, &SearchFileResponse{
			Status: FAIL,
			Errors: errs,
		})
	}

	sf := &softpak.SearchFile{
		Type:      sfd.Type,
		Filenames: sfd.Filenames,
	}
	files, errs := sf.GetSearchResult()
	var res SearchFileResponse
	if len(errs) > 0 {
		res = SearchFileResponse{
			Status: FAIL,
			Errors: errs,
		}
	} else {
		res = SearchFileResponse{
			Status: SUCCESS,
			Files:  files,
		}
	}

	// success
	return c.JSON(http.StatusOK, res)
}
