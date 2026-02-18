package parser

import "encoding/json"

// rawResponse is the top-level API response.
type rawResponse struct {
	OgdSearchResult rawSearchResult `json:"OgdSearchResult"`
}

// rawSearchResult contains the document results.
type rawSearchResult struct {
	OgdDocumentResults rawDocumentResults `json:"OgdDocumentResults"`
}

// rawDocumentResults contains hit count and document references.
type rawDocumentResults struct {
	Hits json.RawMessage                    `json:"Hits"`
	Docs FlexibleArray[rawDocumentReference] `json:"OgdDocumentReference"`
}

// rawHitsObject represents the Hits when it's an object with attributes.
type rawHitsObject struct {
	Text       FlexibleString `json:"#text"`
	PageNumber FlexibleString `json:"@pageNumber"`
	PageSize   FlexibleString `json:"@pageSize"`
}

// rawDocumentReference represents a single document in the search results.
type rawDocumentReference struct {
	Data rawData `json:"Data"`
}

// rawData contains metadata and document list.
type rawData struct {
	Metadaten    rawMetadaten    `json:"Metadaten"`
	Dokumentliste json.RawMessage `json:"Dokumentliste"`
}

// rawMetadaten contains technical, general, and domain-specific metadata.
type rawMetadaten struct {
	Technisch    rawTechnisch    `json:"Technisch"`
	Allgemein    rawAllgemein    `json:"Allgemein"`
	Bundesrecht  json.RawMessage `json:"Bundesrecht,omitempty"`
	Landesrecht  json.RawMessage `json:"Landesrecht,omitempty"`
	Judikatur    json.RawMessage `json:"Judikatur,omitempty"`
}

type rawTechnisch struct {
	ID          string `json:"ID"`
	Applikation string `json:"Applikation"`
}

type rawAllgemein struct {
	DokumentURL string `json:"DokumentUrl"`
}

// rawBundesrecht is the Bundesrecht metadata section.
type rawBundesrecht struct {
	Kurztitel string         `json:"Kurztitel"`
	Langtitel string         `json:"Langtitel"`
	Titel     FlexibleString `json:"Titel"`
	Eli       string         `json:"Eli"`
	BrKons    *rawBrKons     `json:"BrKons,omitempty"`
	Begut     *rawBrKons     `json:"Begut,omitempty"`
	BgblAuth  *rawBrKons     `json:"BgblAuth,omitempty"`
	Erv       *rawBrKons     `json:"Erv,omitempty"`
	BgblPdf   *rawBrKons     `json:"BgblPdf,omitempty"`
	BgblAlt   *rawBrKons     `json:"BgblAlt,omitempty"`
	RegV      *rawBrKons     `json:"RegV,omitempty"`
}

// rawBrKons is a Bundesrecht sub-application section.
type rawBrKons struct {
	Kundmachungsorgan             string         `json:"Kundmachungsorgan"`
	ArtikelParagraphAnlage        FlexibleString `json:"ArtikelParagraphAnlage"`
	Inkrafttretensdatum           string         `json:"Inkrafttretensdatum"`
	Ausserkrafttretensdatum       string         `json:"Ausserkrafttretensdatum"`
	GesamteRechtsvorschriftURL    string         `json:"GesamteRechtsvorschriftUrl"`
}

// rawLandesrecht is the Landesrecht metadata section.
type rawLandesrecht struct {
	Kurztitel string         `json:"Kurztitel"`
	Langtitel string         `json:"Langtitel"`
	Titel     FlexibleString `json:"Titel"`
	Eli       string         `json:"Eli"`
	LrKons    *rawLrKons     `json:"LrKons,omitempty"`
	LgblAuth  *rawLrKons     `json:"LgblAuth,omitempty"`
	Lgbl      *rawLrKons     `json:"Lgbl,omitempty"`
	LgblNO    *rawLrKons     `json:"LgblNO,omitempty"`
	Vbl       *rawLrKons     `json:"Vbl,omitempty"`
	Gr        *rawLrKons     `json:"Gr,omitempty"`
	GrA       *rawLrKons     `json:"GrA,omitempty"`
}

// rawLrKons is a Landesrecht sub-application section.
type rawLrKons struct {
	Kundmachungsorgan             string         `json:"Kundmachungsorgan"`
	ArtikelParagraphAnlage        FlexibleString `json:"ArtikelParagraphAnlage"`
	Inkrafttretensdatum           string         `json:"Inkrafttretensdatum"`
	Ausserkrafttretensdatum       string         `json:"Ausserkrafttretensdatum"`
	GesamteRechtsvorschriftURL    string         `json:"GesamteRechtsvorschriftUrl"`
}

// rawJudikatur is the Judikatur metadata section.
type rawJudikatur struct {
	Kurztitel      string            `json:"Kurztitel"`
	Langtitel      string            `json:"Langtitel"`
	Titel          FlexibleString    `json:"Titel"`
	Geschaeftszahl rawGeschaeftszahl `json:"Geschaeftszahl"`
	Justiz         *rawJustizApp     `json:"Justiz,omitempty"`
	Vfgh           *rawJustizApp     `json:"Vfgh,omitempty"`
	Vwgh           *rawJustizApp     `json:"Vwgh,omitempty"`
	Bvwg           *rawJustizApp     `json:"Bvwg,omitempty"`
	Lvwg           *rawJustizApp     `json:"Lvwg,omitempty"`
	Dsk            *rawJustizApp     `json:"Dsk,omitempty"`
	Gbk            *rawJustizApp     `json:"Gbk,omitempty"`
	Pvak           *rawJustizApp     `json:"Pvak,omitempty"`
	AsylGH         *rawJustizApp     `json:"AsylGH,omitempty"`
	Dok            *rawJustizApp     `json:"Dok,omitempty"`
}

// rawGeschaeftszahl handles the polymorphic Geschaeftszahl field.
type rawGeschaeftszahl struct {
	Value string
}

func (g *rawGeschaeftszahl) UnmarshalJSON(data []byte) error {
	// Try plain string.
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		g.Value = s
		return nil
	}

	// Try object with item field.
	var obj struct {
		Item json.RawMessage `json:"item"`
	}
	if err := json.Unmarshal(data, &obj); err == nil && obj.Item != nil {
		// Item can be string or []string.
		var itemStr string
		if err := json.Unmarshal(obj.Item, &itemStr); err == nil {
			g.Value = itemStr
			return nil
		}
		var items []string
		if err := json.Unmarshal(obj.Item, &items); err == nil && len(items) > 0 {
			g.Value = items[0]
			return nil
		}
	}

	g.Value = ""
	return nil
}

// rawJustizApp is a Judikatur sub-application section.
type rawJustizApp struct {
	Entscheidungsdatum string         `json:"Entscheidungsdatum"`
	Leitsatz           FlexibleString `json:"Leitsatz"`
	Norm               FlexibleString `json:"Norm"`
}

// rawDokumentliste contains content references.
type rawDokumentliste struct {
	ContentReference FlexibleArray[rawContentReference] `json:"ContentReference"`
}

// rawContentReference describes a document content entry.
type rawContentReference struct {
	ContentType string         `json:"ContentType"`
	Name        FlexibleString `json:"Name"`
	Urls        rawURLs        `json:"Urls"`
}

// rawURLs contains the content URLs.
type rawURLs struct {
	ContentUrl FlexibleArray[rawContentURL] `json:"ContentUrl"`
}

// rawContentURL is a single content URL entry.
type rawContentURL struct {
	DataType string `json:"DataType"`
	URL      string `json:"Url"`
}
