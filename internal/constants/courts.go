package constants

// Courts maps CLI court values to their API Applikation values
// for the Judikatur endpoint.
var Courts = map[string]string{
	"justiz":      "Justiz",
	"vfgh":        "Vfgh",
	"vwgh":        "Vwgh",
	"bvwg":        "Bvwg",
	"lvwg":        "Lvwg",
	"dsk":         "Dsk",
	"asylgh":      "AsylGH",
	"normenliste": "Normenliste",
	"pvak":        "Pvak",
	"gbk":         "Gbk",
	"dok":         "Dok",
}

// LeitsatzCourts lists courts that support Leitsatz extraction.
var LeitsatzCourts = map[string]bool{
	"Vfgh":   true,
	"Vwgh":   true,
	"Justiz": true,
	"Bvwg":   true,
}

// ValidCourts returns a list of valid CLI court values.
func ValidCourts() []string {
	return []string{
		"justiz", "vfgh", "vwgh", "bvwg", "lvwg",
		"dsk", "asylgh", "normenliste", "pvak", "gbk", "dok",
	}
}
