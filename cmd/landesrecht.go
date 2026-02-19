package cmd

import (
	"strings"

	"github.com/philrox/ris-cli/internal/api"
	"github.com/philrox/ris-cli/internal/constants"
	"github.com/spf13/cobra"
)

var landesrechtCmd = &cobra.Command{
	Use:   "landesrecht",
	Short: "Landesgesetze durchsuchen",
	Long: `Österreichische Landesgesetze (Landesrecht) durchsuchen.

Beispiele:
  ris landesrecht --search "Bauordnung" --state salzburg
  ris landesrecht --title "Raumordnung" --state wien --json`,
	RunE: runLandesrecht,
}

func init() {
	f := landesrechtCmd.Flags()
	f.StringP("search", "s", "", "Volltextsuche")
	f.StringP("title", "t", "", "Suche in Gesetzestitel")
	f.String("state", "", "Bundesland (z.B. wien, salzburg, tirol)")

	rootCmd.AddCommand(landesrechtCmd)
}

func runLandesrecht(cmd *cobra.Command, args []string) error {
	search, _ := cmd.Flags().GetString("search")
	title, _ := cmd.Flags().GetString("title")
	state, _ := cmd.Flags().GetString("state")

	// At least one required.
	if search == "" && title == "" && state == "" {
		return errValidation("Fehler: mindestens --search, --title oder --state erforderlich")
	}

	params := api.NewParams()
	params.Set("Applikation", "LrKons")

	if search != "" {
		params.Set("Suchworte", search)
	}
	if title != "" {
		params.Set("Titel", title)
	}

	if state != "" {
		paramName, ok := constants.LandesrechtStates[strings.ToLower(state)]
		if !ok {
			return errValidation("Fehler: ungültiger --state Wert %q\nGültige Bundesländer: wien, niederoesterreich, oberoesterreich, salzburg, tirol, vorarlberg, kaernten, steiermark, burgenland", state)
		}
		params.Set(paramName, "true")
	}

	return executeSearch(cmd, "Landesrecht", "Suche in Landesrecht...", params)
}
