package format

import (
	"fmt"
	"io"
	"strings"

	"github.com/fatih/color"
	"github.com/philrox/risgo/internal/model"
)

// Color functions for styled output.
var (
	bold       = color.New(color.Bold).SprintFunc()
	boldWhite  = color.New(color.Bold, color.FgHiWhite).SprintFunc()
	cyan       = color.New(color.FgCyan).SprintFunc()
	yellow     = color.New(color.FgYellow).SprintFunc()
	green      = color.New(color.FgGreen).SprintFunc()
	dim        = color.New(color.Faint).SprintFunc()
	boldYellow = color.New(color.Bold, color.FgYellow).SprintFunc()
)

const (
	// separatorWidth is the character width for horizontal rule separators.
	separatorWidth = 60
	// maxLeitsatzPreview is the maximum character length for Leitsatz previews in search results.
	maxLeitsatzPreview = 200
)

// Text writes search results as human-readable text to the writer.
func Text(w io.Writer, result model.SearchResult) error {
	if len(result.Documents) == 0 {
		fmt.Fprintln(w, "Keine Ergebnisse gefunden.")
		return nil
	}

	// Header with pagination info.
	fmt.Fprintln(w, bold(fmt.Sprintf("Ergebnisse: %d gesamt (Seite %d, zeige %d)",
		result.TotalHits, result.Page, len(result.Documents))))
	fmt.Fprintln(w, dim(strings.Repeat("─", separatorWidth)))

	for i, doc := range result.Documents {
		fmt.Fprintf(w, "\n[%d] %s\n", i+1, boldWhite(docTitle(doc)))

		if doc.Dokumentnummer != "" {
			fmt.Fprintf(w, "    Nr: %s\n", cyan(doc.Dokumentnummer))
		}

		citation := FormatCitation(doc.Citation)
		if citation != "" {
			fmt.Fprintf(w, "    Zitat: %s\n", citation)
		}

		if doc.Geschaeftszahl != "" {
			fmt.Fprintf(w, "    GZ: %s\n", green(doc.Geschaeftszahl))
		}

		dates := FormatDates(doc.Citation)
		if dates != "" {
			fmt.Fprintf(w, "    Geltung: %s\n", dim(dates))
		}

		if doc.Citation != nil && doc.Citation.Eli != "" {
			fmt.Fprintf(w, "    ELI: %s\n", dim(doc.Citation.Eli))
		}

		if doc.Leitsatz != "" {
			leitsatz := doc.Leitsatz
			if len(leitsatz) > maxLeitsatzPreview {
				leitsatz = leitsatz[:maxLeitsatzPreview] + "..."
			}
			fmt.Fprintf(w, "    Leitsatz: %s\n", leitsatz)
		}
	}

	fmt.Fprintln(w)
	if result.HasMore {
		nextPage := result.Page + 1
		fmt.Fprintln(w, boldYellow(fmt.Sprintf("Weitere Ergebnisse verfügbar. Nächste Seite: --page %d", nextPage)))
	}

	return nil
}

// TextDocument writes a single document with its content as human-readable text.
func TextDocument(w io.Writer, doc model.Document, content string) error {
	title := docTitle(doc)
	fmt.Fprintln(w, bold(title))
	fmt.Fprintln(w, dim(strings.Repeat("═", len(title))))
	fmt.Fprintln(w)

	if doc.Dokumentnummer != "" {
		fmt.Fprintf(w, "Dokument: %s\n", cyan(doc.Dokumentnummer))
	}

	citation := FormatCitation(doc.Citation)
	if citation != "" {
		fmt.Fprintf(w, "Zitat: %s\n", citation)
	}

	dates := FormatDates(doc.Citation)
	if dates != "" {
		fmt.Fprintf(w, "Geltung: %s\n", dim(dates))
	}

	if doc.Citation != nil && doc.Citation.Eli != "" {
		fmt.Fprintf(w, "ELI: %s\n", dim(doc.Citation.Eli))
	}

	if content != "" {
		fmt.Fprintln(w)
		fmt.Fprintln(w, dim(strings.Repeat("─", separatorWidth)))
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
