package model

import "strings"

// DocumentRoute describes how to construct a direct URL or search fallback
// for a document number based on its prefix.
type DocumentRoute struct {
	URLPath     string // Path segment in the direct URL (e.g., "Bundesnormen")
	Endpoint    string // API endpoint for search fallback (e.g., "Bundesrecht")
	Applikation string // API Applikation value for search fallback
}

// prefixRoutes maps document number prefixes to their routing information.
// Entries are ordered by prefix length (longest first) for correct matching.
var prefixRoutes = []struct {
	Prefix string
	Route  DocumentRoute
}{
	// Longest prefixes first to ensure correct matching
	{"BGBLPDF", DocumentRoute{URLPath: "BgblPdf", Endpoint: "Bundesrecht", Applikation: "BgblPdf"}},
	{"BGBLA", DocumentRoute{URLPath: "BgblAuth", Endpoint: "Bundesrecht", Applikation: "BgblAuth"}},
	{"BGBL", DocumentRoute{URLPath: "BgblAlt", Endpoint: "Bundesrecht", Applikation: "BgblAlt"}},
	{"ASYLGH", DocumentRoute{URLPath: "AsylGH", Endpoint: "Judikatur", Applikation: "AsylGH"}},
	{"BVWG", DocumentRoute{URLPath: "Bvwg", Endpoint: "Judikatur", Applikation: "Bvwg"}},
	{"LVWG", DocumentRoute{URLPath: "Lvwg", Endpoint: "Judikatur", Applikation: "Lvwg"}},
	{"PVAK", DocumentRoute{URLPath: "Pvak", Endpoint: "Judikatur", Applikation: "Pvak"}},
	{"KMGER", DocumentRoute{URLPath: "KmGer", Endpoint: "Sonstige", Applikation: "KmGer"}},
	{"PRUEF", DocumentRoute{URLPath: "PruefGewO", Endpoint: "Sonstige", Applikation: "PruefGewO"}},
	{"REGV", DocumentRoute{URLPath: "RegV", Endpoint: "Bundesrecht", Applikation: "RegV"}},
	{"AVSV", DocumentRoute{URLPath: "Avsv", Endpoint: "Sonstige", Applikation: "Avsv"}},
	{"NOR", DocumentRoute{URLPath: "Bundesnormen", Endpoint: "Bundesrecht", Applikation: "BrKons"}},
	{"LBG", DocumentRoute{URLPath: "LrBgld", Endpoint: "Landesrecht", Applikation: "LrKons"}},
	{"LKT", DocumentRoute{URLPath: "LrK", Endpoint: "Landesrecht", Applikation: "LrKons"}},
	{"LNO", DocumentRoute{URLPath: "LrNO", Endpoint: "Landesrecht", Applikation: "LrKons"}},
	{"LOO", DocumentRoute{URLPath: "LrOO", Endpoint: "Landesrecht", Applikation: "LrKons"}},
	{"LSB", DocumentRoute{URLPath: "LrSbg", Endpoint: "Landesrecht", Applikation: "LrKons"}},
	{"LST", DocumentRoute{URLPath: "LrStmk", Endpoint: "Landesrecht", Applikation: "LrKons"}},
	{"LTI", DocumentRoute{URLPath: "LrT", Endpoint: "Landesrecht", Applikation: "LrKons"}},
	{"LVB", DocumentRoute{URLPath: "LrVbg", Endpoint: "Landesrecht", Applikation: "LrKons"}},
	{"LWI", DocumentRoute{URLPath: "LrW", Endpoint: "Landesrecht", Applikation: "LrKons"}},
	{"JWR", DocumentRoute{URLPath: "Vwgh", Endpoint: "Judikatur", Applikation: "Vwgh"}},
	{"JFR", DocumentRoute{URLPath: "Vfgh", Endpoint: "Judikatur", Applikation: "Vfgh"}},
	{"JFT", DocumentRoute{URLPath: "Vfgh", Endpoint: "Judikatur", Applikation: "Vfgh"}},
	{"JWT", DocumentRoute{URLPath: "Justiz", Endpoint: "Judikatur", Applikation: "Justiz"}},
	{"JJR", DocumentRoute{URLPath: "Justiz", Endpoint: "Judikatur", Applikation: "Justiz"}},
	{"DSB", DocumentRoute{URLPath: "Dsk", Endpoint: "Judikatur", Applikation: "Dsk"}},
	{"GBK", DocumentRoute{URLPath: "Gbk", Endpoint: "Judikatur", Applikation: "Gbk"}},
	{"BVB", DocumentRoute{URLPath: "Bvb", Endpoint: "Bezirke", Applikation: "Bvb"}},
	{"VBL", DocumentRoute{URLPath: "Vbl", Endpoint: "Landesrecht", Applikation: "Vbl"}},
	{"MRP", DocumentRoute{URLPath: "Mrp", Endpoint: "Sonstige", Applikation: "Mrp"}},
	{"ERL", DocumentRoute{URLPath: "Erlaesse", Endpoint: "Sonstige", Applikation: "Erlaesse"}},
	{"SPG", DocumentRoute{URLPath: "Spg", Endpoint: "Sonstige", Applikation: "Spg"}},
}

// defaultRoute is used when no prefix matches.
var defaultRoute = DocumentRoute{
	Endpoint:    "Judikatur",
	Applikation: "Justiz",
}

// matchPrefix finds the routing information for a document number.
func matchPrefix(dokumentnummer string) (DocumentRoute, bool) {
	upper := strings.ToUpper(dokumentnummer)
	for _, entry := range prefixRoutes {
		if strings.HasPrefix(upper, entry.Prefix) {
			return entry.Route, true
		}
	}
	return defaultRoute, false
}

// DirectURLFromPrefix constructs a direct document URL from a document number
// using the prefix routing table. Returns empty string if no match found.
func DirectURLFromPrefix(dokumentnummer string) string {
	route, ok := matchPrefix(dokumentnummer)
	if !ok || route.URLPath == "" {
		return ""
	}
	return "https://ris.bka.gv.at/Dokumente/" + route.URLPath + "/" + dokumentnummer + "/" + dokumentnummer + ".html"
}

// SearchFallback returns the endpoint and applikation for search-based document
// retrieval when direct URL construction fails.
func SearchFallback(dokumentnummer string) (endpoint, applikation string) {
	route, _ := matchPrefix(dokumentnummer)
	return route.Endpoint, route.Applikation
}
