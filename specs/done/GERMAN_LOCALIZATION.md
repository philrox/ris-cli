# Deutsche Lokalisierung — Umsetzungsplan

## Ziel

Alle user-facing Texte (Help-Texte, Fehlermeldungen, Ausgabe-Labels) von Englisch auf Deutsch umstellen. Die CLI richtet sich an deutschsprachige Nutzer, die mit österreichischem Recht arbeiten.

## Scope

### Was wird übersetzt

| Kategorie | Dateien | Beispiel vorher → nachher |
|-----------|---------|---------------------------|
| Command Short/Long | alle `cmd/*.go` | `"Search federal laws"` → `"Bundesgesetze durchsuchen"` |
| Flag-Beschreibungen | alle `cmd/*.go`, `cmd/root.go` | `"Full-text search terms"` → `"Volltextsuche"` |
| Fehlermeldungen | alle `cmd/*.go`, `internal/api/client.go` | `"Error: at least one of --search or --title is required"` → `"Fehler: mindestens --search oder --title erforderlich"` |
| Ausgabe-Labels | `internal/format/text.go` | `"Results: 42 total"` → `"Ergebnisse: 42 gesamt"` |
| Technische Fehler | `internal/api/client.go` | `"HTTP request failed"` → `"HTTP-Anfrage fehlgeschlagen"` |

### Was NICHT übersetzt wird

| Kategorie | Grund | Beispiele |
|-----------|-------|-----------|
| Flag-Namen | CLI-Konvention, stabil für Scripting | `--search`, `--json`, `--page` |
| Command-Namen | Sind bereits deutsch/fachspezifisch | `bundesrecht`, `judikatur`, `bgbl` |
| JSON-Feldnamen | API-Surface für AI Agents und Scripts | `"dokumentnummer"`, `"total_hits"` |
| API-Parameter | Technische Identifier der RIS API | `"Suchworte"`, `"Applikation"` |
| Enum-Werte in Flags | Werden als Eingabe getippt, müssen stabil bleiben | `brkons`, `vfgh`, `spoe` |
| URLs, Regex, Hostnames | Technisch | `data.bka.gv.at`, `^[A-Z][A-Z0-9_]+$` |

## Umsetzung

### Schritt 1: `cmd/root.go` — Globale Texte

Root-Command und globale Flags.

```go
// VORHER:
Short: "Search and retrieve Austrian legal documents from the RIS API",
Long: `ris-cli — CLI for the Austrian Legal Information System (RIS)

Search and retrieve Austrian legal documents from the RIS OGD API.
Supports federal laws, state laws, court decisions, law gazettes, and more.

Primary output modes:
  Default    Formatted terminal output with colors
  --json     Machine-readable JSON (for AI agents and scripts)
  --plain    Plain text without colors (for piping)`,

// NACHHER:
Short: "Österreichische Rechtsdokumente aus dem RIS suchen und abrufen",
Long: `ris — CLI für das Rechtsinformationssystem des Bundes (RIS)

Suche und Abruf österreichischer Rechtsdokumente über die RIS OGD API.
Unterstützt Bundesrecht, Landesrecht, Judikatur, Gesetzblätter und mehr.

Ausgabemodi:
  Standard   Formatierte Terminalausgabe mit Farben
  --json     Maschinenlesbares JSON (für AI-Agents und Skripte)
  --plain    Klartext ohne Farben (für Piping)`,
```

Flag-Beschreibungen:

| Vorher | Nachher |
|--------|---------|
| `"Output as JSON (machine-readable)"` | `"Ausgabe als JSON (maschinenlesbar)"` |
| `"Output as plain text (stable, no colors)"` | `"Ausgabe als Klartext (stabil, ohne Farben)"` |
| `"Suppress non-essential output"` | `"Nicht-essentielle Ausgaben unterdrücken"` |
| `"Show HTTP request details on stderr"` | `"HTTP-Anfragen auf stderr anzeigen"` |
| `"Disable colored output (also respects NO_COLOR env)"` | `"Farbige Ausgabe deaktivieren (respektiert auch NO_COLOR)"` |
| `"HTTP request timeout"` | `"HTTP-Timeout"` |
| `"Page number for paginated results"` | `"Seitennummer für paginierte Ergebnisse"` |
| `"Results per page (10, 20, 50, 100)"` | `"Ergebnisse pro Seite (10, 20, 50, 100)"` |

Validierung:

| Vorher | Nachher |
|--------|---------|
| `"invalid limit %d: must be 10, 20, 50, or 100"` | `"ungültiges Limit %d: muss 10, 20, 50 oder 100 sein"` |

### Schritt 2: `internal/format/text.go` — Ausgabe-Labels

Diese Texte sieht jeder User bei jeder Suche.

| Vorher | Nachher |
|--------|---------|
| `"No results found."` | `"Keine Ergebnisse gefunden."` |
| `"Results: %d total (page %d, showing %d)"` | `"Ergebnisse: %d gesamt (Seite %d, zeige %d)"` |
| `"Nr: %s"` | `"Nr: %s"` (bleibt — ist deutsch) |
| `"Citation: %s"` | `"Zitat: %s"` |
| `"Case: %s"` | `"GZ: %s"` |
| `"Date: %s"` | `"Datum: %s"` |
| `"Summary: %s"` | `"Leitsatz: %s"` |
| `"More results available. Use --page %d to see next page."` | `"Weitere Ergebnisse verfügbar. Nächste Seite: --page %d"` |
| `"Document: %s"` | `"Dokument: %s"` |
| `"(untitled)"` | `"(ohne Titel)"` |

### Schritt 3: `internal/api/client.go` — Technische Fehler

| Vorher | Nachher |
|--------|---------|
| `"HTTP request failed: %w"` | `"HTTP-Anfrage fehlgeschlagen: %w"` |
| `"failed to read response body: %w"` | `"Antwort konnte nicht gelesen werden: %w"` |
| `"failed to read document body: %w"` | `"Dokument konnte nicht gelesen werden: %w"` |
| `"request timed out: %s"` | `"Zeitüberschreitung: %s"` |
| `"invalid URL: %w"` | `"Ungültige URL: %w"` |
| `"only HTTPS URLs are allowed, got %q"` | `"Nur HTTPS-URLs erlaubt, erhalten: %q"` |
| `"URL host %q is not allowed (allowed: ...)"` | `"Host %q nicht erlaubt (erlaubt: data.bka.gv.at, www.ris.bka.gv.at, ris.bka.gv.at)"` |

### Schritt 4: Command-Dateien — Short/Long/Errors

Pro Command: Short-Beschreibung, Long-Beschreibung (inkl. Beispiele), Flag-Beschreibungen, Fehlermeldungen.

#### `cmd/bundesrecht.go`

```go
Short: "Bundesgesetze durchsuchen (ABGB, StGB, etc.)",
Long: `Österreichische Bundesgesetze (Bundesrecht) durchsuchen.

Beispiele:
  ris bundesrecht --search "Mietrecht"
  ris bundesrecht --title "ABGB" --paragraph 1295
  ris bundesrecht --search "Schadenersatz" --app begut
  ris bundesrecht --search "Mietrecht" --date 2024-01-15 --json`,
```

| Flag vorher | Flag nachher |
|-------------|-------------|
| `"Full-text search terms"` | `"Volltextsuche"` |
| `"Search in law titles"` | `"Suche in Gesetzestitel"` |
| `"Paragraph number (e.g., \"1295\")"` | `"Paragraphennummer (z.B. \"1295\")"` |
| `"Application: brkons, begut, bgblauth, erv"` | `"Applikation: brkons, begut, bgblauth, erv"` |
| `"Historical version date (YYYY-MM-DD)"` | `"Fassungsdatum (JJJJ-MM-TT)"` |

| Error vorher | Error nachher |
|-------------|--------------|
| `"Error: at least one of --search, --title, or --paragraph is required"` | `"Fehler: mindestens --search, --title oder --paragraph erforderlich"` |
| `"Error: invalid --app value %q (valid: ...)"` | `"Fehler: ungültiger --app Wert %q (gültig: ...)"` |
| `"API request failed: %w"` | `"API-Anfrage fehlgeschlagen: %w"` |
| `"failed to parse response: %w"` | `"Antwort konnte nicht verarbeitet werden: %w"` |

#### `cmd/judikatur.go`

```go
Short: "Gerichtsentscheidungen durchsuchen",
Long: `Österreichische Gerichtsentscheidungen durchsuchen.

Beispiele:
  ris judikatur --search "Grundrecht" --court vfgh
  ris judikatur --case-number "5Ob234/20b"
  ris judikatur --norm "1319a ABGB" --from 2020-01-01 --to 2024-12-31`,
```

| Flag | Deutsch |
|------|---------|
| `"Full-text search terms"` | `"Volltextsuche"` |
| `"Legal norm reference"` | `"Normverweis"` |
| `"Case number (Geschaeftszahl)"` | `"Geschäftszahl"` |
| `"Court type: justiz, vfgh, ..."` | `"Gerichtstyp: justiz, vfgh, ..."` |
| `"Decision date from (YYYY-MM-DD)"` | `"Entscheidungsdatum von (JJJJ-MM-TT)"` |
| `"Decision date to (YYYY-MM-DD)"` | `"Entscheidungsdatum bis (JJJJ-MM-TT)"` |

#### `cmd/bgbl.go`

```go
Short: "Bundesgesetzblätter durchsuchen",
Long: `Bundesgesetzblätter (BGBl) durchsuchen.

Beispiele:
  ris bgbl --number 120 --year 2023 --part 1
  ris bgbl --search "Klimaschutz" --json`,
```

| Flag | Deutsch |
|------|---------|
| `"Gazette number"` | `"BGBl-Nummer"` |
| `"Year"` | `"Jahrgang"` |
| `"Title search"` | `"Titelsuche"` |
| `"Teil: 1 (Laws), 2 (Ordinances), 3 (Treaties)"` | `"Teil: 1 (Gesetze), 2 (Verordnungen), 3 (Staatsverträge)"` |
| `"Application: bgblauth, bgblpdf, bgblalt"` | `"Applikation: bgblauth, bgblpdf, bgblalt"` |

#### `cmd/landesrecht.go`

```go
Short: "Landesgesetze durchsuchen",
Long: `Österreichische Landesgesetze (Landesrecht) durchsuchen.

Beispiele:
  ris landesrecht --search "Bauordnung" --state salzburg
  ris landesrecht --title "Raumordnung" --state wien --json`,
```

| Flag | Deutsch |
|------|---------|
| `"Federal state filter (e.g., wien, salzburg, tirol)"` | `"Bundesland (z.B. wien, salzburg, tirol)"` |

#### `cmd/lgbl.go`

```go
Short: "Landesgesetzblätter durchsuchen",
Long: `Landesgesetzblätter (LGBl) durchsuchen.

Beispiele:
  ris lgbl --number 50 --year 2023 --state wien
  ris lgbl --search "Bauordnung" --state salzburg`,
```

#### `cmd/regvorl.go`

```go
Short: "Regierungsvorlagen durchsuchen",
Long: `Regierungsvorlagen durchsuchen.

Beispiele:
  ris regvorl --search "Klimaschutz"
  ris regvorl --ministry bmf --from 2024-01-01`,
```

| Flag | Deutsch |
|------|---------|
| `"Submitting ministry (e.g., bmf, bmi, bmj)"` | `"Einbringendes Ministerium (z.B. bmf, bmi, bmj)"` |
| `"Time filter: einerwoche, ..."` | `"Zeitfilter: einerwoche, ..."` |
| `"Sort direction: asc, desc"` | `"Sortierrichtung: asc, desc"` |
| `"Sort column: kurztitel, stelle, datum"` | `"Sortierspalte: kurztitel, stelle, datum"` |
| `"Decision date from (YYYY-MM-DD)"` | `"Beschlussdatum von (JJJJ-MM-TT)"` |
| `"Decision date to (YYYY-MM-DD)"` | `"Beschlussdatum bis (JJJJ-MM-TT)"` |

#### `cmd/dokument.go`

```go
Short: "Volltext eines Dokuments abrufen",
Long: `Volltext eines Rechtsdokuments abrufen.

Beispiele:
  ris dokument NOR40052761
  ris dokument NOR40052761 --json
  ris dokument --url "https://ris.bka.gv.at/Dokumente/Bundesnormen/NOR40052761/NOR40052761.html"`,
```

| Vorher | Nachher |
|--------|---------|
| `"Direct URL to document content"` | `"Direkte URL zum Dokumentinhalt"` |
| `"Error: either a document number argument or --url is required"` | `"Fehler: Dokumentnummer oder --url erforderlich"` |
| `"Error: document %q not found"` | `"Fehler: Dokument %q nicht gefunden"` |
| `"document number must be 5-50 characters, got %d"` | `"Dokumentnummer muss 5-50 Zeichen lang sein, erhalten: %d"` |
| `"invalid document number format %q ..."` | `"Ungültiges Dokumentnummer-Format %q ..."` |
| `"Direct URL failed (%v), trying search fallback..."` | `"Direkte URL fehlgeschlagen (%v), versuche Suche als Fallback..."` |

#### `cmd/bezirke.go`

```go
Short: "Bezirksverwaltungsbehörden-Kundmachungen durchsuchen",
Long: `Kundmachungen der Bezirksverwaltungsbehörden durchsuchen.

Beispiele:
  ris bezirke --state niederoesterreich --search "Bauordnung"
  ris bezirke --authority "Bezirkshauptmannschaft Innsbruck"`,
```

| Flag | Deutsch |
|------|---------|
| `"Federal state"` | `"Bundesland"` |
| `"District authority name"` | `"Bezirksverwaltungsbehörde"` |
| `"Announcement number"` | `"Kundmachungsnummer"` |
| `"Date from (YYYY-MM-DD)"` | `"Datum von (JJJJ-MM-TT)"` |
| `"Date to (YYYY-MM-DD)"` | `"Datum bis (JJJJ-MM-TT)"` |
| `"Time filter"` | `"Zeitfilter"` |

#### `cmd/gemeinden.go`

```go
Short: "Gemeinderecht durchsuchen",
Long: `Österreichisches Gemeinderecht durchsuchen.

Beispiele:
  ris gemeinden --municipality "Graz" --search "Parkgebuehren"
  ris gemeinden --state tirol --title "Gebuehrenordnung"`,
```

| Flag | Deutsch |
|------|---------|
| `"Municipality name"` | `"Gemeindename"` |
| `"File number (Gr only)"` | `"Geschäftszahl (nur Gr)"` |
| `"Subject area index (Gr only)"` | `"Sachbereichsindex (nur Gr)"` |
| `"District (GrA only)"` | `"Bezirk (nur GrA)"` |
| `"Municipal association (GrA only)"` | `"Gemeindeverband (nur GrA)"` |
| `"Announcement number (GrA only)"` | `"Kundmachungsnummer (nur GrA)"` |
| `"Historical version date (Gr only, YYYY-MM-DD)"` | `"Fassungsdatum (nur Gr, JJJJ-MM-TT)"` |
| `"Sort column (Gr only): ..."` | `"Sortierspalte (nur Gr): ..."` |

#### `cmd/history.go`

```go
Short: "Dokumentänderungshistorie durchsuchen",
Long: `Änderungshistorie von Dokumenten durchsuchen.

Beispiele:
  ris history --app bundesnormen --from 2024-01-01 --to 2024-01-31
  ris history --app justiz --from 2024-06-01 --include-deleted`,
```

| Flag | Deutsch |
|------|---------|
| `"Application to search (required)"` | `"Anwendung (erforderlich)"` |
| `"Changes from date (YYYY-MM-DD)"` | `"Änderungen von (JJJJ-MM-TT)"` |
| `"Changes to date (YYYY-MM-DD)"` | `"Änderungen bis (JJJJ-MM-TT)"` |
| `"Include deleted documents"` | `"Gelöschte Dokumente einschließen"` |

| Error | Deutsch |
|-------|---------|
| `"Error: --app is required"` | `"Fehler: --app ist erforderlich"` |
| `"Error: at least one of --from or --to is required"` | `"Fehler: mindestens --from oder --to erforderlich"` |

#### `cmd/verordnungen.go`

```go
Short: "Verordnungsblätter durchsuchen",
Long: `Verordnungsblätter der Länder durchsuchen.

Beispiele:
  ris verordnungen --search "Wolf" --state tirol
  ris verordnungen --number 25 --from 2024-01-01`,
```

| Flag | Deutsch |
|------|---------|
| `"Publication number"` | `"Kundmachungsnummer"` |

#### `cmd/sonstige.go`

Parent-Command:
```go
Short: "Sonstige Rechtssammlungen durchsuchen",
Long: `Sonstige Rechtssammlungen durchsuchen (8 Teil-Applikationen).

Unterbefehle:
  mrp          Ministerratsprotokolle
  erlaesse     Erlässe
  upts         Parteientransparenz-Entscheidungen
  kmger        Gerichtskundmachungen
  avsv         Sozialversicherungs-Kundmachungen
  avn          Veterinär-Kundmachungen
  spg          Gesundheitsstrukturpläne
  pruefgewo    Gewerbeprüfungen`,
```

Sub-Commands Short:

| Vorher | Nachher |
|--------|---------|
| `"Search Council of Ministers protocols"` | `"Ministerratsprotokolle durchsuchen"` |
| `"Search ministerial decrees"` | `"Erlässe durchsuchen"` |
| `"Search party transparency decisions"` | `"Parteientransparenz-Entscheidungen durchsuchen"` |
| `"Search court announcements"` | `"Gerichtskundmachungen durchsuchen"` |
| `"Search social insurance announcements"` | `"Sozialversicherungs-Kundmachungen durchsuchen"` |
| `"Search veterinary notices"` | `"Veterinär-Kundmachungen durchsuchen"` |
| `"Search health structure plans"` | `"Gesundheitsstrukturpläne durchsuchen"` |
| `"Search trade licensing examinations"` | `"Gewerbeprüfungen durchsuchen"` |

App-spezifische Flags:

| Vorher | Nachher |
|--------|---------|
| `"Submitter/ministry"` | `"Einbringer/Ministerium"` |
| `"Session number"` | `"Sitzungsnummer"` |
| `"Legislative period"` | `"Gesetzgebungsperiode"` |
| `"File number"` | `"Geschäftszahl"` |
| `"Federal ministry"` | `"Bundesministerium"` |
| `"Department"` | `"Abteilung"` |
| `"Source reference"` | `"Fundstelle"` |
| `"Legal norm"` | `"Norm"` |
| `"Version date (YYYY-MM-DD)"` | `"Fassungsdatum (JJJJ-MM-TT)"` |
| `"Political party: ..."` | `"Politische Partei: ..."` |
| `"Court name"` | `"Gericht"` |
| `"Document type"` | `"Dokumentart"` |
| `"Author/institution: ..."` | `"Urheber/Institution: ..."` |
| `"Federal state for RSG"` | `"Bundesland für RSG"` |

Gemeinsame Flags (Loop über alle 8 Sub-Commands):

| Vorher | Nachher |
|--------|---------|
| `"Full-text search terms"` | `"Volltextsuche"` |
| `"Title search"` | `"Titelsuche"` |
| `"Date from (YYYY-MM-DD)"` | `"Datum von (JJJJ-MM-TT)"` |
| `"Date to (YYYY-MM-DD)"` | `"Datum bis (JJJJ-MM-TT)"` |
| `"Time filter"` | `"Zeitfilter"` |
| `"Sort direction: asc, desc"` | `"Sortierrichtung: asc, desc"` |

Fehlermeldungen (alle Sonstige Sub-Commands):

| Vorher | Nachher |
|--------|---------|
| `"Error: invalid --ministry value %q"` | `"Fehler: ungültiger --ministry Wert %q"` |
| `"Error: invalid --party value %q"` | `"Fehler: ungültiger --party Wert %q"` |
| `"Error: invalid --type value %q ..."` | `"Fehler: ungültiger --type Wert %q ..."` |
| `"Error: invalid --author value %q"` | `"Fehler: ungültiger --author Wert %q"` |
| `"Error: invalid --osg-type value %q ..."` | `"Fehler: ungültiger --osg-type Wert %q ..."` |
| `"Error: invalid --rsg-type value %q ..."` | `"Fehler: ungültiger --rsg-type Wert %q ..."` |

#### `cmd/version.go`

```go
Short: "Versionsinformationen anzeigen",
Long:  "Version, Commit-Hash und Build-Datum der ris CLI anzeigen.",
```

## Betroffene Dateien (Übersicht)

| Datei | Strings | Aufwand |
|-------|---------|--------|
| `cmd/root.go` | ~11 | Klein |
| `cmd/bundesrecht.go` | ~9 | Klein |
| `cmd/judikatur.go` | ~11 | Klein |
| `cmd/bgbl.go` | ~11 | Klein |
| `cmd/landesrecht.go` | ~8 | Klein |
| `cmd/lgbl.go` | ~11 | Klein |
| `cmd/regvorl.go` | ~16 | Mittel |
| `cmd/dokument.go` | ~12 | Mittel |
| `cmd/bezirke.go` | ~11 | Klein |
| `cmd/gemeinden.go` | ~20 | Mittel |
| `cmd/sonstige.go` | ~45 | Groß |
| `cmd/history.go` | ~10 | Klein |
| `cmd/verordnungen.go` | ~8 | Klein |
| `cmd/version.go` | ~4 | Klein |
| `internal/format/text.go` | ~11 | Klein |
| `internal/api/client.go` | ~8 | Klein |
| **Gesamt** | **~148 Ersetzungen** | |

## Nicht-Betroffene Dateien

Diese Dateien brauchen KEINE Änderungen:

- `cmd/helpers.go` — keine user-facing Strings
- `internal/format/json.go` — nur JSON-Serialisierung
- `internal/format/citation.go` — formatiert bereits deutsche Daten
- `internal/format/html.go` — reine HTML→Text Konvertierung
- `internal/parser/*` — reine Datenverarbeitung
- `internal/model/*` — reine Typdefinitionen
- `internal/constants/*` — API-Werte, bleiben wie sie sind
- `internal/api/params.go` — nur Query-Parameter-Building
- `internal/api/endpoints.go` — nur Konstanten
- `main.go` — nur Entry-Point

## Reihenfolge

1. `cmd/root.go` + `internal/format/text.go` — höchste Sichtbarkeit
2. `internal/api/client.go` — Fehlermeldungen zentralisiert
3. `cmd/bundesrecht.go` — Template für alle anderen Commands
4. Restliche `cmd/*.go` — nach dem Muster von bundesrecht.go
5. `go build && go vet` — Validierung
6. Manueller Test: `ris --help`, `ris bundesrecht --help`, `ris bundesrecht --search "ABGB"`

## Konventionen

- **Fehler:** Immer mit `"Fehler: "` Prefix (statt `"Error: "`)
- **Datumsformat-Hinweis:** `JJJJ-MM-TT` statt `YYYY-MM-DD`
- **Validierungsfehler:** `"ungültiger --flag Wert %q (gültig: ...)"` als Pattern
- **Technische Fehler:** `"... fehlgeschlagen: %w"` / `"... konnte nicht ... werden: %w"` als Pattern
- **"at least one of":** `"mindestens ... erforderlich"` als Pattern
- **Beispiele-Header:** `"Beispiele:"` statt `"Examples:"`
