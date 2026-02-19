package cmd

import (
	"errors"
	"testing"
)

// executeCommand runs a cobra command with the given args and returns the error.
func executeCommand(args ...string) error {
	rootCmd.SetArgs(args)
	return rootCmd.Execute()
}

// assertValidationError checks that the error is a *ValidationError and contains substr.
func assertValidationError(t *testing.T, err error, substr string) {
	t.Helper()
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	var ve *ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("expected *ValidationError, got %T: %v", err, err)
	}
	if substr != "" && !containsSubstr(err.Error(), substr) {
		t.Errorf("expected error to contain %q, got %q", substr, err.Error())
	}
}

func containsSubstr(s, substr string) bool {
	return len(s) >= len(substr) && (substr == "" || findSubstr(s, substr))
}

func findSubstr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestBundesrecht_NoArgs_ReturnsValidationError(t *testing.T) {
	err := executeCommand("bundesrecht")
	assertValidationError(t, err, "mindestens --search, --title oder --paragraph erforderlich")
}

func TestBundesrecht_InvalidApp_ReturnsValidationError(t *testing.T) {
	err := executeCommand("bundesrecht", "--search", "test", "--app", "invalid")
	assertValidationError(t, err, "ungültiger --app Wert")
}

func TestLandesrecht_NoArgs_ReturnsValidationError(t *testing.T) {
	err := executeCommand("landesrecht")
	assertValidationError(t, err, "mindestens --search, --title oder --state erforderlich")
}

func TestLandesrecht_InvalidState_ReturnsValidationError(t *testing.T) {
	err := executeCommand("landesrecht", "--search", "test", "--state", "invalid")
	assertValidationError(t, err, "ungültiger --state Wert")
}

func TestJudikatur_NoArgs_ReturnsValidationError(t *testing.T) {
	err := executeCommand("judikatur")
	assertValidationError(t, err, "mindestens --search, --norm oder --case-number erforderlich")
}

func TestJudikatur_InvalidCourt_ReturnsValidationError(t *testing.T) {
	err := executeCommand("judikatur", "--search", "test", "--court", "invalid")
	assertValidationError(t, err, "ungültiger --court Wert")
}

func TestBgbl_NoArgs_ReturnsValidationError(t *testing.T) {
	err := executeCommand("bgbl")
	assertValidationError(t, err, "mindestens --number, --year, --search oder --title erforderlich")
}

func TestBgbl_InvalidApp_ReturnsValidationError(t *testing.T) {
	err := executeCommand("bgbl", "--search", "test", "--app", "invalid")
	assertValidationError(t, err, "ungültiger --app Wert")
}

func TestBgbl_InvalidPart_ReturnsValidationError(t *testing.T) {
	err := executeCommand("bgbl", "--search", "test", "--app", "bgblauth", "--part", "99")
	assertValidationError(t, err, "ungültiger --part Wert")
}

func TestLgbl_NoArgs_ReturnsValidationError(t *testing.T) {
	err := executeCommand("lgbl")
	assertValidationError(t, err, "mindestens --number, --year, --state, --search oder --title erforderlich")
}

func TestHistory_NoApp_ReturnsValidationError(t *testing.T) {
	err := executeCommand("history")
	assertValidationError(t, err, "--app ist erforderlich")
}

func TestHistory_NoDateRange_ReturnsValidationError(t *testing.T) {
	err := executeCommand("history", "--app", "bundesnormen")
	assertValidationError(t, err, "mindestens --from oder --to erforderlich")
}

func TestHistory_InvalidApp_ReturnsValidationError(t *testing.T) {
	err := executeCommand("history", "--app", "invalid", "--from", "2024-01-01")
	assertValidationError(t, err, "ungültiger --app Wert")
}

func TestDokument_NoArgs_ReturnsValidationError(t *testing.T) {
	err := executeCommand("dokument")
	assertValidationError(t, err, "Dokumentnummer oder --url erforderlich")
}

func TestVerordnungen_NoArgs_ReturnsValidationError(t *testing.T) {
	err := executeCommand("verordnungen")
	assertValidationError(t, err, "mindestens --search, --title, --state, --number oder --from erforderlich")
}

func TestGemeinden_NoArgs_ReturnsValidationError(t *testing.T) {
	err := executeCommand("gemeinden")
	assertValidationError(t, err, "mindestens ein Suchparameter erforderlich")
}

func TestGemeinden_InvalidApp_ReturnsValidationError(t *testing.T) {
	err := executeCommand("gemeinden", "--search", "test", "--app", "invalid")
	assertValidationError(t, err, "ungültiger --app Wert")
}

func TestBezirke_NoArgs_ReturnsValidationError(t *testing.T) {
	err := executeCommand("bezirke")
	assertValidationError(t, err, "mindestens --search, --title, --state, --authority oder --number erforderlich")
}

func TestRegvorl_NoArgs_ReturnsValidationError(t *testing.T) {
	err := executeCommand("regvorl")
	assertValidationError(t, err, "mindestens --search, --title, --from, --ministry oder --since erforderlich")
}

func TestRegvorl_InvalidMinistry_ReturnsValidationError(t *testing.T) {
	err := executeCommand("regvorl", "--ministry", "invalid", "--search", "test")
	assertValidationError(t, err, "ungültiger --ministry Wert")
}

func TestSonstige_InvalidSince_ReturnsValidationError(t *testing.T) {
	err := executeCommand("sonstige", "mrp", "--search", "test", "--since", "invalid")
	assertValidationError(t, err, "ungültiger --since Wert")
}

func TestSonstige_InvalidSortDir_ReturnsValidationError(t *testing.T) {
	// Use erlaesse (not mrp) to avoid Cobra flag state leaking from the --since test above.
	err := executeCommand("sonstige", "erlaesse", "--search", "test", "--sort-dir", "invalid")
	assertValidationError(t, err, "ungültiger --sort-dir Wert")
}
