package format

import (
	"testing"

	"github.com/fatih/color"
	"github.com/philrox/ris-cli/internal/model"
)

func init() {
	// Disable colors in tests for predictable output.
	color.NoColor = true
}

func TestFormatCitation_Nil(t *testing.T) {
	got := FormatCitation(nil)
	if got != "" {
		t.Errorf("FormatCitation(nil) = %q, want empty", got)
	}
}

func TestFormatCitation_ParagraphAndKurztitel(t *testing.T) {
	c := &model.Citation{
		Kurztitel:         "ABGB",
		Paragraph:         "§ 1295",
		Kundmachungsorgan: "JGS Nr. 946/1811",
	}
	got := FormatCitation(c)
	want := "§ 1295 ABGB (JGS Nr. 946/1811)"
	if got != want {
		t.Errorf("FormatCitation() = %q, want %q", got, want)
	}
}

func TestFormatCitation_KurztitelOnly(t *testing.T) {
	c := &model.Citation{Kurztitel: "StGB"}
	got := FormatCitation(c)
	if got != "StGB" {
		t.Errorf("FormatCitation() = %q, want %q", got, "StGB")
	}
}

func TestFormatCitation_LangtitelFallback(t *testing.T) {
	c := &model.Citation{Langtitel: "Allgemeines buergerliches Gesetzbuch"}
	got := FormatCitation(c)
	if got != "Allgemeines buergerliches Gesetzbuch" {
		t.Errorf("FormatCitation() = %q, want Langtitel fallback", got)
	}
}

func TestFormatCitation_Geschaeftszahl(t *testing.T) {
	c := &model.Citation{Geschaeftszahl: "5Ob234/20b"}
	got := FormatCitation(c)
	if got != "5Ob234/20b" {
		t.Errorf("FormatCitation() = %q, want %q", got, "5Ob234/20b")
	}
}

func TestFormatCitation_Entscheidungsdatum(t *testing.T) {
	c := &model.Citation{
		Kurztitel:          "VfGH",
		Entscheidungsdatum: "2024-01-15",
	}
	got := FormatCitation(c)
	want := "VfGH vom 2024-01-15"
	if got != want {
		t.Errorf("FormatCitation() = %q, want %q", got, want)
	}
}

func TestFormatCitation_GeschaeftszahlWithKurztitel(t *testing.T) {
	c := &model.Citation{Kurztitel: "VfGH", Geschaeftszahl: "G123/24"}
	got := FormatCitation(c)
	want := "VfGH G123/24"
	if got != want {
		t.Errorf("FormatCitation() = %q, want %q", got, want)
	}
}

func TestFormatCitation_Empty(t *testing.T) {
	c := &model.Citation{}
	got := FormatCitation(c)
	if got != "" {
		t.Errorf("FormatCitation(empty) = %q, want empty", got)
	}
}

func TestFormatDates_Nil(t *testing.T) {
	got := FormatDates(nil)
	if got != "" {
		t.Errorf("FormatDates(nil) = %q, want empty", got)
	}
}

func TestFormatDates_Inkrafttreten(t *testing.T) {
	c := &model.Citation{Inkrafttreten: "1812-01-01"}
	got := FormatDates(c)
	want := "in Kraft seit 1812-01-01"
	if got != want {
		t.Errorf("FormatDates() = %q, want %q", got, want)
	}
}

func TestFormatDates_Both(t *testing.T) {
	akt := "2024-12-31"
	c := &model.Citation{
		Inkrafttreten:     "2020-01-01",
		Ausserkrafttreten: &akt,
	}
	got := FormatDates(c)
	want := "in Kraft seit 2020-01-01, außer Kraft seit 2024-12-31"
	if got != want {
		t.Errorf("FormatDates() = %q, want %q", got, want)
	}
}

func TestFormatDates_NilAusserkrafttreten(t *testing.T) {
	c := &model.Citation{
		Inkrafttreten:     "2020-01-01",
		Ausserkrafttreten: nil,
	}
	got := FormatDates(c)
	want := "in Kraft seit 2020-01-01"
	if got != want {
		t.Errorf("FormatDates() = %q, want %q", got, want)
	}
}

func TestFormatDates_Empty(t *testing.T) {
	c := &model.Citation{}
	got := FormatDates(c)
	if got != "" {
		t.Errorf("FormatDates(empty) = %q, want empty", got)
	}
}
