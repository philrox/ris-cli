package parser

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/philrox/risgo/internal/model"
)

// noExpiryDate is the sentinel value the API uses for laws with no expiry.
const noExpiryDate = "9999-12-31"

// ParseSearchResponse parses the raw JSON API response into a SearchResult.
func ParseSearchResponse(data []byte) (model.SearchResult, error) {
	var raw rawResponse
	if err := json.Unmarshal(data, &raw); err != nil {
		return model.SearchResult{}, fmt.Errorf("failed to parse API response: %w", err)
	}

	results := raw.OgdSearchResult.OgdDocumentResults

	// Parse hits metadata.
	totalHits, page, pageSize := parseHits(results.Hits)

	// Parse documents.
	var docs []model.Document
	for _, ref := range results.Docs {
		doc := parseDocumentReference(ref)
		docs = append(docs, doc)
	}

	hasMore := (page * pageSize) < totalHits

	return model.SearchResult{
		TotalHits: totalHits,
		Page:      page,
		PageSize:  pageSize,
		HasMore:   hasMore,
		Documents: docs,
	}, nil
}

// parseHits extracts total hits, page number, and page size from the
// polymorphic Hits field.
func parseHits(raw json.RawMessage) (totalHits, page, pageSize int) {
	if raw == nil {
		return 0, 1, 20
	}

	// Try as object with attributes.
	var obj rawHitsObject
	if err := json.Unmarshal(raw, &obj); err == nil && obj.Text != "" {
		fmt.Sscanf(obj.Text.String(), "%d", &totalHits)
		fmt.Sscanf(obj.PageNumber.String(), "%d", &page)
		fmt.Sscanf(obj.PageSize.String(), "%d", &pageSize)
		if page == 0 {
			page = 1
		}
		if pageSize == 0 {
			pageSize = 20
		}
		return
	}

	// Try as plain number.
	var n FlexibleInt
	if err := json.Unmarshal(raw, &n); err == nil {
		return int(n), 1, 20
	}

	return 0, 1, 20
}

// parseDocumentReference converts a raw document reference into a model.Document.
func parseDocumentReference(ref rawDocumentReference) model.Document {
	data := ref.Data
	meta := data.Metadaten

	doc := model.Document{
		Dokumentnummer: meta.Technisch.ID,
		Applikation:    meta.Technisch.Applikation,
		DokumentURL:    meta.Allgemein.DokumentURL,
	}

	// Parse domain-specific metadata.
	if meta.Bundesrecht != nil {
		parseBundesrecht(meta.Bundesrecht, &doc)
	}
	if meta.Landesrecht != nil {
		parseLandesrecht(meta.Landesrecht, &doc)
	}
	if meta.Judikatur != nil {
		parseJudikatur(meta.Judikatur, &doc)
	}

	// Parse content URLs from Dokumentliste.
	if data.Dokumentliste != nil {
		parseContentURLs(data.Dokumentliste, &doc)
	}

	return doc
}

func parseBundesrecht(raw json.RawMessage, doc *model.Document) {
	var br rawBundesrecht
	if err := json.Unmarshal(raw, &br); err != nil {
		return
	}

	doc.Kurztitel = br.Kurztitel
	doc.Titel = br.Titel.String()

	cit := &model.Citation{
		Kurztitel: br.Kurztitel,
		Langtitel: br.Langtitel,
		Eli:       br.Eli,
	}

	// Find the active sub-application section.
	subApp := firstNonNil(br.BrKons, br.Begut, br.BgblAuth, br.Erv, br.BgblPdf, br.BgblAlt, br.RegV)
	if subApp != nil {
		cit.Kundmachungsorgan = subApp.Kundmachungsorgan
		cit.Paragraph = subApp.ArtikelParagraphAnlage.String()
		cit.Inkrafttreten = subApp.Inkrafttretensdatum
		akt := subApp.Ausserkrafttretensdatum
		if akt != "" && akt != noExpiryDate {
			cit.Ausserkrafttreten = &akt
		}
		doc.GesamteRechtsvorschriftURL = subApp.GesamteRechtsvorschriftURL
	}

	doc.Citation = cit
}

func parseLandesrecht(raw json.RawMessage, doc *model.Document) {
	var lr rawLandesrecht
	if err := json.Unmarshal(raw, &lr); err != nil {
		return
	}

	doc.Kurztitel = lr.Kurztitel
	doc.Titel = lr.Titel.String()

	cit := &model.Citation{
		Kurztitel: lr.Kurztitel,
		Langtitel: lr.Langtitel,
		Eli:       lr.Eli,
	}

	subApp := firstNonNil(lr.LrKons, lr.LgblAuth, lr.Lgbl, lr.LgblNO, lr.Vbl, lr.Gr, lr.GrA)
	if subApp != nil {
		cit.Kundmachungsorgan = subApp.Kundmachungsorgan
		cit.Paragraph = subApp.ArtikelParagraphAnlage.String()
		cit.Inkrafttreten = subApp.Inkrafttretensdatum
		akt := subApp.Ausserkrafttretensdatum
		if akt != "" && akt != noExpiryDate {
			cit.Ausserkrafttreten = &akt
		}
		doc.GesamteRechtsvorschriftURL = subApp.GesamteRechtsvorschriftURL
	}

	doc.Citation = cit
}

func parseJudikatur(raw json.RawMessage, doc *model.Document) {
	var jud rawJudikatur
	if err := json.Unmarshal(raw, &jud); err != nil {
		return
	}

	doc.Kurztitel = jud.Kurztitel
	doc.Titel = jud.Titel.String()
	doc.Geschaeftszahl = jud.Geschaeftszahl.Value

	cit := &model.Citation{
		Kurztitel:      jud.Kurztitel,
		Langtitel:      jud.Langtitel,
		Geschaeftszahl: jud.Geschaeftszahl.Value,
	}

	// Find active sub-application. Leitsatz only for Vfgh, Vwgh, Justiz, Bvwg.
	type judApp struct {
		app         *rawJustizApp
		hasLeitsatz bool
	}
	apps := []judApp{
		{jud.Vfgh, true},
		{jud.Vwgh, true},
		{jud.Justiz, true},
		{jud.Bvwg, true},
		{jud.Lvwg, false},
		{jud.Dsk, false},
		{jud.Gbk, false},
		{jud.Pvak, false},
		{jud.AsylGH, false},
		{jud.Dok, false},
	}

	for _, a := range apps {
		if a.app != nil {
			cit.Entscheidungsdatum = a.app.Entscheidungsdatum
			cit.Inkrafttreten = a.app.Entscheidungsdatum // Semantic mapping
			if a.hasLeitsatz {
				leitsatz := a.app.Leitsatz.String()
				if leitsatz != "" {
					cit.Leitsatz = leitsatz
					doc.Leitsatz = leitsatz
				}
			}
			break
		}
	}

	doc.Citation = cit
}

func parseContentURLs(raw json.RawMessage, doc *model.Document) {
	var dl rawDokumentliste
	if err := json.Unmarshal(raw, &dl); err != nil {
		return
	}

	for _, cr := range dl.ContentReference {
		if !strings.EqualFold(cr.ContentType, "MainDocument") {
			continue
		}
		for _, cu := range cr.Urls.ContentUrl {
			switch strings.ToLower(cu.DataType) {
			case "html":
				doc.ContentURLs.HTML = cu.URL
			case "xml":
				doc.ContentURLs.XML = cu.URL
			case "pdf":
				doc.ContentURLs.PDF = cu.URL
			case "rtf":
				doc.ContentURLs.RTF = cu.URL
			}
		}
	}
}

func firstNonNil(ptrs ...*rawSubApp) *rawSubApp {
	for _, p := range ptrs {
		if p != nil {
			return p
		}
	}
	return nil
}
