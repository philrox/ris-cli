package format

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

// HTMLToText converts HTML content to plain text.
// Strips script, style, and head elements. Normalizes whitespace.
func HTMLToText(htmlContent string) string {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		// Fallback: strip tags with simple approach.
		return stripTagsSimple(htmlContent)
	}

	var sb strings.Builder
	extractText(doc, &sb)

	// Normalize whitespace: collapse multiple blank lines.
	text := sb.String()
	text = normalizeWhitespace(text)
	return strings.TrimSpace(text)
}

// extractText walks the HTML tree and extracts text content,
// skipping script, style, and head elements.
func extractText(n *html.Node, w io.StringWriter) {
	if n.Type == html.ElementNode {
		switch strings.ToLower(n.Data) {
		case "script", "style", "head", "noscript":
			return
		case "br":
			w.WriteString("\n")
		case "p", "div", "h1", "h2", "h3", "h4", "h5", "h6",
			"li", "tr", "blockquote", "pre", "table":
			w.WriteString("\n")
		}
	}

	if n.Type == html.TextNode {
		text := strings.TrimSpace(n.Data)
		if text != "" {
			w.WriteString(text)
			w.WriteString(" ")
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		extractText(c, w)
	}

	if n.Type == html.ElementNode {
		switch strings.ToLower(n.Data) {
		case "p", "div", "h1", "h2", "h3", "h4", "h5", "h6",
			"li", "tr", "blockquote", "pre", "table":
			w.WriteString("\n")
		}
	}
}

// normalizeWhitespace collapses sequences of blank lines into at most two newlines.
func normalizeWhitespace(s string) string {
	lines := strings.Split(s, "\n")
	var result []string
	blankCount := 0
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			blankCount++
			if blankCount <= 1 {
				result = append(result, "")
			}
		} else {
			blankCount = 0
			result = append(result, trimmed)
		}
	}
	return strings.Join(result, "\n")
}

// stripTagsSimple is a basic fallback HTML tag stripper.
func stripTagsSimple(s string) string {
	var result strings.Builder
	inTag := false
	for _, r := range s {
		switch {
		case r == '<':
			inTag = true
		case r == '>':
			inTag = false
			result.WriteRune(' ')
		case !inTag:
			result.WriteRune(r)
		}
	}
	return result.String()
}
