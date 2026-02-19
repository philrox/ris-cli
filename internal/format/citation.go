package format

import (
	"strings"

	"github.com/fatih/color"
	"github.com/philrox/risgo/internal/model"
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
	} else if c.Langtitel != "" {
		parts = append(parts, citationParagraph(c.Langtitel))
	}

	// Kundmachungsorgan in parentheses: "(JGS Nr. 946/1811)"
	if c.Kundmachungsorgan != "" {
		parts = append(parts, citationOrgan("("+c.Kundmachungsorgan+")"))
	}

	// Geschaeftszahl for court decisions when no paragraph present.
	if c.Geschaeftszahl != "" && c.Paragraph == "" {
		parts = append(parts, c.Geschaeftszahl)
	}

	// Entscheidungsdatum for court decisions.
	if c.Entscheidungsdatum != "" {
		parts = append(parts, citationOrgan("vom "+c.Entscheidungsdatum))
	}

	return strings.Join(parts, " ")
}

// FormatDates formats Inkrafttreten/Ausserkrafttreten into a display string.
// Returns empty string if no date info is available.
func FormatDates(c *model.Citation) string {
	if c == nil {
		return ""
	}

	var parts []string

	if c.Inkrafttreten != "" {
		parts = append(parts, "in Kraft seit "+c.Inkrafttreten)
	}

	if c.Ausserkrafttreten != nil && *c.Ausserkrafttreten != "" {
		parts = append(parts, "außer Kraft seit "+*c.Ausserkrafttreten)
	}

	return strings.Join(parts, ", ")
}
