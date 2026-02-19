package cmd

import (
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
)

var (
	// Global flags
	jsonOutput  bool
	plainOutput bool
	quiet       bool
	verbose     bool
	noColor     bool
	noPager     bool
	timeout     time.Duration
	page        int
	limit       int

	// isTTY is true when stdout is connected to a terminal.
	isTTY bool
)

var rootCmd = &cobra.Command{
	Use:   "risgo",
	Short: "Österreichische Rechtsdokumente aus dem RIS suchen und abrufen",
	Long: `risgo — CLI für das Rechtsinformationssystem des Bundes (RIS)

Suche und Abruf österreichischer Rechtsdokumente über die RIS OGD API.
Unterstützt Bundesrecht, Landesrecht, Judikatur, Gesetzblätter und mehr.

Ausgabemodi:
  Standard   Formatierte Terminalausgabe mit Farben
  --json     Maschinenlesbares JSON (für AI-Agents und Skripte)
  --plain    Klartext ohne Farben (für Piping)`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func Execute() error {
	return rootCmd.Execute()
}

// RootCmd returns the root cobra command for doc generation.
func RootCmd() *cobra.Command {
	return rootCmd
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().BoolVarP(&jsonOutput, "json", "j", false, "Ausgabe als JSON (maschinenlesbar)")
	rootCmd.PersistentFlags().BoolVar(&plainOutput, "plain", false, "Ausgabe als Klartext (stabil, ohne Farben)")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "Nicht-essentielle Ausgaben unterdrücken")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "HTTP-Anfragen auf stderr anzeigen")
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "Farbige Ausgabe deaktivieren (respektiert auch NO_COLOR)")
	rootCmd.PersistentFlags().BoolVar(&noPager, "no-pager", false, "Pager für lange Ausgaben deaktivieren")
	rootCmd.PersistentFlags().DurationVar(&timeout, "timeout", 30*time.Second, "HTTP-Timeout")
	rootCmd.PersistentFlags().IntVarP(&page, "page", "p", 1, "Seitennummer für paginierte Ergebnisse")
	rootCmd.PersistentFlags().IntVarP(&limit, "limit", "l", 20, "Ergebnisse pro Seite (10, 20, 50, 100)")
}

func initConfig() {
	// Detect whether stdout is a terminal.
	isTTY = isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd())

	// When not on a TTY (e.g. piped), behave as if --plain was set.
	if !isTTY || plainOutput {
		color.NoColor = true
	}

	// Respect NO_COLOR environment variable (https://no-color.org/)
	if _, ok := os.LookupEnv("NO_COLOR"); ok {
		noColor = true
	}

	if noColor {
		color.NoColor = true
	}

	// Respect RIS_TIMEOUT environment variable
	if envTimeout := os.Getenv("RIS_TIMEOUT"); envTimeout != "" {
		if d, err := time.ParseDuration(envTimeout); err == nil {
			// Only use env var if flag was not explicitly set
			if !rootCmd.PersistentFlags().Changed("timeout") {
				timeout = d
			}
		}
	}
}
