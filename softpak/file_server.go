package softpak

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"net/http"
	"strings"
	"sysafari.com/softpak/rattler/util"
)

// File server of soft pak

// DownloadExportPdf Download the export PDF file
// Download the export PDF file ,Specify the file name for the download
func DownloadExportPdf(c echo.Context) error {
	dir := viper.GetString("directory.tax-bill")

	origin := c.Param("origin") + ".pdf"
	target := c.Param("target") + ".pdf"

	originFilename := dir + "/" + origin
	fmt.Println(originFilename)

	if util.IsExists(originFilename) {
		return c.Attachment(originFilename, target)
	}

	return c.String(http.StatusNotFound, "The file "+origin+" is not found")
}

// DownloadExportXml Download the export XML file
// Only asl user can access this api
func DownloadExportXml(c echo.Context) error {
	nlExportDir := viper.GetString("directory.export.nl")
	beExportDir := viper.GetString("directory.export.be")

	dc := strings.ToUpper(c.Param("declareCountry"))
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

	return c.String(http.StatusNoContent, fmt.Sprintf("The file %s is not found", filename))
}
