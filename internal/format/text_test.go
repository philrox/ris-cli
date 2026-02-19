package format

import (
	"bytes"
	"strings"
	"testing"

	"github.com/philrox/ris-cli/internal/model"
)

func TestText_EmptyResults(t *testing.T) {
	var buf bytes.Buffer
	result := model.SearchResult{}
	if err := Text(&buf, result); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "Keine Ergebnisse") {
		t.Error("expected 'Keine Ergebnisse' message for empty results")
	}
}

func TestText_SingleResult(t *testing.T) {
	var buf bytes.Buffer
	result := model.SearchResult{
		TotalHits: 1,
		Page:      1,
		PageSize:  20,
		Documents: []model.Document{
			{
				Dokumentnummer: "NOR40052761",
				Titel:          "§ 1295 ABGB",
				Citation: &model.Citation{
					Kurztitel:         "ABGB",
					Paragraph:         "§ 1295",
					Kundmachungsorgan: "JGS Nr. 946/1811",
					Inkrafttreten:     "1812-01-01",
				},
			},
		},
	}

	if err := Text(&buf, result); err != nil {
		t.Fatal(err)
	}

	out := buf.String()

	checks := []string{
		"Ergebnisse: 1 gesamt",
		"§ 1295 ABGB",
		"NOR40052761",
		"JGS Nr. 946/1811",
		"in Kraft seit 1812-01-01",
	}

	for _, check := range checks {
		if !strings.Contains(out, check) {
			t.Errorf("output missing %q", check)
		}
	}
}

func TestText_HasMore(t *testing.T) {
	var buf bytes.Buffer
	result := model.SearchResult{
		TotalHits: 42,
		Page:      1,
		PageSize:  20,
		HasMore:   true,
		Documents: []model.Document{
			{Titel: "Test"},
		},
	}

	if err := Text(&buf, result); err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(buf.String(), "--page 2") {
		t.Error("expected pagination hint for HasMore=true")
	}
}

func TestText_NoHasMore(t *testing.T) {
	var buf bytes.Buffer
	result := model.SearchResult{
		TotalHits: 1,
		Page:      1,
		PageSize:  20,
		HasMore:   false,
		Documents: []model.Document{
			{Titel: "Test"},
		},
	}

	if err := Text(&buf, result); err != nil {
		t.Fatal(err)
	}

	if strings.Contains(buf.String(), "--page") {
		t.Error("should not show pagination hint when HasMore=false")
	}
}

func TestText_ELIDisplay(t *testing.T) {
	var buf bytes.Buffer
	result := model.SearchResult{
		TotalHits: 1,
		Page:      1,
		Documents: []model.Document{
			{
				Titel: "Test",
				Citation: &model.Citation{
					Kurztitel: "ABGB",
					Eli:       "eli/bgbl/1811/946/main",
				},
			},
		},
	}

	if err := Text(&buf, result); err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(buf.String(), "ELI:") {
		t.Error("expected ELI display in output")
	}
}

func TestText_Leitsatz(t *testing.T) {
	var buf bytes.Buffer
	result := model.SearchResult{
		TotalHits: 1,
		Page:      1,
		Documents: []model.Document{
			{
				Titel:    "Test",
				Leitsatz: "Dies ist ein Leitsatz.",
			},
		},
	}

	if err := Text(&buf, result); err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(buf.String(), "Leitsatz:") {
		t.Error("expected Leitsatz in output")
	}
}

func TestText_LeitsatzTruncation(t *testing.T) {
	var buf bytes.Buffer
	longText := strings.Repeat("A", 250)
	result := model.SearchResult{
		TotalHits: 1,
		Page:      1,
		Documents: []model.Document{
			{
				Titel:    "Test",
				Leitsatz: longText,
			},
		},
	}

	if err := Text(&buf, result); err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(buf.String(), "...") {
		t.Error("expected truncated Leitsatz with ellipsis")
	}
}

func TestTextDocument_Basic(t *testing.T) {
	var buf bytes.Buffer
	doc := model.Document{
		Dokumentnummer: "NOR40052761",
		Titel:          "§ 1295 ABGB",
		Citation: &model.Citation{
			Kurztitel:         "ABGB",
			Paragraph:         "§ 1295",
			Kundmachungsorgan: "JGS Nr. 946/1811",
			Inkrafttreten:     "1812-01-01",
			Eli:               "eli/bgbl/1811/946/main",
		},
	}

	if err := TextDocument(&buf, doc, "Wer einem andern durch Verschulden..."); err != nil {
		t.Fatal(err)
	}

	out := buf.String()

	checks := []string{
		"§ 1295 ABGB",
		"NOR40052761",
		"in Kraft seit 1812-01-01",
		"ELI:",
		"Wer einem andern durch Verschulden",
	}

	for _, check := range checks {
		if !strings.Contains(out, check) {
			t.Errorf("TextDocument output missing %q", check)
		}
	}
}

func TestTextDocument_EmptyContent(t *testing.T) {
	var buf bytes.Buffer
	doc := model.Document{Titel: "Test"}

	if err := TextDocument(&buf, doc, ""); err != nil {
		t.Fatal(err)
	}

	// Should not contain separator when no content.
	if strings.Count(buf.String(), "─") > 0 {
		t.Error("should not show content separator when content is empty")
	}
}

func TestText_LeitsatzTruncationBoundary(t *testing.T) {
	tests := []struct {
		name       string
		length     int
		wantEllipsis bool
	}{
		{"exactly 200 chars - no truncation", 200, false},
		{"201 chars - truncated", 201, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			leitsatz := strings.Repeat("X", tt.length)
			result := model.SearchResult{
				TotalHits: 1,
				Page:      1,
				Documents: []model.Document{
					{
						Titel:    "Test",
						Leitsatz: leitsatz,
					},
				},
			}

			if err := Text(&buf, result); err != nil {
				t.Fatal(err)
			}

			out := buf.String()
			hasEllipsis := strings.Contains(out, "...")

			if hasEllipsis != tt.wantEllipsis {
				t.Errorf("len=%d: ellipsis present = %v, want %v", tt.length, hasEllipsis, tt.wantEllipsis)
			}

			if tt.wantEllipsis {
				// When truncated, the output should contain exactly 200 X's followed by "..."
				truncated := strings.Repeat("X", 200) + "..."
				if !strings.Contains(out, truncated) {
					t.Errorf("len=%d: expected truncated output to contain 200 chars + ellipsis", tt.length)
				}
			} else {
				// When not truncated, all original chars should be present
				if !strings.Contains(out, leitsatz) {
					t.Errorf("len=%d: expected full Leitsatz in output", tt.length)
				}
			}
		})
	}
}

func TestDocTitle_Fallbacks(t *testing.T) {
	tests := []struct {
		name string
		doc  model.Document
		want string
	}{
		{"Titel", model.Document{Titel: "foo"}, "foo"},
		{"Kurztitel", model.Document{Kurztitel: "bar"}, "bar"},
		{"Dokumentnummer", model.Document{Dokumentnummer: "NOR123"}, "NOR123"},
		{"Empty", model.Document{}, "(ohne Titel)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := docTitle(tt.doc)
			if got != tt.want {
				t.Errorf("docTitle() = %q, want %q", got, tt.want)
			}
		})
	}
}
