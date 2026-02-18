package cmd

import (
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/philrox/ris-cli/internal/api"
	"github.com/philrox/ris-cli/internal/format"
	"github.com/philrox/ris-cli/internal/model"
	"github.com/philrox/ris-cli/internal/parser"
	"github.com/philrox/ris-cli/internal/ui"
	"github.com/spf13/cobra"
)

var dokumentCmd = &cobra.Command{
	Use:   "dokument [document-number]",
	Short: "Volltext eines Dokuments abrufen",
	Long: `Volltext eines Rechtsdokuments abrufen.

Beispiele:
  ris dokument NOR40052761
  ris dokument NOR40052761 --json
  ris dokument --url "https://ris.bka.gv.at/Dokumente/Bundesnormen/NOR40052761/NOR40052761.html"`,
	Args: cobra.MaximumNArgs(1),
	RunE: runDokument,
}

func init() {
	f := dokumentCmd.Flags()
	f.String("url", "", "Direkte URL zum Dokumentinhalt")

	rootCmd.AddCommand(dokumentCmd)
}

var docNumberRegex = regexp.MustCompile(`^[A-Z][A-Z0-9_]+$`)

// allowedHosts for SSRF protection.
var allowedHosts = map[string]bool{
	"data.bka.gv.at":     true,
	"www.ris.bka.gv.at":  true,
	"ris.bka.gv.at":      true,
}

func validateDocNumber(nr string) error {
	if len(nr) < 5 || len(nr) > 50 {
		return fmt.Errorf("Dokumentnummer muss 5-50 Zeichen lang sein, erhalten: %d", len(nr))
	}
	if !docNumberRegex.MatchString(nr) {
		return fmt.Errorf("Ungültiges Dokumentnummer-Format %q (muss ^[A-Z][A-Z0-9_]+$ entsprechen)", nr)
	}
	return nil
}

func validateURL(rawURL string) error {
	u, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("Ungültige URL: %w", err)
	}
	if u.Scheme != "https" {
		return fmt.Errorf("Nur HTTPS-URLs erlaubt, erhalten: %q", u.Scheme)
	}
	host := strings.ToLower(u.Hostname())
	if !allowedHosts[host] {
		return fmt.Errorf("Host %q nicht erlaubt (erlaubt: data.bka.gv.at, www.ris.bka.gv.at, ris.bka.gv.at)", host)
	}
	return nil
}

func runDokument(cmd *cobra.Command, args []string) error {
	docURL, _ := cmd.Flags().GetString("url")
	var docNumber string
	if len(args) > 0 {
		docNumber = args[0]
	}

	if docNumber == "" && docURL == "" {
		fmt.Fprintln(os.Stderr, "Fehler: Dokumentnummer oder --url erforderlich")
		os.Exit(2)
	}

	client := newClient(cmd)

	if docURL != "" {
		// Direct URL fetch.
		if err := validateURL(docURL); err != nil {
			fmt.Fprintf(os.Stderr, "Fehler: %v\n", err)
			os.Exit(2)
		}
		return fetchAndOutputDocument(cmd, client, docURL, docNumber)
	}

	// Document number strategy.
	if err := validateDocNumber(docNumber); err != nil {
		fmt.Fprintf(os.Stderr, "Fehler: %v\n", err)
		os.Exit(2)
	}

	// Step 1: Try direct URL from prefix routing table.
	directURL := model.DirectURLFromPrefix(docNumber)
	if directURL != "" {
		s := startSpinner(cmd, "Lade Dokument...")
		htmlContent, err := client.FetchDocument(directURL)
		stopSpinner(s)
		if err == nil {
			return outputDocumentContent(cmd, docNumber, directURL, htmlContent)
		}
		// Direct URL failed, fall through to search.
		if isVerbose(cmd) {
			fmt.Fprintf(os.Stderr, "Direkte URL fehlgeschlagen (%v), versuche Suche als Fallback...\n", err)
		}
	}

	// Step 2: Fallback to search API.
	endpoint, applikation := model.SearchFallback(docNumber)
	params := api.NewParams()
	params.Set("Applikation", applikation)
	params.Set("Dokumentnummer", docNumber)
	params.Set("DokumenteProSeite", "Ten")

	s2 := startSpinner(cmd, "Suche Dokument-URL...")
	body, err := client.Search(endpoint, params)
	stopSpinner(s2)
	if err != nil {
		return fmt.Errorf("Such-API-Anfrage fehlgeschlagen: %w", err)
	}

	result, err := parser.ParseSearchResponse(body)
	if err != nil {
		return fmt.Errorf("Suchantwort konnte nicht verarbeitet werden: %w", err)
	}

	if len(result.Documents) == 0 {
		fmt.Fprintf(os.Stderr, "Fehler: Dokument %q nicht gefunden\n", docNumber)
		os.Exit(3)
	}

	// Find HTML content URL from search result.
	doc := result.Documents[0]
	htmlURL := ""
	if doc.ContentURLs.HTML != "" {
		htmlURL = doc.ContentURLs.HTML
	} else if doc.DokumentURL != "" {
		htmlURL = doc.DokumentURL
	}

	if htmlURL == "" {
		// No content URL found; output metadata only.
		if useJSON(cmd) {
			return format.JSONDocument(os.Stdout, doc, "")
		}
		w, cleanup := ui.NewPagerWriter(!usePager(cmd))
		defer cleanup()
		return format.TextDocument(w, doc, "")
	}

	return fetchAndOutputDocument(cmd, client, htmlURL, docNumber)
}

func fetchAndOutputDocument(cmd *cobra.Command, client *api.Client, docURL, docNumber string) error {
	s := startSpinner(cmd, "Lade Dokument...")
	htmlContent, err := client.FetchDocument(docURL)
	stopSpinner(s)
	if err != nil {
		return fmt.Errorf("Dokument konnte nicht abgerufen werden: %w", err)
	}
	return outputDocumentContent(cmd, docNumber, docURL, htmlContent)
}

// usePager returns true when pager should be used for document output.
func usePager(cmd *cobra.Command) bool {
	return !useJSON(cmd) && !plainOutput && !quiet && !noPager
}

func outputDocumentContent(cmd *cobra.Command, docNumber, docURL, htmlContent string) error {
	textContent := format.HTMLToText(htmlContent)

	if useJSON(cmd) {
		doc := model.Document{
			Dokumentnummer: docNumber,
			DokumentURL:    docURL,
		}
		return format.JSONDocument(os.Stdout, doc, textContent)
	}

	w, cleanup := ui.NewPagerWriter(!usePager(cmd))
	defer cleanup()

	fmt.Fprintln(w, textContent)
	return nil
}
