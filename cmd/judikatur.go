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

var judikaturCmd = &cobra.Command{
	Use:   "judikatur",
	Short: "Gerichtsentscheidungen durchsuchen",
	Long: `Österreichische Gerichtsentscheidungen durchsuchen.

Beispiele:
  ris judikatur --search "Grundrecht" --court vfgh
  ris judikatur --case-number "5Ob234/20b"
  ris judikatur --norm "1319a ABGB" --from 2020-01-01 --to 2024-12-31`,
	RunE: runJudikatur,
}

func init() {
	f := judikaturCmd.Flags()
	f.StringP("search", "s", "", "Volltextsuche")
	f.StringP("norm", "n", "", "Normverweis")
	f.String("case-number", "", "Geschäftszahl")
	f.StringP("court", "c", "justiz", "Gerichtstyp: justiz, vfgh, vwgh, bvwg, lvwg, dsk, asylgh, normenliste, pvak, gbk, dok")
	f.String("from", "", "Entscheidungsdatum von (JJJJ-MM-TT)")
	f.String("to", "", "Entscheidungsdatum bis (JJJJ-MM-TT)")

	rootCmd.AddCommand(judikaturCmd)
}

func runJudikatur(cmd *cobra.Command, args []string) error {
	search, _ := cmd.Flags().GetString("search")
	norm, _ := cmd.Flags().GetString("norm")
	caseNumber, _ := cmd.Flags().GetString("case-number")
	court, _ := cmd.Flags().GetString("court")
	from, _ := cmd.Flags().GetString("from")
	to, _ := cmd.Flags().GetString("to")

	// At least one of search/norm/case-number required.
	if search == "" && norm == "" && caseNumber == "" {
		return errValidation("Fehler: mindestens --search, --norm oder --case-number erforderlich")
	}

	// Resolve court to Applikation value.
	courtValue, ok := constants.Courts[strings.ToLower(court)]
	if !ok {
		return errValidation("Fehler: ungültiger --court Wert %q\nGültige Gerichte: justiz, vfgh, vwgh, bvwg, lvwg, dsk, asylgh, normenliste, pvak, gbk, dok", court)
	}

	client := newClient(cmd)
	params := api.NewParams()
	params.Set("Applikation", courtValue)

	if search != "" {
		params.Set("Suchworte", search)
	}
	if norm != "" {
		params.Set("Norm", norm)
	}
	if caseNumber != "" {
		params.Set("Geschaeftszahl", caseNumber)
	}
	if from != "" {
		params.Set("EntscheidungsdatumVon", from)
	}
	if to != "" {
		params.Set("EntscheidungsdatumBis", to)
	}

	setPageParams(cmd, params)

	s := startSpinner(cmd, "Suche in Judikatur...")
	body, err := client.Search("Judikatur", params)
	stopSpinner(s)
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
