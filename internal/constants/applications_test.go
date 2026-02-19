package constants

import "testing"

func TestIsValidHistoryApp(t *testing.T) {
	valid := []string{
		"bundesnormen", "landesnormen", "justiz", "vfgh", "vwgh",
		"bvwg", "lvwg", "bgblauth", "bgblalt", "bgblpdf",
		"lgblauth", "lgbl", "lgblno", "gemeinderecht", "gemeinderechtauth",
		"bvb", "vbl", "regv", "mrp", "erlaesse",
		"pruefgewo", "avsv", "spg", "kmger", "dsk",
		"gbk", "dok", "pvak", "normenliste", "asylgh",
	}
	for _, app := range valid {
		if !IsValidHistoryApp(app) {
			t.Errorf("IsValidHistoryApp(%q) = false, want true", app)
		}
	}

	invalid := []string{"", "unknown", "bundesrecht", "BUNDESNORMEN", "vfGH"}
	for _, app := range invalid {
		if IsValidHistoryApp(app) {
			t.Errorf("IsValidHistoryApp(%q) = true, want false", app)
		}
	}
}
