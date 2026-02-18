package format

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/philrox/ris-cli/internal/model"
)

func TestJSON_SearchResult(t *testing.T) {
	var buf bytes.Buffer
	result := model.SearchResult{
		TotalHits: 1,
		Page:      1,
		PageSize:  20,
		HasMore:   false,
		Documents: []model.Document{
			{
				Dokumentnummer: "NOR40052761",
				Titel:          "§ 1295 ABGB",
			},
		},
	}

	if err := JSON(&buf, result); err != nil {
		t.Fatal(err)
	}

	// Verify it's valid JSON.
	var parsed model.SearchResult
	if err := json.Unmarshal(buf.Bytes(), &parsed); err != nil {
		t.Fatalf("JSON output is not valid JSON: %v", err)
	}

	if parsed.TotalHits != 1 {
		t.Errorf("TotalHits = %d, want 1", parsed.TotalHits)
	}
	if parsed.Documents[0].Dokumentnummer != "NOR40052761" {
		t.Errorf("Dokumentnummer = %q, want NOR40052761", parsed.Documents[0].Dokumentnummer)
	}
}

func TestJSON_EmptyResult(t *testing.T) {
	var buf bytes.Buffer
	result := model.SearchResult{}

	if err := JSON(&buf, result); err != nil {
		t.Fatal(err)
	}

	var parsed model.SearchResult
	if err := json.Unmarshal(buf.Bytes(), &parsed); err != nil {
		t.Fatalf("JSON output is not valid JSON: %v", err)
	}
}

func TestJSONDocument(t *testing.T) {
	var buf bytes.Buffer
	doc := model.Document{
		Dokumentnummer: "NOR40052761",
		Titel:          "§ 1295 ABGB",
	}

	if err := JSONDocument(&buf, doc, "Document content here"); err != nil {
		t.Fatal(err)
	}

	var parsed model.DocumentContent
	if err := json.Unmarshal(buf.Bytes(), &parsed); err != nil {
		t.Fatalf("JSON output is not valid JSON: %v", err)
	}

	if parsed.Metadata.Dokumentnummer != "NOR40052761" {
		t.Errorf("Metadata.Dokumentnummer = %q, want NOR40052761", parsed.Metadata.Dokumentnummer)
	}
	if parsed.Content != "Document content here" {
		t.Errorf("Content = %q, want 'Document content here'", parsed.Content)
	}
}

func TestJSON_PrettyPrinted(t *testing.T) {
	var buf bytes.Buffer
	result := model.SearchResult{
		TotalHits: 1,
		Page:      1,
		Documents: []model.Document{
			{Titel: "Test"},
		},
	}

	if err := JSON(&buf, result); err != nil {
		t.Fatal(err)
	}

	// Pretty-printed JSON should contain indentation.
	out := buf.String()
	if !contains(out, "  ") {
		t.Error("JSON output should be pretty-printed with indentation")
	}
}

func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && bytes.Contains([]byte(s), []byte(substr))
}
