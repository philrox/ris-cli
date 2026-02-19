package cmd

import (
	"strings"

	"github.com/philrox/risgo/internal/api"
	"github.com/philrox/risgo/internal/constants"
	"github.com/spf13/cobra"
)

var bgblCmd = &cobra.Command{
	Use:   "bgbl",
	Short: "Bundesgesetzblätter durchsuchen",
	Long: `Bundesgesetzblätter (BGBl) durchsuchen.

Beispiele:
  risgo bgbl --number 120 --year 2023 --part 1
  risgo bgbl --search "Klimaschutz" --json`,
	RunE: runBgbl,
}

func init() {
	f := bgblCmd.Flags()
	f.String("number", "", "BGBl-Nummer")
	f.String("year", "", "Jahrgang")
	f.StringP("search", "s", "", "Volltextsuche")
	f.StringP("title", "t", "", "Titelsuche")
	f.String("part", "", "Teil: 1 (Gesetze), 2 (Verordnungen), 3 (Staatsverträge)")
	f.String("app", "bgblauth", "Applikation: bgblauth, bgblpdf, bgblalt")

	rootCmd.AddCommand(bgblCmd)
}

func runBgbl(cmd *cobra.Command, args []string) error {
	number, _ := cmd.Flags().GetString("number")
	year, _ := cmd.Flags().GetString("year")
	search, _ := cmd.Flags().GetString("search")
	title, _ := cmd.Flags().GetString("title")
	part, _ := cmd.Flags().GetString("part")
	app, _ := cmd.Flags().GetString("app")

	// At least one of number/year/search/title required.
	if number == "" && year == "" && search == "" && title == "" {
		return errValidation("Fehler: mindestens --number, --year, --search oder --title erforderlich")
	}

	appValue, ok := constants.BgblApps[strings.ToLower(app)]
	if !ok {
		return errValidation("Fehler: ungültiger --app Wert %q (gültig: bgblauth, bgblpdf, bgblalt)", app)
	}

	params := api.NewParams()
	params.Set("Applikation", appValue)

	if number != "" {
		params.Set("Bgblnummer", number)
	}
	if year != "" {
		params.Set("Jahrgang", year)
	}
	if search != "" {
		params.Set("Suchworte", search)
	}
	if title != "" {
		params.Set("Titel", title)
	}
	if part != "" {
		teilValue, ok := constants.BgblTeile[part]
		if !ok {
			return errValidation("Fehler: ungültiger --part Wert %q (gültig: 1, 2, 3)", part)
		}
		params.Set("Teil", teilValue)
	}

	return executeSearch(cmd, "Bundesrecht", "Suche in Bundesgesetzblättern...", params)
}
