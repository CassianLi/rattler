package web

import (
	"github.com/go-playground/validator"
	"sysafari.com/softpak/rattler/softpak"
)

const (
	SUCCESS = "success"
	FAIL    = "fail"
)

type (
	SearchFileRequest struct {
		// Type TAX_BILL, EXPORT_XML
		Type      string   `json:"type" validate:"required"`
		Filenames []string `json:"filenames" validate:"required"`
	}

	SearchFileValidator struct {
		Validator *validator.Validate
	}

	SearchFileResponse struct {
		Status string                     `json:"status"`
		Errors []string                   `json:"errors"`
		Files  []softpak.SearchFileResult `json:"files"`
	}
)
