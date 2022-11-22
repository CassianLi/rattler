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

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

// DownloadTaxPdf Download the export PDF file
// @Summary      下载税金单文件
// @Description  通过申报国家确定税金文件路径
// @Tags         download
// @Accept       json
// @Produce      json
// @Param        origin   path  string   true  "下载的源文件名,不带文件后缀"
// @Param        target   path  string   true  "下载文件后，将文件重命名的文件名，没有后缀将自动添加pdf作为后缀"
// @Param        dc   	  query string   false  "申报国家(BE|NL),默认为NL"
// @Success      200
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /download/pdf/{origin}/{target} [get]
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
// @Summary      下载export文件
// @Description  通过文件名前缀确定文件路径
// @Tags         download
// @Accept       json
// @Produce      json
// @Param        dc   	  path  string   true  "申报国家(BE|NL)"
// @Param        filename path  string   true  "export文件的文件名"
// @Param        download query string   false  "是否下载，1表示直接下载"
// @Success      200
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /download/xml/{dc}/{filename} [get]
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
// @Summary      搜索文件
// @Description  可检索税金单文件以及export报关结果文件，可使用文件名部分做模糊匹配。建议使用Job number 进行检索
// @Tags         search
// @Accept       json
// @Produce      json
// @Param        message  body  SearchFileRequest  true  "检索内容"
// @Success      200 {object} []softpak.SearchFileResult
// @Failure      400 {object} util.ResponseError
// @Failure      404 {object} util.ResponseError
// @Failure      500 {object} util.ResponseError
// @Router       /search/file [post]
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
		return c.JSON(http.StatusBadRequest, &util.ResponseError{
			Status: util.FAIL,
			Errors: errs,
		})
	}

	sf := &softpak.SearchFile{
		DeclareCountry: sfd.DeclareCountry,
		Type:           sfd.Type,
		Filenames:      sfd.Filenames,
	}
	files, errs := sf.GetSearchResult()
	if len(errs) > 0 {
		return c.JSON(http.StatusBadRequest, &util.ResponseError{
			Status: util.FAIL,
			Errors: errs,
		})
	}

	// success
	return c.JSON(http.StatusOK, files)
}

// ExportListenFiles 获取Export监听路径下当前文件列表
// @Summary      获取Export监听路径下当前文件列表
// @Description  通过指定的申报国家获取其当前Export监听路径下的文件列表
// @Tags         export
// @Accept       json
// @Produce      json
// @Param        dc  path  	  string   true  "申报国家(BE|NL)"
// @Success      200 {object} []softpak.ExportFileListDTO
// @Failure      400 {object} util.ResponseError
// @Failure      404 {object} util.ResponseError
// @Failure      500 {object} util.ResponseError
// @Router       /export/list/{dc} [get]
func ExportListenFiles(c echo.Context) error {
	dc := strings.ToUpper(c.Param("dc"))
	fmt.Println(dc)

	data, err := softpak.ExportListenDicFiles(dc)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &util.ResponseError{
			Status: util.FAIL,
			Errors: []string{
				err.Error(),
			},
		})
	}

	// success
	return c.JSON(http.StatusOK, data)
}
