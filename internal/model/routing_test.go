package model

import "testing"

// TestJWRPrefixRoutesToVwgh verifies the bug fix: JWR documents are VwGH
// documents, not Justiz. The Applikation must be "Vwgh".
func TestJWRPrefixRoutesToVwgh(t *testing.T) {
	route, ok := matchPrefix("JWR_20230101_1234")
	if !ok {
		t.Fatal("expected JWR prefix to match, got no match")
	}
	if route.Applikation != "Vwgh" {
		t.Errorf("JWR Applikation = %q, want %q", route.Applikation, "Vwgh")
	}
	if route.URLPath != "Vwgh" {
		t.Errorf("JWR URLPath = %q, want %q", route.URLPath, "Vwgh")
	}
	if route.Endpoint != "Judikatur" {
		t.Errorf("JWR Endpoint = %q, want %q", route.Endpoint, "Judikatur")
	}
}

// TestJudikaturPrefixRouting verifies that all Judikatur prefixes route to
// their correct Applikation values.
func TestJudikaturPrefixRouting(t *testing.T) {
	tests := []struct {
		prefix      string
		wantURLPath string
		wantApp     string
	}{
		{"JWR", "Vwgh", "Vwgh"},
		{"JFR", "Vfgh", "Vfgh"},
		{"JFT", "Vfgh", "Vfgh"},
		{"JWT", "Justiz", "Justiz"},
		{"JJR", "Justiz", "Justiz"},
		{"DSB", "Dsk", "Dsk"},
		{"GBK", "Gbk", "Gbk"},
		{"ASYLGH", "AsylGH", "AsylGH"},
		{"BVWG", "Bvwg", "Bvwg"},
		{"LVWG", "Lvwg", "Lvwg"},
		{"PVAK", "Pvak", "Pvak"},
	}
	for _, tt := range tests {
		t.Run(tt.prefix, func(t *testing.T) {
			doc := tt.prefix + "_20230101_0001"
			route, ok := matchPrefix(doc)
			if !ok {
				t.Fatalf("expected prefix %q to match", tt.prefix)
			}
			if route.Endpoint != "Judikatur" {
				t.Errorf("Endpoint = %q, want %q", route.Endpoint, "Judikatur")
			}
			if route.URLPath != tt.wantURLPath {
				t.Errorf("URLPath = %q, want %q", route.URLPath, tt.wantURLPath)
			}
			if route.Applikation != tt.wantApp {
				t.Errorf("Applikation = %q, want %q", route.Applikation, tt.wantApp)
			}
		})
	}
}

// TestAllPrefixRoutes verifies every entry in the routing table maps to the
// expected DocumentRoute.
func TestAllPrefixRoutes(t *testing.T) {
	tests := []struct {
		prefix    string
		wantRoute DocumentRoute
	}{
		{"BGBLPDF", DocumentRoute{"BgblPdf", "Bundesrecht", "BgblPdf"}},
		{"BGBLA", DocumentRoute{"BgblAuth", "Bundesrecht", "BgblAuth"}},
		{"BGBL", DocumentRoute{"BgblAlt", "Bundesrecht", "BgblAlt"}},
		{"ASYLGH", DocumentRoute{"AsylGH", "Judikatur", "AsylGH"}},
		{"BVWG", DocumentRoute{"Bvwg", "Judikatur", "Bvwg"}},
		{"LVWG", DocumentRoute{"Lvwg", "Judikatur", "Lvwg"}},
		{"PVAK", DocumentRoute{"Pvak", "Judikatur", "Pvak"}},
		{"KMGER", DocumentRoute{"KmGer", "Sonstige", "KmGer"}},
		{"PRUEF", DocumentRoute{"PruefGewO", "Sonstige", "PruefGewO"}},
		{"REGV", DocumentRoute{"RegV", "Bundesrecht", "RegV"}},
		{"AVSV", DocumentRoute{"Avsv", "Sonstige", "Avsv"}},
		{"NOR", DocumentRoute{"Bundesnormen", "Bundesrecht", "BrKons"}},
		{"LBG", DocumentRoute{"LrBgld", "Landesrecht", "LrKons"}},
		{"LKT", DocumentRoute{"LrK", "Landesrecht", "LrKons"}},
		{"LNO", DocumentRoute{"LrNO", "Landesrecht", "LrKons"}},
		{"LOO", DocumentRoute{"LrOO", "Landesrecht", "LrKons"}},
		{"LSB", DocumentRoute{"LrSbg", "Landesrecht", "LrKons"}},
		{"LST", DocumentRoute{"LrStmk", "Landesrecht", "LrKons"}},
		{"LTI", DocumentRoute{"LrT", "Landesrecht", "LrKons"}},
		{"LVB", DocumentRoute{"LrVbg", "Landesrecht", "LrKons"}},
		{"LWI", DocumentRoute{"LrW", "Landesrecht", "LrKons"}},
		{"JWR", DocumentRoute{"Vwgh", "Judikatur", "Vwgh"}},
		{"JFR", DocumentRoute{"Vfgh", "Judikatur", "Vfgh"}},
		{"JFT", DocumentRoute{"Vfgh", "Judikatur", "Vfgh"}},
		{"JWT", DocumentRoute{"Justiz", "Judikatur", "Justiz"}},
		{"JJR", DocumentRoute{"Justiz", "Judikatur", "Justiz"}},
		{"DSB", DocumentRoute{"Dsk", "Judikatur", "Dsk"}},
		{"GBK", DocumentRoute{"Gbk", "Judikatur", "Gbk"}},
		{"BVB", DocumentRoute{"Bvb", "Bezirke", "Bvb"}},
		{"VBL", DocumentRoute{"Vbl", "Landesrecht", "Vbl"}},
		{"MRP", DocumentRoute{"Mrp", "Sonstige", "Mrp"}},
		{"ERL", DocumentRoute{"Erlaesse", "Sonstige", "Erlaesse"}},
		{"SPG", DocumentRoute{"Spg", "Sonstige", "Spg"}},
	}
	for _, tt := range tests {
		t.Run(tt.prefix, func(t *testing.T) {
			doc := tt.prefix + "_20230101_0001"
			route, ok := matchPrefix(doc)
			if !ok {
				t.Fatalf("expected prefix %q to match", tt.prefix)
			}
			if route != tt.wantRoute {
				t.Errorf("route = %+v, want %+v", route, tt.wantRoute)
			}
		})
	}
}

// TestMatchPrefixCaseInsensitive verifies that matchPrefix is case-insensitive.
func TestMatchPrefixCaseInsensitive(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantApp string
	}{
		{"uppercase", "NOR40026024", "BrKons"},
		{"lowercase", "nor40026024", "BrKons"},
		{"mixed case", "Nor40026024", "BrKons"},
		{"jwr lowercase", "jwr_20230101", "Vwgh"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			route, ok := matchPrefix(tt.input)
			if !ok {
				t.Fatalf("expected %q to match", tt.input)
			}
			if route.Applikation != tt.wantApp {
				t.Errorf("Applikation = %q, want %q", route.Applikation, tt.wantApp)
			}
		})
	}
}

// TestMatchPrefixLongestFirst verifies that longer prefixes take priority.
// "BGBLPDF" must match before "BGBLA" which must match before "BGBL".
func TestMatchPrefixLongestFirst(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantApp string
	}{
		{"BGBLPDF matches BgblPdf", "BGBLPDF_I_2023_42", "BgblPdf"},
		{"BGBLA matches BgblAuth", "BGBLA_2023_I_42", "BgblAuth"},
		{"BGBL matches BgblAlt", "BGBL_1990_42", "BgblAlt"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			route, ok := matchPrefix(tt.input)
			if !ok {
				t.Fatalf("expected %q to match", tt.input)
			}
			if route.Applikation != tt.wantApp {
				t.Errorf("Applikation = %q, want %q", route.Applikation, tt.wantApp)
			}
		})
	}
}

// TestMatchPrefixUnknown verifies that unknown prefixes return the default
// route with ok=false.
func TestMatchPrefixUnknown(t *testing.T) {
	tests := []string{"XYZ_123", "UNKNOWN", "", "ZZZ999"}
	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			route, ok := matchPrefix(input)
			if ok {
				t.Errorf("expected no match for %q, but got match", input)
			}
			if route.Endpoint != "Judikatur" {
				t.Errorf("default Endpoint = %q, want %q", route.Endpoint, "Judikatur")
			}
			if route.Applikation != "Justiz" {
				t.Errorf("default Applikation = %q, want %q", route.Applikation, "Justiz")
			}
		})
	}
}

// TestDirectURLFromPrefixKnown verifies that DirectURLFromPrefix constructs
// the correct URL for known prefixes.
func TestDirectURLFromPrefixKnown(t *testing.T) {
	tests := []struct {
		docNr   string
		wantURL string
	}{
		{
			"NOR40026024",
			"https://ris.bka.gv.at/Dokumente/Bundesnormen/NOR40026024/NOR40026024.html",
		},
		{
			"JWR_20230101_1234",
			"https://ris.bka.gv.at/Dokumente/Vwgh/JWR_20230101_1234/JWR_20230101_1234.html",
		},
		{
			"JFR_20230101_5678",
			"https://ris.bka.gv.at/Dokumente/Vfgh/JFR_20230101_5678/JFR_20230101_5678.html",
		},
		{
			"BGBLA_2023_I_42",
			"https://ris.bka.gv.at/Dokumente/BgblAuth/BGBLA_2023_I_42/BGBLA_2023_I_42.html",
		},
		{
			"BGBLPDF_I_2023_42",
			"https://ris.bka.gv.at/Dokumente/BgblPdf/BGBLPDF_I_2023_42/BGBLPDF_I_2023_42.html",
		},
		{
			"LBG40001234",
			"https://ris.bka.gv.at/Dokumente/LrBgld/LBG40001234/LBG40001234.html",
		},
		{
			"BVB_20230101_001",
			"https://ris.bka.gv.at/Dokumente/Bvb/BVB_20230101_001/BVB_20230101_001.html",
		},
		{
			"ERL_2023_001",
			"https://ris.bka.gv.at/Dokumente/Erlaesse/ERL_2023_001/ERL_2023_001.html",
		},
	}
	for _, tt := range tests {
		t.Run(tt.docNr, func(t *testing.T) {
			got := DirectURLFromPrefix(tt.docNr)
			if got != tt.wantURL {
				t.Errorf("DirectURLFromPrefix(%q)\n  got  %q\n  want %q", tt.docNr, got, tt.wantURL)
			}
		})
	}
}

// TestDirectURLFromPrefixUnknown verifies that DirectURLFromPrefix returns
// an empty string for document numbers with no matching prefix.
func TestDirectURLFromPrefixUnknown(t *testing.T) {
	tests := []string{"XYZ_123", "UNKNOWN_DOC", "", "999_NUMERIC"}
	for _, docNr := range tests {
		t.Run(docNr, func(t *testing.T) {
			got := DirectURLFromPrefix(docNr)
			if got != "" {
				t.Errorf("DirectURLFromPrefix(%q) = %q, want empty string", docNr, got)
			}
		})
	}
}

// TestSearchFallbackKnownPrefix verifies that SearchFallback returns the
// correct endpoint and applikation for known prefixes.
func TestSearchFallbackKnownPrefix(t *testing.T) {
	tests := []struct {
		docNr       string
		wantEndpt   string
		wantApp     string
	}{
		{"NOR40026024", "Bundesrecht", "BrKons"},
		{"JWR_20230101_1234", "Judikatur", "Vwgh"},
		{"JFR_20230101_5678", "Judikatur", "Vfgh"},
		{"JWT_20230101_9999", "Judikatur", "Justiz"},
		{"BGBLA_2023_I_42", "Bundesrecht", "BgblAuth"},
		{"LVWG_20230101_001", "Judikatur", "Lvwg"},
		{"LBG40001234", "Landesrecht", "LrKons"},
		{"KMGER_2023_001", "Sonstige", "KmGer"},
		{"BVB_20230101_001", "Bezirke", "Bvb"},
	}
	for _, tt := range tests {
		t.Run(tt.docNr, func(t *testing.T) {
			endpoint, applikation := SearchFallback(tt.docNr)
			if endpoint != tt.wantEndpt {
				t.Errorf("endpoint = %q, want %q", endpoint, tt.wantEndpt)
			}
			if applikation != tt.wantApp {
				t.Errorf("applikation = %q, want %q", applikation, tt.wantApp)
			}
		})
	}
}

// TestSearchFallbackUnknownPrefix verifies that SearchFallback returns the
// default route (Judikatur/Justiz) for unknown prefixes.
func TestSearchFallbackUnknownPrefix(t *testing.T) {
	tests := []string{"XYZ_123", "UNKNOWN_DOC", "", "999_NUMERIC"}
	for _, docNr := range tests {
		t.Run(docNr, func(t *testing.T) {
			endpoint, applikation := SearchFallback(docNr)
			if endpoint != "Judikatur" {
				t.Errorf("endpoint = %q, want %q", endpoint, "Judikatur")
			}
			if applikation != "Justiz" {
				t.Errorf("applikation = %q, want %q", applikation, "Justiz")
			}
		})
	}
}

// TestPrefixRoutesOrdering verifies that the routing table has longer prefixes
// before shorter ones with the same leading characters. This is critical for
// correct matching (e.g., BGBLPDF before BGBLA before BGBL).
func TestPrefixRoutesOrdering(t *testing.T) {
	// Build a map of prefixes that share leading characters. For each group,
	// the longer prefix must appear before any shorter one in the slice.
	groups := map[string][]int{} // prefix -> list of indices
	for i, entry := range prefixRoutes {
		for j, other := range prefixRoutes {
			if i != j && len(entry.Prefix) < len(other.Prefix) &&
				other.Prefix[:len(entry.Prefix)] == entry.Prefix {
				// other.Prefix is a longer version of entry.Prefix.
				// other must appear before entry (j < i).
				key := entry.Prefix + "/" + other.Prefix
				if _, exists := groups[key]; !exists {
					groups[key] = []int{j, i}
				}
			}
		}
	}
	for key, indices := range groups {
		longIdx, shortIdx := indices[0], indices[1]
		if longIdx >= shortIdx {
			t.Errorf("ordering violation for %s: longer prefix at index %d, shorter at %d",
				key, longIdx, shortIdx)
		}
	}
}

// TestDirectURLFromPrefixURLFormat verifies the URL structure is correct
// (scheme, host, path pattern).
func TestDirectURLFromPrefixURLFormat(t *testing.T) {
	url := DirectURLFromPrefix("NOR40026024")
	const prefix = "https://ris.bka.gv.at/Dokumente/"
	if len(url) < len(prefix) || url[:len(prefix)] != prefix {
		t.Errorf("URL does not start with %q: got %q", prefix, url)
	}
	const suffix = ".html"
	if len(url) < len(suffix) || url[len(url)-len(suffix):] != suffix {
		t.Errorf("URL does not end with %q: got %q", suffix, url)
	}
}
