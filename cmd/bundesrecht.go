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

var bundesrechtCmd = &cobra.Command{
	Use:   "bundesrecht",
	Short: "Bundesgesetze durchsuchen (ABGB, StGB, etc.)",
	Long: `Österreichische Bundesgesetze (Bundesrecht) durchsuchen.

Beispiele:
  ris bundesrecht --search "Mietrecht"
  ris bundesrecht --title "ABGB" --paragraph 1295
  ris bundesrecht --search "Schadenersatz" --app begut
  ris bundesrecht --search "Mietrecht" --date 2024-01-15 --json`,
	RunE: runBundesrecht,
}

func init() {
	f := bundesrechtCmd.Flags()
	f.StringP("search", "s", "", "Volltextsuche")
	f.StringP("title", "t", "", "Suche in Gesetzestitel")
	f.String("paragraph", "", "Paragraphennummer (z.B. \"1295\")")
	f.String("app", "brkons", "Applikation: brkons, begut, bgblauth, erv")
	f.String("date", "", "Fassungsdatum (JJJJ-MM-TT)")

	rootCmd.AddCommand(bundesrechtCmd)
}

func runBundesrecht(cmd *cobra.Command, args []string) error {
	search, _ := cmd.Flags().GetString("search")
	title, _ := cmd.Flags().GetString("title")
	paragraph, _ := cmd.Flags().GetString("paragraph")
	app, _ := cmd.Flags().GetString("app")
	date, _ := cmd.Flags().GetString("date")

	// At least one of search/title/paragraph required.
	if search == "" && title == "" && paragraph == "" {
		fmt.Fprintln(os.Stderr, "Fehler: mindestens --search, --title oder --paragraph erforderlich")
		os.Exit(2)
	}

	// Resolve application value.
	appValue, ok := constants.BundesrechtApps[strings.ToLower(app)]
	if !ok {
		fmt.Fprintf(os.Stderr, "Fehler: ungültiger --app Wert %q (gültig: brkons, begut, bgblauth, erv)\n", app)
		os.Exit(2)
	}

	client := newClient(cmd)
	params := api.NewParams()
	params.Set("Applikation", appValue)

	if search != "" {
		params.Set("Suchworte", search)
	}
	if title != "" {
		params.Set("Titel", title)
	}
	if paragraph != "" {
		params.Set("Abschnitt.Von", paragraph)
		params.Set("Abschnitt.Bis", paragraph)
		params.Set("Abschnitt.Typ", "Paragraph")
	}
	if date != "" {
		params.Set("FassungVom", date)
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
