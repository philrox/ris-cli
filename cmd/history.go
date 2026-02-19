package cmd

import (
	"strings"

	"github.com/philrox/ris-cli/internal/api"
	"github.com/philrox/ris-cli/internal/constants"
	"github.com/spf13/cobra"
)

var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "Dokumentänderungshistorie durchsuchen",
	Long: `Änderungshistorie von Dokumenten durchsuchen.

Beispiele:
  ris history --app bundesnormen --from 2024-01-01 --to 2024-01-31
  ris history --app justiz --from 2024-06-01 --include-deleted`,
	RunE: runHistory,
}

func init() {
	f := historyCmd.Flags()
	f.StringP("app", "a", "", "Anwendung (erforderlich)")
	f.String("from", "", "Änderungen von (JJJJ-MM-TT)")
	f.String("to", "", "Änderungen bis (JJJJ-MM-TT)")
	f.Bool("include-deleted", false, "Gelöschte Dokumente einschließen")

	rootCmd.AddCommand(historyCmd)
}

func runHistory(cmd *cobra.Command, args []string) error {
	app, _ := cmd.Flags().GetString("app")
	from, _ := cmd.Flags().GetString("from")
	to, _ := cmd.Flags().GetString("to")
	includeDeleted, _ := cmd.Flags().GetBool("include-deleted")

	// --app is required.
	if app == "" {
		return errValidation("Fehler: --app ist erforderlich")
	}

	// At least one of --from or --to required.
	if from == "" && to == "" {
		return errValidation("Fehler: mindestens --from oder --to erforderlich")
	}

	// Validate app value.
	appLower := strings.ToLower(app)
	if !constants.IsValidHistoryApp(appLower) {
		return errValidation("Fehler: ungültiger --app Wert %q\nGültig: bundesnormen, landesnormen, justiz, vfgh, vwgh, bvwg, lvwg, bgblauth, bgblalt, bgblpdf, lgblauth, lgbl, lgblno, gemeinderecht, gemeinderechtauth, bvb, vbl, regv, mrp, erlaesse, pruefgewo, avsv, spg, kmger, dsk, gbk, dok, pvak, normenliste, asylgh", app)
	}

	params := api.NewParams()

	// History uses Anwendung, NOT Applikation.
	params.Set("Anwendung", appLower)

	if from != "" {
		params.Set("AenderungenVon", from)
	}
	if to != "" {
		params.Set("AenderungenBis", to)
	}
	if includeDeleted {
		params.Set("IncludeDeletedDocuments", "true")
	}

	return executeSearch(cmd, "History", "Suche in Änderungshistorie...", params)
}
