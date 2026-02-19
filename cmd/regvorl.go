package cmd

import (
	"strings"

	"github.com/philrox/risgo/internal/api"
	"github.com/philrox/risgo/internal/constants"
	"github.com/spf13/cobra"
)

var regvorlCmd = &cobra.Command{
	Use:   "regvorl",
	Short: "Regierungsvorlagen durchsuchen",
	Long: `Regierungsvorlagen durchsuchen.

Beispiele:
  risgo regvorl --search "Klimaschutz"
  risgo regvorl --ministry bmf --from 2024-01-01`,
	RunE: runRegvorl,
}

func init() {
	f := regvorlCmd.Flags()
	f.StringP("search", "s", "", "Volltextsuche")
	f.StringP("title", "t", "", "Titelsuche")
	f.String("from", "", "Beschlussdatum von (JJJJ-MM-TT)")
	f.String("to", "", "Beschlussdatum bis (JJJJ-MM-TT)")
	f.String("ministry", "", "Einbringendes Ministerium (z.B. bmf, bmi, bmj)")
	f.String("since", "", "Zeitfilter: einerwoche, zweiwochen, einemmonat, dreimonaten, sechsmonaten, einemjahr")
	f.String("sort-dir", "", "Sortierrichtung: asc, desc")
	f.String("sort-by", "", "Sortierspalte: kurztitel, stelle, datum")

	rootCmd.AddCommand(regvorlCmd)
}

func runRegvorl(cmd *cobra.Command, args []string) error {
	search, _ := cmd.Flags().GetString("search")
	title, _ := cmd.Flags().GetString("title")
	from, _ := cmd.Flags().GetString("from")
	to, _ := cmd.Flags().GetString("to")
	ministry, _ := cmd.Flags().GetString("ministry")
	since, _ := cmd.Flags().GetString("since")
	sortDir, _ := cmd.Flags().GetString("sort-dir")
	sortBy, _ := cmd.Flags().GetString("sort-by")

	// At least one required.
	if search == "" && title == "" && from == "" && ministry == "" && since == "" {
		return errValidation("Fehler: mindestens --search, --title, --from, --ministry oder --since erforderlich")
	}

	params := api.NewParams()
	params.Set("Applikation", "RegV")

	if search != "" {
		params.Set("Suchworte", search)
	}
	if title != "" {
		params.Set("Titel", title)
	}
	if from != "" {
		params.Set("BeschlussdatumVon", from)
	}
	if to != "" {
		params.Set("BeschlussdatumBis", to)
	}
	if ministry != "" {
		value, ok := constants.RegvorlMinistries[strings.ToLower(ministry)]
		if !ok {
			return errValidation("Fehler: ungültiger --ministry Wert %q\nGültig: bka, bmkoes, bmeia, bmaw, bmbwf, bmf, bmi, bmj, bmk, bmlv, bml, bmsgpk, bmffim, bmeuv", ministry)
		}
		params.Set("EinbringendeStelle", value)
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
		value, ok := constants.RegvorlSortColumns[strings.ToLower(sortBy)]
		if !ok {
			return errValidation("Fehler: ungültiger --sort-by Wert %q (gültig: kurztitel, stelle, datum)", sortBy)
		}
		params.Set("Sortierung.SortedByColumn", value)
	}

	return executeSearch(cmd, "Bundesrecht", "Suche in Regierungsvorlagen...", params)
}
