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

var bgblCmd = &cobra.Command{
	Use:   "bgbl",
	Short: "Bundesgesetzblätter durchsuchen",
	Long: `Bundesgesetzblätter (BGBl) durchsuchen.

Beispiele:
  ris bgbl --number 120 --year 2023 --part 1
  ris bgbl --search "Klimaschutz" --json`,
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
		fmt.Fprintln(os.Stderr, "Fehler: mindestens --number, --year, --search oder --title erforderlich")
		os.Exit(2)
	}

	appValue, ok := constants.BgblApps[strings.ToLower(app)]
	if !ok {
		fmt.Fprintf(os.Stderr, "Fehler: ungültiger --app Wert %q (gültig: bgblauth, bgblpdf, bgblalt)\n", app)
		os.Exit(2)
	}

	client := newClient(cmd)
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
			fmt.Fprintf(os.Stderr, "Fehler: ungültiger --part Wert %q (gültig: 1, 2, 3)\n", part)
			os.Exit(2)
		}
		params.Set("Teil", teilValue)
	}

	setPageParams(cmd, params)

	body, err := client.Search("Bundesrecht", params)
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
