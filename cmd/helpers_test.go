package cmd

import (
	"errors"
	"testing"
)

func TestValidationError_Error(t *testing.T) {
	err := &ValidationError{msg: "test error"}
	if err.Error() != "test error" {
		t.Errorf("expected %q, got %q", "test error", err.Error())
	}
}

func TestErrValidation_Simple(t *testing.T) {
	err := errValidation("Fehler: --app ist erforderlich")
	if err.Error() != "Fehler: --app ist erforderlich" {
		t.Errorf("unexpected message: %q", err.Error())
	}

	var ve *ValidationError
	if !errors.As(err, &ve) {
		t.Error("expected error to be a *ValidationError")
	}
}

func TestErrValidation_Formatted(t *testing.T) {
	err := errValidation("Fehler: ungültiger --app Wert %q", "invalid")
	expected := `Fehler: ungültiger --app Wert "invalid"`
	if err.Error() != expected {
		t.Errorf("expected %q, got %q", expected, err.Error())
	}
}

func TestErrValidation_IsNotGenericError(t *testing.T) {
	valErr := errValidation("validation issue")
	genericErr := errors.New("generic issue")

	var ve *ValidationError
	if !errors.As(valErr, &ve) {
		t.Error("errValidation should produce *ValidationError")
	}
	if errors.As(genericErr, &ve) {
		t.Error("generic error should not match *ValidationError")
	}
}
