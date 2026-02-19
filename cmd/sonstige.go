package cmd

import (
	"strings"

	"github.com/philrox/risgo/internal/api"
	"github.com/philrox/risgo/internal/constants"
	"github.com/spf13/cobra"
)

var sonstigeCmd = &cobra.Command{
	Use:   "sonstige",
	Short: "Sonstige Rechtssammlungen durchsuchen",
	Long: `Sonstige Rechtssammlungen durchsuchen (8 Teil-Applikationen).

Unterbefehle:
  mrp          Ministerratsprotokolle
  erlaesse     Erlässe
  upts         Parteientransparenz-Entscheidungen
  kmger        Gerichtskundmachungen
  avsv         Sozialversicherungs-Kundmachungen
  avn          Veterinär-Kundmachungen
  spg          Gesundheitsstrukturpläne
  pruefgewo    Gewerbeprüfungen`,
}

// --- mrp ---
var mrpCmd = &cobra.Command{
	Use:   "mrp",
	Short: "Ministerratsprotokolle durchsuchen",
	Long: `Ministerratsprotokolle durchsuchen.

Beispiele:
  risgo sonstige mrp --search "Budget"
  risgo sonstige mrp --session 42`,
	RunE: runMrp,
}

// --- erlaesse ---
var erlCmd = &cobra.Command{
	Use:   "erlaesse",
	Short: "Erlässe durchsuchen",
	Long: `Erlässe durchsuchen.

Beispiele:
  risgo sonstige erlaesse --ministry bmf
  risgo sonstige erlaesse --search "Steuer"`,
	RunE: runErlaesse,
}

// --- upts ---
var uptsCmd = &cobra.Command{
	Use:   "upts",
	Short: "Parteientransparenz-Entscheidungen durchsuchen",
	Long: `Parteientransparenz-Entscheidungen durchsuchen (UPTS).

Beispiele:
  risgo sonstige upts --party spoe`,
	RunE: runUpts,
}

// --- kmger ---
var kmgerCmd = &cobra.Command{
	Use:   "kmger",
	Short: "Gerichtskundmachungen durchsuchen",
	Long: `Gerichtskundmachungen durchsuchen (KmGer).

Beispiele:
  risgo sonstige kmger --type geschaeftsordnung`,
	RunE: runKmger,
}

// --- avsv ---
var avsvCmd = &cobra.Command{
	Use:   "avsv",
	Short: "Sozialversicherungs-Kundmachungen durchsuchen",
	Long: `Sozialversicherungs-Kundmachungen durchsuchen (AVSV).

Beispiele:
  risgo sonstige avsv --author dvsv`,
	RunE: runAvsv,
}

// --- avn ---
var avnCmd = &cobra.Command{
	Use:   "avn",
	Short: "Veterinär-Kundmachungen durchsuchen",
	Long: `Veterinär-Kundmachungen durchsuchen (AVN).

Beispiele:
  risgo sonstige avn --type kundmachung`,
	RunE: runAvn,
}

// --- spg ---
var spgCmd = &cobra.Command{
	Use:   "spg",
	Short: "Gesundheitsstrukturpläne durchsuchen",
	Long: `Gesundheitsstrukturpläne durchsuchen (SPG).

Beispiele:
  risgo sonstige spg --osg-type oesg`,
	RunE: runSpg,
}

// --- pruefgewo ---
var pruefgewoCmd = &cobra.Command{
	Use:   "pruefgewo",
	Short: "Gewerbeprüfungen durchsuchen",
	Long: `Gewerbeprüfungen durchsuchen (PrüfGewO).

Beispiele:
  risgo sonstige pruefgewo --type befaehigung`,
	RunE: runPruefgewo,
}

func init() {
	// Common flags for all sub-commands.
	for _, sub := range []*cobra.Command{mrpCmd, erlCmd, uptsCmd, kmgerCmd, avsvCmd, avnCmd, spgCmd, pruefgewoCmd} {
		f := sub.Flags()
		f.StringP("search", "s", "", "Volltextsuche")
		f.StringP("title", "t", "", "Titelsuche")
		f.String("from", "", "Datum von (JJJJ-MM-TT)")
		f.String("to", "", "Datum bis (JJJJ-MM-TT)")
		f.String("since", "", "Zeitfilter")
		f.String("sort-dir", "", "Sortierrichtung: asc, desc")
	}

	// App-specific flags.
	mrpCmd.Flags().String("submitter", "", "Einbringer/Ministerium")
	mrpCmd.Flags().String("session", "", "Sitzungsnummer")
	mrpCmd.Flags().String("period", "", "Gesetzgebungsperiode")
	mrpCmd.Flags().String("file-number", "", "Geschäftszahl")

	erlCmd.Flags().String("ministry", "", "Bundesministerium")
	erlCmd.Flags().String("department", "", "Abteilung")
	erlCmd.Flags().String("source", "", "Fundstelle")
	erlCmd.Flags().String("norm", "", "Norm")
	erlCmd.Flags().String("date", "", "Fassungsdatum (JJJJ-MM-TT)")

	uptsCmd.Flags().String("party", "", "Politische Partei: spoe, oevp, fpoe, gruene, neos, bzoe")
	uptsCmd.Flags().String("file-number", "", "Geschäftszahl")
	uptsCmd.Flags().String("norm", "", "Norm")

	kmgerCmd.Flags().String("type", "", "KmGer-Typ: geschaeftsordnung, geschaeftsverteilung")
	kmgerCmd.Flags().String("court-name", "", "Gericht")
	kmgerCmd.Flags().String("file-number", "", "Geschäftszahl")

	avsvCmd.Flags().String("doc-type", "", "Dokumentart")
	avsvCmd.Flags().String("author", "", "Urheber/Institution: dvsv, pva, oegk, auva, svs, bvaeb")
	avsvCmd.Flags().String("avsv-number", "", "AVSV-Nummer")

	avnCmd.Flags().String("avn-number", "", "AVN-Nummer")
	avnCmd.Flags().String("type", "", "AVN-Typ: kundmachung, verordnung, erlass")

	spgCmd.Flags().String("spg-number", "", "SPG-Nummer")
	spgCmd.Flags().String("osg-type", "", "OSG-Typ: oesg, oesg-grossgeraete")
	spgCmd.Flags().String("rsg-type", "", "RSG-Typ: rsg, rsg-grossgeraete")
	spgCmd.Flags().String("rsg-state", "", "Bundesland für RSG")

	pruefgewoCmd.Flags().String("type", "", "PrüfGewO-Typ: befaehigung, eignung, meister")

	sonstigeCmd.AddCommand(mrpCmd, erlCmd, uptsCmd, kmgerCmd, avsvCmd, avnCmd, spgCmd, pruefgewoCmd)
	rootCmd.AddCommand(sonstigeCmd)
}

// setCommonSonstigeParams sets the common parameters shared by all sonstige sub-commands.
func setCommonSonstigeParams(cmd *cobra.Command, params *api.Params) error {
	search, _ := cmd.Flags().GetString("search")
	title, _ := cmd.Flags().GetString("title")
	since, _ := cmd.Flags().GetString("since")
	sortDir, _ := cmd.Flags().GetString("sort-dir")

	if search != "" {
		params.Set("Suchworte", search)
	}
	if title != "" {
		params.Set("Titel", title)
	}
	if since != "" {
		value, ok := constants.ImRisSeit[strings.ToLower(since)]
		if !ok {
			return errValidation("Fehler: ungültiger --since Wert %q", since)
		}
		params.Set("ImRisSeit", value)
	}
	if sortDir != "" {
		value, ok := constants.SortDirections[strings.ToLower(sortDir)]
		if !ok {
			return errValidation("Fehler: ungültiger --sort-dir Wert %q (gültig: asc, desc)", sortDir)
		}
		params.Set("Sortierung.SortDirection", value)
	}
	return nil
}

func executeSonstigeSearch(cmd *cobra.Command, params *api.Params) error {
	return executeSearch(cmd, "Sonstige", "Suche in Sonstige Rechtsquellen...", params)
}

func runMrp(cmd *cobra.Command, args []string) error {
	params := api.NewParams()
	params.Set("Applikation", "Mrp")
	if err := setCommonSonstigeParams(cmd, params); err != nil {
		return err
	}

	from, _ := cmd.Flags().GetString("from")
	to, _ := cmd.Flags().GetString("to")
	submitter, _ := cmd.Flags().GetString("submitter")
	session, _ := cmd.Flags().GetString("session")
	period, _ := cmd.Flags().GetString("period")
	fileNumber, _ := cmd.Flags().GetString("file-number")

	if from != "" {
		params.Set("Sitzungsdatum.Von", from)
	}
	if to != "" {
		params.Set("Sitzungsdatum.Bis", to)
	}
	if submitter != "" {
		params.Set("Einbringer", submitter)
	}
	if session != "" {
		params.Set("Sitzungsnummer", session)
	}
	if period != "" {
		params.Set("Gesetzgebungsperiode", period)
	}
	if fileNumber != "" {
		params.Set("Geschaeftszahl", fileNumber)
	}

	return executeSonstigeSearch(cmd, params)
}

func runErlaesse(cmd *cobra.Command, args []string) error {
	params := api.NewParams()
	params.Set("Applikation", "Erlaesse")
	if err := setCommonSonstigeParams(cmd, params); err != nil {
		return err
	}

	from, _ := cmd.Flags().GetString("from")
	to, _ := cmd.Flags().GetString("to")
	ministry, _ := cmd.Flags().GetString("ministry")
	department, _ := cmd.Flags().GetString("department")
	source, _ := cmd.Flags().GetString("source")
	norm, _ := cmd.Flags().GetString("norm")
	date, _ := cmd.Flags().GetString("date")

	if from != "" {
		params.Set("VonInkrafttretensdatum", from)
	}
	if to != "" {
		params.Set("BisInkrafttretensdatum", to)
	}
	if ministry != "" {
		value, ok := constants.ErlMinistries[strings.ToLower(ministry)]
		if !ok {
			return errValidation("Fehler: ungültiger --ministry Wert %q\nGültig: bka, bmkoes, bmeia, bmaw, bmbwf, bmf, bmi, bmj, bmk, bmlv, bml, bmsgpk", ministry)
		}
		params.Set("Bundesministerium", value)
	}
	if department != "" {
		params.Set("Abteilung", department)
	}
	if source != "" {
		params.Set("Fundstelle", source)
	}
	if norm != "" {
		params.Set("Norm", norm)
	}
	if date != "" {
		params.Set("FassungVom", date)
	}

	return executeSonstigeSearch(cmd, params)
}

func runUpts(cmd *cobra.Command, args []string) error {
	params := api.NewParams()
	params.Set("Applikation", "Upts")
	if err := setCommonSonstigeParams(cmd, params); err != nil {
		return err
	}

	from, _ := cmd.Flags().GetString("from")
	to, _ := cmd.Flags().GetString("to")
	party, _ := cmd.Flags().GetString("party")
	fileNumber, _ := cmd.Flags().GetString("file-number")
	norm, _ := cmd.Flags().GetString("norm")

	if from != "" {
		params.Set("Entscheidungsdatum.Von", from)
	}
	if to != "" {
		params.Set("Entscheidungsdatum.Bis", to)
	}
	if party != "" {
		value, ok := constants.UptsParties[strings.ToLower(party)]
		if !ok {
			return errValidation("Fehler: ungültiger --party Wert %q\nGültig: spoe, oevp, fpoe, gruene, neos, bzoe", party)
		}
		params.Set("Partei", value)
	}
	if fileNumber != "" {
		params.Set("Geschaeftszahl", fileNumber)
	}
	if norm != "" {
		params.Set("Norm", norm)
	}

	return executeSonstigeSearch(cmd, params)
}

func runKmger(cmd *cobra.Command, args []string) error {
	params := api.NewParams()
	params.Set("Applikation", "KmGer")
	if err := setCommonSonstigeParams(cmd, params); err != nil {
		return err
	}

	from, _ := cmd.Flags().GetString("from")
	to, _ := cmd.Flags().GetString("to")
	typ, _ := cmd.Flags().GetString("type")
	courtName, _ := cmd.Flags().GetString("court-name")
	fileNumber, _ := cmd.Flags().GetString("file-number")

	if from != "" {
		params.Set("Kundmachungsdatum.Von", from)
	}
	if to != "" {
		params.Set("Kundmachungsdatum.Bis", to)
	}
	if typ != "" {
		value, ok := constants.KmgerTypes[strings.ToLower(typ)]
		if !ok {
			return errValidation("Fehler: ungültiger --type Wert %q (gültig: geschaeftsordnung, geschaeftsverteilung)", typ)
		}
		params.Set("Typ", value)
	}
	if courtName != "" {
		params.Set("Gericht", courtName)
	}
	if fileNumber != "" {
		params.Set("Geschaeftszahl", fileNumber)
	}

	return executeSonstigeSearch(cmd, params)
}

func runAvsv(cmd *cobra.Command, args []string) error {
	params := api.NewParams()
	params.Set("Applikation", "Avsv")
	if err := setCommonSonstigeParams(cmd, params); err != nil {
		return err
	}

	from, _ := cmd.Flags().GetString("from")
	to, _ := cmd.Flags().GetString("to")
	docType, _ := cmd.Flags().GetString("doc-type")
	author, _ := cmd.Flags().GetString("author")
	avsvNumber, _ := cmd.Flags().GetString("avsv-number")

	if from != "" {
		params.Set("Kundmachung.Von", from)
	}
	if to != "" {
		params.Set("Kundmachung.Bis", to)
	}
	if docType != "" {
		params.Set("Dokumentart", docType)
	}
	if author != "" {
		value, ok := constants.AvsvAuthors[strings.ToLower(author)]
		if !ok {
			return errValidation("Fehler: ungültiger --author Wert %q\nGültig: dvsv, pva, oegk, auva, svs, bvaeb", author)
		}
		params.Set("Urheber", value)
	}
	if avsvNumber != "" {
		params.Set("Avsvnummer", avsvNumber)
	}

	return executeSonstigeSearch(cmd, params)
}

func runAvn(cmd *cobra.Command, args []string) error {
	params := api.NewParams()
	params.Set("Applikation", "Avn")
	if err := setCommonSonstigeParams(cmd, params); err != nil {
		return err
	}

	from, _ := cmd.Flags().GetString("from")
	to, _ := cmd.Flags().GetString("to")
	avnNumber, _ := cmd.Flags().GetString("avn-number")
	typ, _ := cmd.Flags().GetString("type")

	if from != "" {
		params.Set("Kundmachung.Von", from)
	}
	if to != "" {
		params.Set("Kundmachung.Bis", to)
	}
	if avnNumber != "" {
		params.Set("Avnnummer", avnNumber)
	}
	if typ != "" {
		value, ok := constants.AvnTypes[strings.ToLower(typ)]
		if !ok {
			return errValidation("Fehler: ungültiger --type Wert %q (gültig: kundmachung, verordnung, erlass)", typ)
		}
		params.Set("Typ", value)
	}

	return executeSonstigeSearch(cmd, params)
}

func runSpg(cmd *cobra.Command, args []string) error {
	params := api.NewParams()
	params.Set("Applikation", "Spg")
	if err := setCommonSonstigeParams(cmd, params); err != nil {
		return err
	}

	from, _ := cmd.Flags().GetString("from")
	to, _ := cmd.Flags().GetString("to")
	spgNumber, _ := cmd.Flags().GetString("spg-number")
	osgType, _ := cmd.Flags().GetString("osg-type")
	rsgType, _ := cmd.Flags().GetString("rsg-type")
	rsgState, _ := cmd.Flags().GetString("rsg-state")

	if from != "" {
		params.Set("Kundmachungsdatum.Von", from)
	}
	if to != "" {
		params.Set("Kundmachungsdatum.Bis", to)
	}
	if spgNumber != "" {
		params.Set("Spgnummer", spgNumber)
	}
	if osgType != "" {
		value, ok := constants.OsgTypes[strings.ToLower(osgType)]
		if !ok {
			return errValidation("Fehler: ungültiger --osg-type Wert %q (gültig: oesg, oesg-grossgeraete)", osgType)
		}
		params.Set("OsgTyp", value)
	}
	if rsgType != "" {
		value, ok := constants.RsgTypes[strings.ToLower(rsgType)]
		if !ok {
			return errValidation("Fehler: ungültiger --rsg-type Wert %q (gültig: rsg, rsg-grossgeraete)", rsgType)
		}
		params.Set("RsgTyp", value)
	}
	if rsgState != "" {
		params.Set("RsgLand", rsgState)
	}

	return executeSonstigeSearch(cmd, params)
}

func runPruefgewo(cmd *cobra.Command, args []string) error {
	params := api.NewParams()
	params.Set("Applikation", "PruefGewO")
	if err := setCommonSonstigeParams(cmd, params); err != nil {
		return err
	}

	from, _ := cmd.Flags().GetString("from")
	to, _ := cmd.Flags().GetString("to")
	typ, _ := cmd.Flags().GetString("type")

	if from != "" {
		params.Set("Kundmachungsdatum.Von", from)
	}
	if to != "" {
		params.Set("Kundmachungsdatum.Bis", to)
	}
	if typ != "" {
		value, ok := constants.PruefgewoTypes[strings.ToLower(typ)]
		if !ok {
			return errValidation("Fehler: ungültiger --type Wert %q (gültig: befaehigung, eignung, meister)", typ)
		}
		params.Set("Typ", value)
	}

	return executeSonstigeSearch(cmd, params)
}
