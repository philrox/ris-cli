package cmd

import (
	"strings"

	"github.com/philrox/risgo/internal/api"
	"github.com/philrox/risgo/internal/constants"
	"github.com/spf13/cobra"
)

var lgblCmd = &cobra.Command{
	Use:   "lgbl",
	Short: "Landesgesetzblätter durchsuchen",
	Long: `Landesgesetzblätter (LGBl) durchsuchen.

Beispiele:
  risgo lgbl --number 50 --year 2023 --state wien
  risgo lgbl --search "Bauordnung" --state salzburg`,
	RunE: runLgbl,
}

func init() {
	f := lgblCmd.Flags()
	f.String("number", "", "LGBl-Nummer")
	f.String("year", "", "Jahrgang")
	f.String("state", "", "Bundesland (z.B. wien, salzburg)")
	f.StringP("search", "s", "", "Volltextsuche")
	f.StringP("title", "t", "", "Titelsuche")
	f.String("app", "lgblauth", "Applikation: lgblauth, lgbl, lgblno")

	rootCmd.AddCommand(lgblCmd)
}

func runLgbl(cmd *cobra.Command, args []string) error {
	number, _ := cmd.Flags().GetString("number")
	year, _ := cmd.Flags().GetString("year")
	state, _ := cmd.Flags().GetString("state")
	search, _ := cmd.Flags().GetString("search")
	title, _ := cmd.Flags().GetString("title")
	app, _ := cmd.Flags().GetString("app")

	// At least one required.
	if number == "" && year == "" && state == "" && search == "" && title == "" {
		return errValidation("Fehler: mindestens --number, --year, --state, --search oder --title erforderlich")
	}

	appValue, ok := constants.LgblApps[strings.ToLower(app)]
	if !ok {
		return errValidation("Fehler: ungültiger --app Wert %q (gültig: lgblauth, lgbl, lgblno)", app)
	}

	params := api.NewParams()
	params.Set("Applikation", appValue)

	if number != "" {
		params.Set("Lgblnummer", number)
	}
	if year != "" {
		params.Set("Jahrgang", year)
	}
	if state != "" {
		paramName, ok := constants.LandesrechtStates[strings.ToLower(state)]
		if !ok {
			return errValidation("Fehler: ungültiger --state Wert %q\nGültige Bundesländer: wien, niederoesterreich, oberoesterreich, salzburg, tirol, vorarlberg, kaernten, steiermark, burgenland", state)
		}
		params.Set(paramName, "true")
	}
	if search != "" {
		params.Set("Suchworte", search)
	}
	if title != "" {
		params.Set("Titel", title)
	}

	return executeSearch(cmd, "Landesrecht", "Suche in Landesgesetzblättern...", params)
}
