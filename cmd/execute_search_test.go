package cmd

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/philrox/risgo/internal/api"
	"github.com/spf13/cobra"
)

// minimalAPIResponse is a minimal valid RIS API JSON response with zero results.
const minimalAPIResponse = `{
	"OgdSearchResult": {
		"OgdDocumentResults": {
			"Hits": "0",
			"OgdDocumentReference": []
		}
	}
}`

// oneHitAPIResponse is a valid RIS API JSON response with one document.
const oneHitAPIResponse = `{
	"OgdSearchResult": {
		"OgdDocumentResults": {
			"Hits": {
				"#text": "1",
				"@pageNumber": "1",
				"@pageSize": "20"
			},
			"OgdDocumentReference": {
				"Data": {
					"Metadaten": {
						"Technisch": {
							"ID": "test-id-1",
							"Applikation": "BrKons"
						},
						"Allgemein": {
							"DokumentUrl": "https://example.com/doc1",
							"Dokumenttyp": "BVG"
						}
					},
					"Dokumentliste": null
				}
			}
		}
	}
}`

// setupTestCmd creates a root cobra command with the required global flags
// and overrides RIS_BASE_URL to point to the given test server.
func setupTestCmd(baseURL string) *cobra.Command {
	cmd := &cobra.Command{Use: "test"}
	cmd.PersistentFlags().Bool("json", false, "")
	cmd.PersistentFlags().Bool("verbose", false, "")
	cmd.PersistentFlags().Duration("timeout", 0, "")
	cmd.PersistentFlags().Int("page", 1, "")
	cmd.PersistentFlags().Int("limit", 20, "")

	os.Setenv("RIS_BASE_URL", baseURL)
	return cmd
}

func TestExecuteSearch_Success_TextOutput(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify endpoint path.
		if !strings.HasSuffix(r.URL.Path, "/Bundesrecht") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		// Verify params are passed through.
		if r.URL.Query().Get("Suchworte") != "test" {
			t.Errorf("expected Suchworte=test, got %q", r.URL.Query().Get("Suchworte"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(minimalAPIResponse))
	}))
	defer srv.Close()

	cmd := setupTestCmd(srv.URL)
	defer os.Unsetenv("RIS_BASE_URL")

	params := api.NewParams()
	params.Set("Suchworte", "test")

	err := executeSearch(cmd, "Bundesrecht", "Suche...", params)
	if err != nil {
		t.Fatalf("executeSearch returned error: %v", err)
	}
}

func TestExecuteSearch_Success_JSONOutput(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(oneHitAPIResponse))
	}))
	defer srv.Close()

	cmd := setupTestCmd(srv.URL)
	defer os.Unsetenv("RIS_BASE_URL")

	// Enable JSON mode.
	cmd.PersistentFlags().Set("json", "true")

	params := api.NewParams()
	params.Set("Suchworte", "test")

	err := executeSearch(cmd, "Bundesrecht", "Suche...", params)
	if err != nil {
		t.Fatalf("executeSearch returned error: %v", err)
	}
}

func TestExecuteSearch_APIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	cmd := setupTestCmd(srv.URL)
	defer os.Unsetenv("RIS_BASE_URL")

	params := api.NewParams()
	params.Set("Suchworte", "test")

	err := executeSearch(cmd, "Bundesrecht", "Suche...", params)
	if err == nil {
		t.Fatal("expected error for HTTP 500, got nil")
	}
	if !strings.Contains(err.Error(), "API-Anfrage fehlgeschlagen") {
		t.Errorf("expected API error message, got: %v", err)
	}
}

func TestExecuteSearch_ParseError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`not json`))
	}))
	defer srv.Close()

	cmd := setupTestCmd(srv.URL)
	defer os.Unsetenv("RIS_BASE_URL")

	params := api.NewParams()
	params.Set("Suchworte", "test")

	err := executeSearch(cmd, "Bundesrecht", "Suche...", params)
	if err == nil {
		t.Fatal("expected parse error, got nil")
	}
	if !strings.Contains(err.Error(), "Antwort konnte nicht verarbeitet werden") {
		t.Errorf("expected parse error message, got: %v", err)
	}
}

func TestExecuteSearch_SetsPageParams(t *testing.T) {
	var receivedPage, receivedLimit string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPage = r.URL.Query().Get("Seitennummer")
		receivedLimit = r.URL.Query().Get("DokumenteProSeite")
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(minimalAPIResponse))
	}))
	defer srv.Close()

	cmd := setupTestCmd(srv.URL)
	defer os.Unsetenv("RIS_BASE_URL")

	cmd.PersistentFlags().Set("page", "3")
	cmd.PersistentFlags().Set("limit", "50")

	params := api.NewParams()
	params.Set("Suchworte", "test")

	err := executeSearch(cmd, "Bundesrecht", "Suche...", params)
	if err != nil {
		t.Fatalf("executeSearch returned error: %v", err)
	}
	if receivedPage != "3" {
		t.Errorf("expected Seitennummer=3, got %q", receivedPage)
	}
	if receivedLimit != "Fifty" {
		t.Errorf("expected DokumenteProSeite=Fifty, got %q", receivedLimit)
	}
}

func TestExecuteSearch_InvalidLimit(t *testing.T) {
	cmd := setupTestCmd("http://unused")
	defer os.Unsetenv("RIS_BASE_URL")

	cmd.PersistentFlags().Set("limit", "3")

	params := api.NewParams()
	params.Set("Suchworte", "test")

	err := executeSearch(cmd, "Bundesrecht", "Suche...", params)
	if err == nil {
		t.Fatal("expected validation error for --limit 3, got nil")
	}
	var ve *ValidationError
	if !errors.As(err, &ve) {
		t.Errorf("expected ValidationError, got %T: %v", err, err)
	}
	if !strings.Contains(err.Error(), "--limit") {
		t.Errorf("expected error to mention --limit, got: %v", err)
	}
}
