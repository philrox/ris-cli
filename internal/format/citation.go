package format

import (
	"strings"

	"github.com/philrox/ris-cli/internal/model"
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
		parts = append(parts, c.Paragraph+" "+c.Kurztitel)
	} else if c.Kurztitel != "" {
		parts = append(parts, c.Kurztitel)
	} else if c.Paragraph != "" {
		parts = append(parts, c.Paragraph)
	}

	// Kundmachungsorgan in parentheses: "(JGS Nr. 946/1811)"
	if c.Kundmachungsorgan != "" {
		parts = append(parts, "("+c.Kundmachungsorgan+")")
	}

	// Geschaeftszahl for court decisions.
	if c.Geschaeftszahl != "" && c.Paragraph == "" {
		if len(parts) == 0 {
			parts = append(parts, c.Geschaeftszahl)
		}
	}

	return strings.Join(parts, " ")
}
