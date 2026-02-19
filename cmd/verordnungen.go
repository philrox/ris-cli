package cmd

import (
	"strings"

	"github.com/philrox/ris-cli/internal/api"
	"github.com/philrox/ris-cli/internal/constants"
	"github.com/spf13/cobra"
)

var verordnungenCmd = &cobra.Command{
	Use:   "verordnungen",
	Short: "Verordnungsblätter durchsuchen",
	Long: `Verordnungsblätter der Länder durchsuchen.

Beispiele:
  ris verordnungen --search "Wolf" --state tirol
  ris verordnungen --number 25 --from 2024-01-01`,
	RunE: runVerordnungen,
}

func init() {
	f := verordnungenCmd.Flags()
	f.StringP("search", "s", "", "Volltextsuche")
	f.StringP("title", "t", "", "Titelsuche")
	f.String("state", "", "Bundesland")
	f.String("number", "", "Kundmachungsnummer")
	f.String("from", "", "Datum von (JJJJ-MM-TT)")
	f.String("to", "", "Datum bis (JJJJ-MM-TT)")

	rootCmd.AddCommand(verordnungenCmd)
}

func runVerordnungen(cmd *cobra.Command, args []string) error {
	search, _ := cmd.Flags().GetString("search")
	title, _ := cmd.Flags().GetString("title")
	state, _ := cmd.Flags().GetString("state")
	number, _ := cmd.Flags().GetString("number")
	from, _ := cmd.Flags().GetString("from")
	to, _ := cmd.Flags().GetString("to")

	// At least one required.
	if search == "" && title == "" && state == "" && number == "" && from == "" {
		return errValidation("Fehler: mindestens --search, --title, --state, --number oder --from erforderlich")
	}

	params := api.NewParams()
	params.Set("Applikation", "Vbl")

	if search != "" {
		params.Set("Suchworte", search)
	}
	if title != "" {
		params.Set("Titel", title)
	}
	if state != "" {
		// Verordnungen uses direct Bundesland values, NOT SucheIn* format.
		value, ok := constants.VerordnungenStates[strings.ToLower(state)]
		if !ok {
			return errValidation("Fehler: ungültiger --state Wert %q\nGültige Bundesländer: wien, niederoesterreich, oberoesterreich, salzburg, tirol, vorarlberg, kaernten, steiermark, burgenland", state)
		}
		params.Set("Bundesland", value)
	}
	if number != "" {
		params.Set("Kundmachungsnummer", number)
	}
	if from != "" {
		params.Set("Kundmachungsdatum.Von", from)
	}
	if to != "" {
		params.Set("Kundmachungsdatum.Bis", to)
	}

	return executeSearch(cmd, "Landesrecht", "Suche in Verordnungsblättern...", params)
}
