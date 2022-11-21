package softpak

type (
	SearchFileResult struct {
		// Type TAX_BILL, EXPORT_XML
		Type       string `json:"type"`
		SearchText string `json:"searchText"`
		Filename   string `json:"filename"`
		Filepath   string `json:"filepath"`
	}
)
