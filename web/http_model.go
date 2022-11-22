package web

import (
	"github.com/go-playground/validator"
)

type (
	CustomValidator struct {
		Validator *validator.Validate
	}

	SearchFileRequest struct {
		// DeclareCountry NL, BE
		DeclareCountry string `json:"declareCountry" validate:"required"`
		// Type TAX_BILL, EXPORT_XML
		Type      string   `json:"type" validate:"required"`
		Filenames []string `json:"filenames" validate:"required"`
	}
)
