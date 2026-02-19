package cmd

import (
	"strings"

	"github.com/philrox/ris-cli/internal/api"
	"github.com/philrox/ris-cli/internal/constants"
	"github.com/spf13/cobra"
)

var gemeindenCmd = &cobra.Command{
	Use:   "gemeinden",
	Short: "Gemeinderecht durchsuchen",
	Long: `Österreichisches Gemeinderecht durchsuchen.

Beispiele:
  ris gemeinden --municipality "Graz" --search "Parkgebuehren"
  ris gemeinden --state tirol --title "Gebuehrenordnung"`,
	RunE: runGemeinden,
}

func init() {
	f := gemeindenCmd.Flags()
	f.StringP("search", "s", "", "Volltextsuche")
	f.StringP("title", "t", "", "Titelsuche")
	f.String("state", "", "Bundesland")
	f.String("municipality", "", "Gemeindename")
	f.String("file-number", "", "Geschäftszahl (nur Gr)")
	f.String("index", "", "Sachbereichsindex (nur Gr)")
	f.String("district", "", "Bezirk (nur GrA)")
	f.String("gemeindeverband", "", "Gemeindeverband (nur GrA)")
	f.String("announcement-nr", "", "Kundmachungsnummer (nur GrA)")
	f.String("app", "gr", "Applikation: gr, gra")
	f.String("date", "", "Fassungsdatum (nur Gr, JJJJ-MM-TT)")
	f.String("from", "", "Datum von (nur GrA, JJJJ-MM-TT)")
	f.String("to", "", "Datum bis (nur GrA, JJJJ-MM-TT)")
	f.String("since", "", "Zeitfilter")
	f.String("sort-dir", "", "Sortierrichtung: asc, desc")
	f.String("sort-by", "", "Sortierspalte (nur Gr): geschaeftszahl, bundesland, gemeinde")

	rootCmd.AddCommand(gemeindenCmd)
}

func runGemeinden(cmd *cobra.Command, args []string) error {
	search, _ := cmd.Flags().GetString("search")
	title, _ := cmd.Flags().GetString("title")
	state, _ := cmd.Flags().GetString("state")
	municipality, _ := cmd.Flags().GetString("municipality")
	fileNumber, _ := cmd.Flags().GetString("file-number")
	index, _ := cmd.Flags().GetString("index")
	district, _ := cmd.Flags().GetString("district")
	gemeindeverband, _ := cmd.Flags().GetString("gemeindeverband")
	announcementNr, _ := cmd.Flags().GetString("announcement-nr")
	app, _ := cmd.Flags().GetString("app")
	date, _ := cmd.Flags().GetString("date")
	from, _ := cmd.Flags().GetString("from")
	to, _ := cmd.Flags().GetString("to")
	since, _ := cmd.Flags().GetString("since")
	sortDir, _ := cmd.Flags().GetString("sort-dir")
	sortBy, _ := cmd.Flags().GetString("sort-by")

	// At least one required.
	if search == "" && title == "" && state == "" && municipality == "" &&
		fileNumber == "" && index == "" && district == "" && gemeindeverband == "" && announcementNr == "" {
		return errValidation("Fehler: mindestens ein Suchparameter erforderlich")
	}

	appValue, ok := constants.GemeindenApps[strings.ToLower(app)]
	if !ok {
		return errValidation("Fehler: ungültiger --app Wert %q (gültig: gr, gra)", app)
	}

	params := api.NewParams()
	params.Set("Applikation", appValue)

	if search != "" {
		params.Set("Suchworte", search)
	}
	if title != "" {
		params.Set("Titel", title)
	}
	if state != "" {
		params.Set("Bundesland", state)
	}
	if municipality != "" {
		params.Set("Gemeinde", municipality)
	}
	if fileNumber != "" {
		params.Set("Geschaeftszahl", fileNumber)
	}
	if index != "" {
		value, ok := constants.GemeindenIndex[strings.ToLower(index)]
		if !ok {
			return errValidation("Fehler: ungültiger --index Wert %q", index)
		}
		params.Set("Index", value)
	}
	if district != "" {
		params.Set("Bezirk", district)
	}
	if gemeindeverband != "" {
		params.Set("Gemeindeverband", gemeindeverband)
	}
	if announcementNr != "" {
		params.Set("Kundmachungsnummer", announcementNr)
	}
	if date != "" {
		params.Set("FassungVom", date)
	}
	if from != "" {
		params.Set("Kundmachungsdatum.Von", from)
	}
	if to != "" {
		params.Set("Kundmachungsdatum.Bis", to)
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
	if sortBy != "" {
		value, ok := constants.GemeindenSortColumns[strings.ToLower(sortBy)]
		if !ok {
			return errValidation("Fehler: ungültiger --sort-by Wert %q (gültig: geschaeftszahl, bundesland, gemeinde)", sortBy)
		}
		params.Set("Sortierung.SortedByColumn", value)
	}

	return executeSearch(cmd, "Gemeinden", "Suche in Gemeinderecht...", params)
}
