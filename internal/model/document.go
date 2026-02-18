package model

// Document represents a single legal document from RIS.
type Document struct {
	Dokumentnummer             string      `json:"dokumentnummer"`
	Applikation                string      `json:"applikation"`
	Titel                      string      `json:"titel"`
	Kurztitel                  string      `json:"kurztitel"`
	Citation                   *Citation   `json:"citation,omitempty"`
	ContentURLs                ContentURLs `json:"content_urls"`
	DokumentURL                string      `json:"dokument_url"`
	GesamteRechtsvorschriftURL string      `json:"gesamte_rechtsvorschrift_url,omitempty"`
	Geschaeftszahl             string      `json:"geschaeftszahl,omitempty"`
	Leitsatz                   string      `json:"leitsatz,omitempty"`
}

// Citation contains structured legal citation information.
type Citation struct {
	Kurztitel          string  `json:"kurztitel"`
	Langtitel          string  `json:"langtitel,omitempty"`
	Kundmachungsorgan  string  `json:"kundmachungsorgan,omitempty"`
	Paragraph          string  `json:"paragraph,omitempty"`
	Eli                string  `json:"eli,omitempty"`
	Inkrafttreten      string  `json:"inkrafttreten,omitempty"`
	Ausserkrafttreten  *string `json:"ausserkrafttreten"`
	Geschaeftszahl     string  `json:"geschaeftszahl,omitempty"`
	Entscheidungsdatum string  `json:"entscheidungsdatum,omitempty"`
	Leitsatz           string  `json:"leitsatz,omitempty"`
}

// ContentURLs holds URLs for different document formats.
type ContentURLs struct {
	HTML string `json:"html"`
	XML  string `json:"xml"`
	PDF  string `json:"pdf"`
	RTF  string `json:"rtf"`
}

// DocumentContent represents a full document with its text content.
type DocumentContent struct {
	Metadata Document `json:"metadata"`
	Content  string   `json:"content"`
}
