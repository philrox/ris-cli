package constants

// PageSizes maps numeric page sizes to their API enum values.
var PageSizes = map[int]string{
	10:  "Ten",
	20:  "Twenty",
	50:  "Fifty",
	100: "OneHundred",
}

// ImRisSeit maps CLI time filter values to their API enum values.
var ImRisSeit = map[string]string{
	"einerwoche":   "EinerWoche",
	"zweiwochen":   "ZweiWochen",
	"einemmonat":   "EinemMonat",
	"dreimonaten":  "DreiMonaten",
	"sechsmonaten": "SechsMonaten",
	"einemjahr":    "EinemJahr",
}

// SortDirections maps CLI sort direction values to their API enum values.
var SortDirections = map[string]string{
	"asc":  "Ascending",
	"desc": "Descending",
}

// BundesrechtApps maps CLI app values to their API Applikation values
// for the bundesrecht command.
var BundesrechtApps = map[string]string{
	"brkons":   "BrKons",
	"begut":    "Begut",
	"bgblauth": "BgblAuth",
	"erv":      "Erv",
}

// BgblApps maps CLI app values to their API Applikation values
// for the bgbl command.
var BgblApps = map[string]string{
	"bgblauth": "BgblAuth",
	"bgblpdf":  "BgblPdf",
	"bgblalt":  "BgblAlt",
}

// BgblTeile maps CLI part values to their API Teil values.
var BgblTeile = map[string]string{
	"1": "Eins",
	"2": "Zwei",
	"3": "Drei",
}

// LgblApps maps CLI app values to their API Applikation values
// for the lgbl command.
var LgblApps = map[string]string{
	"lgblauth": "LgblAuth",
	"lgbl":     "Lgbl",
	"lgblno":   "LgblNO",
}

// GemeindenApps maps CLI app values to their API Applikation values
// for the gemeinden command.
var GemeindenApps = map[string]string{
	"gr":  "Gr",
	"gra": "GrA",
}

// GemeindenIndex maps CLI index values to API Index values for Gemeinden (Gr only).
var GemeindenIndex = map[string]string{
	"undefined":                                        "Undefined",
	"vertretungskoerperundalgemeineverwaltung":          "VertretungskoerperUndAllgemeineVerwaltung",
	"oeffentlicheordnungundsicherheit":                  "OeffentlicheOrdnungUndSicherheit",
	"unterrichterziehungsportunwissenschaft":             "UnterrichtErziehungSportUndWissenschaft",
	"kunstkulturunddkultus":                             "KunstKulturUndKultus",
	"sozialewohfahrtundwohnbaufoerderung":                "SozialeWohlfahrtUndWohnbaufoerderung",
	"gesundheit":                                        "Gesundheit",
	"strassenundwasserbauverkehr":                        "StraßenUndWasserbauVerkehr",
	"wirtschaftsfoerderung":                             "Wirtschaftsfoerderung",
	"dienstleistungen":                                  "Dienstleistungen",
	"finanzwirtschaft":                                  "Finanzwirtschaft",
}

// GemeindenSortColumns maps CLI sort-by values to API SortedByColumn values for Gemeinden (Gr only).
var GemeindenSortColumns = map[string]string{
	"geschaeftszahl": "Geschaeftszahl",
	"bundesland":     "Bundesland",
	"gemeinde":       "Gemeinde",
}

// RegvorlSortColumns maps CLI sort-by values to API SortedByColumn values for regvorl.
var RegvorlSortColumns = map[string]string{
	"kurztitel": "Kurztitel",
	"stelle":    "EinbringendeStelle",
	"datum":     "Beschlussdatum",
}

// UptsParties maps CLI party values to API Partei values.
var UptsParties = map[string]string{
	"spoe":   "SPÖ - Sozialdemokratische Partei Österreichs",
	"oevp":   "ÖVP - Österreichische Volkspartei",
	"fpoe":   "FPÖ - Freiheitliche Partei Österreichs",
	"gruene": "GRÜNE - Die Grünen - Die Grüne Alternative",
	"neos":   "NEOS - NEOS – Das Neue Österreich und Liberales Forum",
	"bzoe":   "BZÖ - Bündnis Zukunft Österreich",
}

// KmgerTypes maps CLI type values to API Typ values.
var KmgerTypes = map[string]string{
	"geschaeftsordnung":    "Geschaeftsordnung",
	"geschaeftsverteilung": "Geschaeftsverteilung",
}

// AvsvAuthors maps CLI author values to API Urheber values.
var AvsvAuthors = map[string]string{
	"dvsv":  "Dachverband der Sozialversicherungsträger (DVSV)",
	"pva":   "Pensionsversicherungsanstalt (PVA)",
	"oegk":  "Österreichische Gesundheitskasse (ÖGK)",
	"auva":  "Allgemeine Unfallversicherungsanstalt (AUVA)",
	"svs":   "Sozialversicherungsanstalt der Selbständigen (SVS)",
	"bvaeb": "Versicherungsanstalt öffentlich Bediensteter, Eisenbahnen und Bergbau (BVAEB)",
}

// AvnTypes maps CLI type values to API Typ values.
var AvnTypes = map[string]string{
	"kundmachung": "Kundmachung",
	"verordnung":  "Verordnung",
	"erlass":      "Erlass",
}

// OsgTypes maps CLI OSG type values to API OsgTyp values.
var OsgTypes = map[string]string{
	"oesg":              "ÖSG",
	"oesg-grossgeraete": "ÖSG - Großgeräteplan",
}

// RsgTypes maps CLI RSG type values to API RsgTyp values.
var RsgTypes = map[string]string{
	"rsg":              "RSG",
	"rsg-grossgeraete": "RSG - Großgeräteplan",
}

// PruefgewoTypes maps CLI type values to API Typ values.
var PruefgewoTypes = map[string]string{
	"befaehigung": "Befähigungsprüfung",
	"eignung":     "Eignungsprüfung",
	"meister":     "Meisterprüfung",
}

// historyApps is the set of valid application values for the History endpoint.
var historyApps = map[string]bool{
	"bundesnormen":      true,
	"landesnormen":      true,
	"justiz":            true,
	"vfgh":              true,
	"vwgh":              true,
	"bvwg":              true,
	"lvwg":              true,
	"bgblauth":          true,
	"bgblalt":           true,
	"bgblpdf":           true,
	"lgblauth":          true,
	"lgbl":              true,
	"lgblno":            true,
	"gemeinderecht":     true,
	"gemeinderechtauth": true,
	"bvb":               true,
	"vbl":               true,
	"regv":              true,
	"mrp":               true,
	"erlaesse":          true,
	"pruefgewo":         true,
	"avsv":              true,
	"spg":               true,
	"kmger":             true,
	"dsk":               true,
	"gbk":               true,
	"dok":               true,
	"pvak":              true,
	"normenliste":       true,
	"asylgh":            true,
}

// IsValidHistoryApp returns true if the given application name is valid for the History endpoint.
func IsValidHistoryApp(app string) bool {
	return historyApps[app]
}
