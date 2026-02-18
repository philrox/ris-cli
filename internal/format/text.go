package format

import (
	"fmt"
	"io"
	"strings"

	"github.com/philrox/ris-cli/internal/model"
)

// Text writes search results as human-readable text to the writer.
func Text(w io.Writer, result model.SearchResult) error {
	if len(result.Documents) == 0 {
		fmt.Fprintln(w, "Keine Ergebnisse gefunden.")
		return nil
	}

	// Header with pagination info.
	fmt.Fprintf(w, "Ergebnisse: %d gesamt (Seite %d, zeige %d)\n",
		result.TotalHits, result.Page, len(result.Documents))
	fmt.Fprintln(w, strings.Repeat("─", 60))

	for i, doc := range result.Documents {
		fmt.Fprintf(w, "\n[%d] %s\n", i+1, docTitle(doc))

		if doc.Dokumentnummer != "" {
			fmt.Fprintf(w, "    Nr: %s\n", doc.Dokumentnummer)
		}

		citation := FormatCitation(doc.Citation)
		if citation != "" {
			fmt.Fprintf(w, "    Zitat: %s\n", citation)
		}

		if doc.Geschaeftszahl != "" {
			fmt.Fprintf(w, "    GZ: %s\n", doc.Geschaeftszahl)
		}

		if doc.Citation != nil && doc.Citation.Entscheidungsdatum != "" {
			fmt.Fprintf(w, "    Datum: %s\n", doc.Citation.Entscheidungsdatum)
		}

		if doc.Leitsatz != "" {
			leitsatz := doc.Leitsatz
			if len(leitsatz) > 200 {
				leitsatz = leitsatz[:200] + "..."
			}
			fmt.Fprintf(w, "    Leitsatz: %s\n", leitsatz)
		}
	}

	fmt.Fprintln(w)
	if result.HasMore {
		nextPage := result.Page + 1
		fmt.Fprintf(w, "Weitere Ergebnisse verfügbar. Nächste Seite: --page %d\n", nextPage)
	}

	return nil
}

// TextDocument writes a single document with its content as human-readable text.
func TextDocument(w io.Writer, doc model.Document, content string) error {
	title := docTitle(doc)
	fmt.Fprintln(w, title)
	fmt.Fprintln(w, strings.Repeat("═", len(title)))
	fmt.Fprintln(w)

	if doc.Dokumentnummer != "" {
		fmt.Fprintf(w, "Dokument: %s\n", doc.Dokumentnummer)
	}

	citation := FormatCitation(doc.Citation)
	if citation != "" {
		fmt.Fprintf(w, "Zitat: %s\n", citation)
	}

	if content != "" {
		fmt.Fprintln(w)
		fmt.Fprintln(w, strings.Repeat("─", 60))
		fmt.Fprintln(w)
		fmt.Fprintln(w, content)
	}

	return nil
}

func docTitle(doc model.Document) string {
	if doc.Titel != "" {
		return doc.Titel
	}
	if doc.Kurztitel != "" {
		return doc.Kurztitel
	}
	if doc.Dokumentnummer != "" {
		return doc.Dokumentnummer
	}
	return "(ohne Titel)"
}
