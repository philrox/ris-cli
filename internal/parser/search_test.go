package parser

import (
	"encoding/json"
	"testing"
)

// buildBundesrechtResponse constructs a minimal rawResponse JSON with
// a single Bundesrecht document whose BrKons.Ausserkrafttretensdatum is set
// to the given value.
func buildBundesrechtResponse(ausserkrafttreten string) []byte {
	resp := rawResponse{
		OgdSearchResult: rawSearchResult{
			OgdDocumentResults: rawDocumentResults{
				Hits: json.RawMessage(`1`),
				Docs: FlexibleArray[rawDocumentReference]{
					{
						Data: rawData{
							Metadaten: rawMetadaten{
								Technisch: rawTechnisch{
									ID:          "TEST001",
									Applikation: "BrKons",
								},
								Allgemein: rawAllgemein{
									DokumentURL: "https://example.com/doc",
								},
							},
						},
					},
				},
			},
		},
	}

	// Build the Bundesrecht metadata separately so we can control the
	// Ausserkrafttretensdatum value precisely.
	br := rawBundesrecht{
		Kurztitel: "TestGesetz",
		Langtitel: "Testgesetz Langform",
		Titel:     FlexibleString("Test Titel"),
		Eli:       "eli/test/2025",
		BrKons: &rawSubApp{
			Kundmachungsorgan:       "BGBl. I Nr. 1/2020",
			Inkrafttretensdatum:     "2020-01-01",
			Ausserkrafttretensdatum: ausserkrafttreten,
		},
	}
	brJSON, _ := json.Marshal(br)
	resp.OgdSearchResult.OgdDocumentResults.Docs[0].Data.Metadaten.Bundesrecht = brJSON

	data, _ := json.Marshal(resp)
	return data
}

func TestNoExpiryDate_SentinelFiltered(t *testing.T) {
	data := buildBundesrechtResponse("9999-12-31")
	result, err := ParseSearchResponse(data)
	if err != nil {
		t.Fatalf("ParseSearchResponse returned error: %v", err)
	}
	if len(result.Documents) != 1 {
		t.Fatalf("expected 1 document, got %d", len(result.Documents))
	}
	doc := result.Documents[0]
	if doc.Citation == nil {
		t.Fatal("expected Citation to be non-nil")
	}
	if doc.Citation.Ausserkrafttreten != nil {
		t.Errorf("expected Ausserkrafttreten to be nil for sentinel date 9999-12-31, got %q", *doc.Citation.Ausserkrafttreten)
	}
}

func TestNoExpiryDate_RealDatePreserved(t *testing.T) {
	data := buildBundesrechtResponse("2025-01-01")
	result, err := ParseSearchResponse(data)
	if err != nil {
		t.Fatalf("ParseSearchResponse returned error: %v", err)
	}
	if len(result.Documents) != 1 {
		t.Fatalf("expected 1 document, got %d", len(result.Documents))
	}
	doc := result.Documents[0]
	if doc.Citation == nil {
		t.Fatal("expected Citation to be non-nil")
	}
	if doc.Citation.Ausserkrafttreten == nil {
		t.Fatal("expected Ausserkrafttreten to be non-nil for real date 2025-01-01")
	}
	if *doc.Citation.Ausserkrafttreten != "2025-01-01" {
		t.Errorf("expected Ausserkrafttreten = %q, got %q", "2025-01-01", *doc.Citation.Ausserkrafttreten)
	}
}

func TestLandesrecht_ParsesWithUnifiedSubApp(t *testing.T) {
	resp := rawResponse{
		OgdSearchResult: rawSearchResult{
			OgdDocumentResults: rawDocumentResults{
				Hits: json.RawMessage(`1`),
				Docs: FlexibleArray[rawDocumentReference]{
					{
						Data: rawData{
							Metadaten: rawMetadaten{
								Technisch: rawTechnisch{
									ID:          "LR001",
									Applikation: "LrKons",
								},
								Allgemein: rawAllgemein{
									DokumentURL: "https://example.com/lr",
								},
							},
						},
					},
				},
			},
		},
	}

	lr := rawLandesrecht{
		Kurztitel: "TestLandesgesetz",
		Langtitel: "Testlandesgesetz Langform",
		Titel:     FlexibleString("LR Titel"),
		Eli:       "eli/lr/2025",
		LrKons: &rawSubApp{
			Kundmachungsorgan:       "LGBl. Nr. 1/2020",
			Inkrafttretensdatum:     "2020-06-01",
			Ausserkrafttretensdatum: "2030-12-31",
		},
	}
	lrJSON, _ := json.Marshal(lr)
	resp.OgdSearchResult.OgdDocumentResults.Docs[0].Data.Metadaten.Landesrecht = lrJSON

	data, _ := json.Marshal(resp)
	result, err := ParseSearchResponse(data)
	if err != nil {
		t.Fatalf("ParseSearchResponse returned error: %v", err)
	}
	if len(result.Documents) != 1 {
		t.Fatalf("expected 1 document, got %d", len(result.Documents))
	}
	doc := result.Documents[0]
	if doc.Kurztitel != "TestLandesgesetz" {
		t.Errorf("expected Kurztitel %q, got %q", "TestLandesgesetz", doc.Kurztitel)
	}
	if doc.Citation == nil {
		t.Fatal("expected Citation to be non-nil")
	}
	if doc.Citation.Kundmachungsorgan != "LGBl. Nr. 1/2020" {
		t.Errorf("expected Kundmachungsorgan %q, got %q", "LGBl. Nr. 1/2020", doc.Citation.Kundmachungsorgan)
	}
	if doc.Citation.Inkrafttreten != "2020-06-01" {
		t.Errorf("expected Inkrafttreten %q, got %q", "2020-06-01", doc.Citation.Inkrafttreten)
	}
	if doc.Citation.Ausserkrafttreten == nil {
		t.Fatal("expected Ausserkrafttreten to be non-nil")
	}
	if *doc.Citation.Ausserkrafttreten != "2030-12-31" {
		t.Errorf("expected Ausserkrafttreten %q, got %q", "2030-12-31", *doc.Citation.Ausserkrafttreten)
	}
}

func TestNoExpiryDate_EmptyStringFiltered(t *testing.T) {
	data := buildBundesrechtResponse("")
	result, err := ParseSearchResponse(data)
	if err != nil {
		t.Fatalf("ParseSearchResponse returned error: %v", err)
	}
	if len(result.Documents) != 1 {
		t.Fatalf("expected 1 document, got %d", len(result.Documents))
	}
	doc := result.Documents[0]
	if doc.Citation == nil {
		t.Fatal("expected Citation to be non-nil")
	}
	if doc.Citation.Ausserkrafttreten != nil {
		t.Errorf("expected Ausserkrafttreten to be nil for empty string, got %q", *doc.Citation.Ausserkrafttreten)
	}
}
