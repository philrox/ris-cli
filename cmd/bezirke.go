package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/philrox/ris-cli/internal/api"
	"github.com/philrox/ris-cli/internal/constants"
	"github.com/philrox/ris-cli/internal/format"
	"github.com/philrox/ris-cli/internal/parser"
	"github.com/spf13/cobra"
)

var bezirkeCmd = &cobra.Command{
	Use:   "bezirke",
	Short: "Bezirksverwaltungsbehörden-Kundmachungen durchsuchen",
	Long: `Kundmachungen der Bezirksverwaltungsbehörden durchsuchen.

Beispiele:
  ris bezirke --state niederoesterreich --search "Bauordnung"
  ris bezirke --authority "Bezirkshauptmannschaft Innsbruck"`,
	RunE: runBezirke,
}

func init() {
	f := bezirkeCmd.Flags()
	f.StringP("search", "s", "", "Volltextsuche")
	f.StringP("title", "t", "", "Titelsuche")
	f.String("state", "", "Bundesland")
	f.String("authority", "", "Bezirksverwaltungsbehörde")
	f.String("number", "", "Kundmachungsnummer")
	f.String("from", "", "Datum von (JJJJ-MM-TT)")
	f.String("to", "", "Datum bis (JJJJ-MM-TT)")
	f.String("since", "", "Zeitfilter")

	rootCmd.AddCommand(bezirkeCmd)
}

func runBezirke(cmd *cobra.Command, args []string) error {
	search, _ := cmd.Flags().GetString("search")
	title, _ := cmd.Flags().GetString("title")
	state, _ := cmd.Flags().GetString("state")
	authority, _ := cmd.Flags().GetString("authority")
	number, _ := cmd.Flags().GetString("number")
	from, _ := cmd.Flags().GetString("from")
	to, _ := cmd.Flags().GetString("to")
	since, _ := cmd.Flags().GetString("since")

	// At least one required.
	if search == "" && title == "" && state == "" && authority == "" && number == "" {
		fmt.Fprintln(os.Stderr, "Fehler: mindestens --search, --title, --state, --authority oder --number erforderlich")
		os.Exit(2)
	}

	client := newClient(cmd)
	params := api.NewParams()
	params.Set("Applikation", "Bvb")

	if search != "" {
		params.Set("Suchworte", search)
	}
	if title != "" {
		params.Set("Titel", title)
	}
	if state != "" {
		// Bezirke uses display names with Umlauts.
		value, ok := constants.BezirkeStates[strings.ToLower(state)]
		if !ok {
			fmt.Fprintf(os.Stderr, "Fehler: ungültiger --state Wert %q\n", state)
			fmt.Fprintln(os.Stderr, "Gültig: wien, niederoesterreich, oberoesterreich, salzburg, tirol, vorarlberg, kaernten, steiermark, burgenland")
			os.Exit(2)
		}
		params.Set("Bundesland", value)
	}
	if authority != "" {
		params.Set("Bezirksverwaltungsbehoerde", authority)
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
	if since != "" {
		value, ok := constants.ImRisSeit[strings.ToLower(since)]
		if !ok {
			fmt.Fprintf(os.Stderr, "Fehler: ungültiger --since Wert %q\n", since)
			os.Exit(2)
		}
		params.Set("ImRisSeit", value)
	}

	setPageParams(cmd, params)

	body, err := client.Search("Bezirke", params)
	if err != nil {
		return fmt.Errorf("API-Anfrage fehlgeschlagen: %w", err)
	}

	result, err := parser.ParseSearchResponse(body)
	if err != nil {
		return fmt.Errorf("Antwort konnte nicht verarbeitet werden: %w", err)
	}

	if useJSON(cmd) {
		return format.JSON(os.Stdout, result)
	}
	return format.Text(os.Stdout, result)
}
