package api

const (
	// DefaultBaseURL is the base URL for the RIS OGD API v2.6.
	DefaultBaseURL = "https://data.bka.gv.at/ris/api/v2.6/"

	// Endpoint path segments (appended to BaseURL).
	EndpointBundesrecht = "Bundesrecht"
	EndpointLandesrecht = "Landesrecht"
	EndpointJudikatur   = "Judikatur"
	EndpointBezirke     = "Bezirke"
	EndpointGemeinden   = "Gemeinden"
	EndpointSonstige    = "Sonstige"
	EndpointHistory     = "History"
)

// AllowedHosts for SSRF protection when fetching document content.
var AllowedHosts = map[string]bool{
	"data.bka.gv.at":     true,
	"www.ris.bka.gv.at":  true,
	"ris.bka.gv.at":      true,
}
