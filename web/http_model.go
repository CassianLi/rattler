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
		// Year exp: 2022
		Year string `json:"year"`
		// Month exp: 09
		Month string `json:"month"`
		// Type TAX_BILL, EXPORT_XML
		Type string `json:"type" validate:"required"`
		// Filenames Support use Job number
		Filenames []string `json:"filenames" validate:"required"`
	}
)
