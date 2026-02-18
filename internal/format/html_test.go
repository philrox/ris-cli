package format

import (
	"strings"
	"testing"
)

func TestHTMLToText_Basic(t *testing.T) {
	html := "<p>Hello World</p>"
	got := HTMLToText(html)
	if !strings.Contains(got, "Hello World") {
		t.Errorf("HTMLToText() = %q, want 'Hello World'", got)
	}
}

func TestHTMLToText_StripScript(t *testing.T) {
	html := "<p>Text</p><script>alert('xss')</script><p>More</p>"
	got := HTMLToText(html)
	if strings.Contains(got, "alert") {
		t.Error("HTMLToText should strip script content")
	}
	if !strings.Contains(got, "Text") || !strings.Contains(got, "More") {
		t.Error("HTMLToText should preserve non-script text")
	}
}

func TestHTMLToText_StripStyle(t *testing.T) {
	html := "<style>body{color:red}</style><p>Content</p>"
	got := HTMLToText(html)
	if strings.Contains(got, "color") {
		t.Error("HTMLToText should strip style content")
	}
	if !strings.Contains(got, "Content") {
		t.Error("HTMLToText should preserve text content")
	}
}

func TestHTMLToText_StripHead(t *testing.T) {
	html := "<html><head><title>Title</title></head><body><p>Body</p></body></html>"
	got := HTMLToText(html)
	if strings.Contains(got, "Title") {
		t.Error("HTMLToText should strip head content")
	}
	if !strings.Contains(got, "Body") {
		t.Error("HTMLToText should preserve body text")
	}
}

func TestHTMLToText_LineBreaks(t *testing.T) {
	html := "Line1<br>Line2<br/>Line3"
	got := HTMLToText(html)
	if !strings.Contains(got, "Line1") || !strings.Contains(got, "Line2") || !strings.Contains(got, "Line3") {
		t.Errorf("HTMLToText() = %q, want all lines", got)
	}
}

func TestHTMLToText_BlockElements(t *testing.T) {
	html := "<div>Div1</div><div>Div2</div>"
	got := HTMLToText(html)
	if !strings.Contains(got, "Div1") || !strings.Contains(got, "Div2") {
		t.Errorf("HTMLToText() = %q, want both divs", got)
	}
}

func TestHTMLToText_WhitespaceNormalization(t *testing.T) {
	html := "<p>A</p><p></p><p></p><p></p><p></p><p>B</p>"
	got := HTMLToText(html)
	// Should not have more than one consecutive blank line.
	if strings.Contains(got, "\n\n\n") {
		t.Errorf("HTMLToText should normalize whitespace, got: %q", got)
	}
}

func TestHTMLToText_EmptyInput(t *testing.T) {
	got := HTMLToText("")
	if got != "" {
		t.Errorf("HTMLToText('') = %q, want empty", got)
	}
}

func TestNormalizeWhitespace(t *testing.T) {
	input := "A\n\n\n\n\nB\n\nC"
	got := normalizeWhitespace(input)
	if strings.Contains(got, "\n\n\n") {
		t.Error("normalizeWhitespace should collapse blank lines")
	}
	if !strings.Contains(got, "A") || !strings.Contains(got, "B") || !strings.Contains(got, "C") {
		t.Error("normalizeWhitespace should preserve content")
	}
}

func TestStripTagsSimple(t *testing.T) {
	input := "<p>Hello</p> <b>World</b>"
	got := stripTagsSimple(input)
	if strings.Contains(got, "<") || strings.Contains(got, ">") {
		t.Errorf("stripTagsSimple should remove all tags, got: %q", got)
	}
	if !strings.Contains(got, "Hello") || !strings.Contains(got, "World") {
		t.Errorf("stripTagsSimple should preserve text, got: %q", got)
	}
}
