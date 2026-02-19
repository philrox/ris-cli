# ris

CLI für das Rechtsinformationssystem des Bundes (RIS) — Suche und Abruf österreichischer Rechtsdokumente über die [RIS OGD API](https://data.bka.gv.at/ris/api/v2.6/).

## Installation

### Homebrew (macOS/Linux)

```bash
brew install philrox/tap/ris
```

### Go

```bash
go install github.com/philrox/ris-cli@latest
```

### Binary Download

Fertige Binaries für Linux, macOS und Windows: [GitHub Releases](https://github.com/philrox/ris-cli/releases)

## Schnellstart

```bash
# Bundesrecht nach "Mietrecht" durchsuchen
ris bundesrecht --search "Mietrecht"

# Bestimmten ABGB-Paragraphen abrufen
ris bundesrecht --title "ABGB" --paragraph 1295

# JSON-Ausgabe für Skripte und AI-Agents
ris bundesrecht --search "Mietrecht" --json

# Volltext eines Dokuments abrufen
ris dokument NOR40052761

# VfGH-Entscheidungen zu Grundrechten
ris judikatur --search "Grundrecht" --court vfgh --from 2020-01-01

# Salzburger Landesrecht
ris landesrecht --search "Bauordnung" --state salzburg
```

## Befehle

| Befehl | Beschreibung |
|--------|-------------|
| `bundesrecht` | Bundesgesetze durchsuchen (ABGB, StGB, etc.) |
| `landesrecht` | Landesgesetze durchsuchen |
| `judikatur` | Gerichtsentscheidungen durchsuchen |
| `bgbl` | Bundesgesetzblätter durchsuchen |
| `lgbl` | Landesgesetzblätter durchsuchen |
| `regvorl` | Regierungsvorlagen durchsuchen |
| `dokument` | Volltext eines Dokuments abrufen |
| `bezirke` | Bezirksverwaltungsbehörden-Kundmachungen |
| `gemeinden` | Gemeinderecht durchsuchen |
| `sonstige` | Sonstige Rechtssammlungen (MRP, Erlässe, etc.) |
| `history` | Dokumentänderungshistorie |
| `verordnungen` | Verordnungsblätter durchsuchen |
| `completion` | Shell-Autovervollständigung generieren |
| `version` | Versionsinformationen anzeigen |

## Ausgabemodi

| Modus | Beschreibung |
|-------|-------------|
| Standard | Formatierte Terminalausgabe mit Farben |
| `--json` | Maschinenlesbares JSON (für AI-Agents und Skripte) |
| `--plain` | Klartext ohne Farben (für Piping) |

Die Ausgabe wird automatisch erkannt: Ist stdout ein Terminal, wird formatierter Text mit Farben ausgegeben. Bei Piping (`|`) wird automatisch Klartext verwendet.

## Beispiele

### Suche und Dokumentabruf

```bash
# Erstes Ergebnis als Dokumentnummer extrahieren
DOC=$(ris bundesrecht --search "Datenschutz" --json | jq -r '.documents[0].dokumentnummer')

# Volltext abrufen
ris dokument "$DOC" --json | jq '.content'
```

### Paginierung

```bash
ris judikatur --search "Schadenersatz" --page 2 --limit 50
```

### Bundesgesetzblatt

```bash
ris bgbl --number 120 --year 2023 --part 1
```

### Regierungsvorlagen

```bash
ris regvorl --ministry bmf --from 2024-01-01
```

### Ministerratsprotokolle

```bash
ris sonstige mrp --search "Budget" --session 42
```

### Dokumenthistorie

```bash
ris history --app bundesnormen --from 2024-01-01 --to 2024-01-31
```

## Globale Flags

| Flag | Kurz | Beschreibung |
|------|------|-------------|
| `--json` | `-j` | JSON-Ausgabe |
| `--plain` | | Klartext-Ausgabe |
| `--quiet` | `-q` | Nicht-essentielle Ausgaben unterdrücken |
| `--verbose` | `-v` | HTTP-Anfragen auf stderr anzeigen |
| `--no-color` | | Farben deaktivieren |
| `--no-pager` | | Pager deaktivieren |
| `--timeout` | | HTTP-Timeout (Standard: 30s) |
| `--page` | `-p` | Seitennummer (Standard: 1) |
| `--limit` | `-l` | Ergebnisse pro Seite (Standard: 20) |

## Umgebungsvariablen

| Variable | Beschreibung | Standard |
|----------|-------------|---------|
| `RIS_TIMEOUT` | HTTP-Timeout | `30s` |
| `RIS_BASE_URL` | API-Base-URL überschreiben | `https://data.bka.gv.at/ris/api/v2.6/` |
| `NO_COLOR` | Farben deaktivieren ([no-color.org](https://no-color.org/)) | — |
| `PAGER` | Pager für lange Ausgaben | `less -FIRX` |

Priorität: Flags > Umgebungsvariablen > Standardwerte

## Shell-Autovervollständigung

```bash
# Bash
source <(ris completion bash)

# Zsh
source <(ris completion zsh)

# Fish
ris completion fish | source
```

## Lizenz

MIT
