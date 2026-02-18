package format

import (
	"strings"

	"github.com/fatih/color"
	"github.com/philrox/ris-cli/internal/model"
)

// Citation color functions.
var (
	citationParagraph = color.New(color.FgYellow, color.Bold).SprintFunc()
	citationOrgan     = color.New(color.Faint).SprintFunc()
)

// FormatCitation formats a Citation into a human-readable Austrian legal citation string.
// Example: "§ 1295 ABGB (JGS Nr. 946/1811)"
func FormatCitation(c *model.Citation) string {
	if c == nil {
		return ""
	}

	var parts []string

	// Paragraph + short title: "§ 1295 ABGB"
	if c.Paragraph != "" && c.Kurztitel != "" {
		parts = append(parts, citationParagraph(c.Paragraph+" "+c.Kurztitel))
	} else if c.Kurztitel != "" {
		parts = append(parts, citationParagraph(c.Kurztitel))
	} else if c.Paragraph != "" {
		parts = append(parts, citationParagraph(c.Paragraph))
	}

	// Kundmachungsorgan in parentheses: "(JGS Nr. 946/1811)"
	if c.Kundmachungsorgan != "" {
		parts = append(parts, citationOrgan("("+c.Kundmachungsorgan+")"))
	}

	// Geschaeftszahl for court decisions.
	if c.Geschaeftszahl != "" && c.Paragraph == "" {
		if len(parts) == 0 {
			parts = append(parts, c.Geschaeftszahl)
		}
	}

	return strings.Join(parts, " ")
}
