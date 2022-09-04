package softpak

import (
	"fmt"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"strings"
	"sysafari.com/softpak/rattler/util"
)

// File server of soft pak

// DownloadTaxPdf Download the export PDF file
// Download the export PDF file ,Specify the file name for the download
func DownloadTaxPdf(c echo.Context) error {
	nlTaxDir := viper.GetString("directory.tax-bill.nl")
	beTaxDir := viper.GetString("directory.tax-bill.be")

	origin := c.Param("origin") + ".pdf"
	target := c.Param("target") + ".pdf"
	dc := strings.ToUpper(c.QueryParam("dc"))

	var filePath string
	if "NL" == dc {
		filePath = nlTaxDir + "/" + origin
	} else if "BE" == dc {
		filePath = beTaxDir + "/" + origin
	} else {
		return c.String(http.StatusNotFound, fmt.Sprintf("%s is not a valid declare country", dc))
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
	nlExportDir := viper.GetString("directory.export.nl")
	beExportDir := viper.GetString("directory.export.be")

	dc := strings.ToUpper(c.Param("dc"))
	filename := c.Param("filename")

	needDownload := c.QueryParam("download")

	var filePath string
	if "NL" == dc {
		filePath = nlExportDir + "/" + filename
	} else if "BE" == dc {
		filePath = beExportDir + "/" + filename
	} else {
		return c.String(http.StatusNotFound, fmt.Sprintf("%s is not a valid declare country", dc))
	}

	if util.IsExists(filePath) {
		if "1" == needDownload {
			return c.Attachment(filePath, filename)
		}
		return c.File(filePath)
	}

	log.Errorf("Download export xl failed,%s is not found.", filePath)

	return c.String(http.StatusNoContent, fmt.Sprintf("The file %s is not found", filename))
}
