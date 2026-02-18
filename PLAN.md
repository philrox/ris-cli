# ris-cli — CLI for the Austrian Legal Information System (RIS)

## 1. Overview

| Field | Value |
|-------|-------|
| **Name** | `ris` |
| **One-liner** | Search and retrieve Austrian legal documents from the RIS API |
| **Language** | Go |
| **Primary users** | AI agents (JSON output), humans (formatted terminal output) |
| **Repository** | `ris-cli` (standalone, separate from `ris-mcp-ts`) |
| **API** | RIS OGD API v2.6 — `https://data.bka.gv.at/ris/api/v2.6/` |

## 2. CLI Spec

### USAGE

```
ris [global flags] <command> [command flags]
```

### Command Tree

```
ris
├── bundesrecht      Search federal laws (ABGB, StGB, etc.)
├── landesrecht      Search state/provincial laws
├── judikatur        Search court decisions
├── bgbl             Search Federal Law Gazettes (Bundesgesetzblatt)
├── lgbl             Search State Law Gazettes (Landesgesetzblatt)
├── regvorl          Search Government Bills (Regierungsvorlagen)
├── dokument         Retrieve full document text by number or URL
├── bezirke          Search district authority announcements
├── gemeinden        Search municipal law
├── sonstige         Search miscellaneous legal collections
├── history          Search document change history
├── verordnungen     Search state ordinance gazettes
└── version          Print version information
```

### Global Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--help` | `-h` | bool | | Show help |
| `--json` | `-j` | bool | false | Output as JSON (machine-readable) |
| `--plain` | | bool | false | Output as plain text (stable, no colors) |
| `--quiet` | `-q` | bool | false | Suppress non-essential output |
| `--verbose` | `-v` | bool | false | Show HTTP request details on stderr |
| `--no-color` | | bool | auto | Disable colored output (also respects `NO_COLOR` env) |
| `--timeout` | | duration | 30s | HTTP request timeout |
| `--page` | `-p` | int | 1 | Page number for paginated results |
| `--limit` | `-l` | int | 20 | Results per page (10, 20, 50, 100) |

### I/O Contract

| Stream | Content |
|--------|---------|
| **stdout** | Primary data: search results, document text. JSON when `--json`, formatted text otherwise. |
| **stderr** | Progress spinners, HTTP debug info (`--verbose`), error messages, warnings. |

### Exit Codes

| Code | Meaning |
|------|---------|
| `0` | Success |
| `1` | API error (timeout, HTTP error, parsing failure) |
| `2` | Invalid usage (bad flags, missing required params) |
| `3` | Document not found |

## 3. Subcommand Specifications

### 3.1 `ris bundesrecht`

Search Austrian federal laws (Bundesrecht).

```
ris bundesrecht --search "Mietrecht"
ris bundesrecht --title "ABGB" --paragraph 1295
ris bundesrecht --search "Schadenersatz" --app begut
ris bundesrecht --search "Mietrecht" --date 2024-01-15 --json
```

| Flag | Short | Type | Required | API Param | Description |
|------|-------|------|----------|-----------|-------------|
| `--search` | `-s` | string | * | `Suchworte` | Full-text search terms |
| `--title` | `-t` | string | * | `Titel` | Search in law titles |
| `--paragraph` | | string | * | `Abschnitt.Von/Bis/Typ` | Paragraph number (e.g., "1295") |
| `--app` | | enum | | `Applikation` | `brkons` (default), `begut`, `bgblauth`, `erv` |
| `--date` | | date | | `FassungVom` | Historical version date (YYYY-MM-DD) |

*At least one of the `*` params required.*

**API Endpoint:** `Bundesrecht`
**Default Applikation:** `BrKons`

**Applikation values:**
| CLI value | API value | Description |
|-----------|-----------|-------------|
| `brkons` | `BrKons` | Consolidated federal law (default) |
| `begut` | `Begut` | Draft legislation |
| `bgblauth` | `BgblAuth` | Federal Law Gazette authentic |
| `erv` | `Erv` | English translations |

### 3.2 `ris landesrecht`

Search Austrian state/provincial laws (Landesrecht).

```
ris landesrecht --search "Bauordnung" --state salzburg
ris landesrecht --title "Raumordnung" --state wien --json
```

| Flag | Short | Type | Required | API Param | Description |
|------|-------|------|----------|-----------|-------------|
| `--search` | `-s` | string | * | `Suchworte` | Full-text search |
| `--title` | `-t` | string | * | `Titel` | Search in titles |
| `--state` | | enum | * | `Bundesland.SucheIn*` | Federal state filter |

*At least one required.*

**API Endpoint:** `Landesrecht`
**Default Applikation:** `LrKons`

**State values and API mapping:**
| CLI value | API param |
|-----------|-----------|
| `wien` | `Bundesland.SucheInWien=true` |
| `niederoesterreich` | `Bundesland.SucheInNiederoesterreich=true` |
| `oberoesterreich` | `Bundesland.SucheInOberoesterreich=true` |
| `salzburg` | `Bundesland.SucheInSalzburg=true` |
| `tirol` | `Bundesland.SucheInTirol=true` |
| `vorarlberg` | `Bundesland.SucheInVorarlberg=true` |
| `kaernten` | `Bundesland.SucheInKaernten=true` |
| `steiermark` | `Bundesland.SucheInSteiermark=true` |
| `burgenland` | `Bundesland.SucheInBurgenland=true` |

### 3.3 `ris judikatur`

Search Austrian court decisions.

```
ris judikatur --search "Grundrecht" --court vfgh
ris judikatur --case-number "5Ob234/20b"
ris judikatur --norm "1319a ABGB" --from 2020-01-01 --to 2024-12-31
```

| Flag | Short | Type | Required | API Param | Description |
|------|-------|------|----------|-----------|-------------|
| `--search` | `-s` | string | * | `Suchworte` | Full-text search |
| `--norm` | `-n` | string | * | `Norm` | Legal norm reference |
| `--case-number` | | string | * | `Geschaeftszahl` | Case number |
| `--court` | `-c` | enum | | `Applikation` | Court type (default: `justiz`) |
| `--from` | | date | | `EntscheidungsdatumVon` | Decision date from |
| `--to` | | date | | `EntscheidungsdatumBis` | Decision date to |

*At least one of `*` required.*

**API Endpoint:** `Judikatur`

**Court values:**
| CLI value | API value | Description |
|-----------|-----------|-------------|
| `justiz` | `Justiz` | Ordinary courts — OGH, OLG, LG, BG (default) |
| `vfgh` | `Vfgh` | Constitutional Court |
| `vwgh` | `Vwgh` | Supreme Administrative Court |
| `bvwg` | `Bvwg` | Federal Administrative Court |
| `lvwg` | `Lvwg` | State Administrative Courts |
| `dsk` | `Dsk` | Data Protection Authority |
| `asylgh` | `AsylGH` | Asylum Court (historical, until 2013) |
| `normenliste` | `Normenliste` | Court norm lists |
| `pvak` | `Pvak` | Personnel Representation Commission |
| `gbk` | `Gbk` | Equal Treatment Commission |
| `dok` | `Dok` | Disciplinary Commission |

### 3.4 `ris bgbl`

Search Federal Law Gazettes (Bundesgesetzblatt).

```
ris bgbl --number 120 --year 2023 --part 1
ris bgbl --search "Klimaschutz" --json
```

| Flag | Short | Type | Required | API Param | Description |
|------|-------|------|----------|-----------|-------------|
| `--number` | | string | * | `Bgblnummer` | Gazette number |
| `--year` | | string | * | `Jahrgang` | Year |
| `--search` | `-s` | string | * | `Suchworte` | Full-text search |
| `--title` | `-t` | string | * | `Titel` | Title search |
| `--part` | | enum | | `Teil` | `1` (Laws), `2` (Ordinances), `3` (Treaties) |
| `--app` | | enum | | `Applikation` | `bgblauth` (default), `bgblpdf`, `bgblalt` |

**API Endpoint:** `Bundesrecht`

**Applikation values:**
| CLI value | API value | Description |
|-----------|-----------|-------------|
| `bgblauth` | `BgblAuth` | Authentic, 2004+ (default) |
| `bgblpdf` | `BgblPdf` | PDF format |
| `bgblalt` | `BgblAlt` | Historical, 1945-2003 |

### 3.5 `ris lgbl`

Search State Law Gazettes (Landesgesetzblatt).

```
ris lgbl --number 50 --year 2023 --state wien
ris lgbl --search "Bauordnung" --state salzburg
```

| Flag | Short | Type | Required | API Param | Description |
|------|-------|------|----------|-----------|-------------|
| `--number` | | string | * | `Lgblnummer` | Gazette number |
| `--year` | | string | * | `Jahrgang` | Year |
| `--state` | | enum | * | `Bundesland.SucheIn*` | Federal state |
| `--search` | `-s` | string | * | `Suchworte` | Full-text search |
| `--title` | `-t` | string | * | `Titel` | Title search |
| `--app` | | enum | | `Applikation` | `lgblauth` (default), `lgbl`, `lgblno` |

**API Endpoint:** `Landesrecht`

**Applikation values:**
| CLI value | API value | Description |
|-----------|-----------|-------------|
| `lgblauth` | `LgblAuth` | Authentic (default) |
| `lgbl` | `Lgbl` | General |
| `lgblno` | `LgblNO` | Lower Austria |

### 3.6 `ris regvorl`

Search Government Bills (Regierungsvorlagen).

```
ris regvorl --search "Klimaschutz"
ris regvorl --ministry bmf --from 2024-01-01
```

| Flag | Short | Type | Required | API Param | Description |
|------|-------|------|----------|-----------|-------------|
| `--search` | `-s` | string | * | `Suchworte` | Full-text search |
| `--title` | `-t` | string | * | `Titel` | Title search |
| `--from` | | date | * | `BeschlussdatumVon` | Decision date from |
| `--to` | | date | | `BeschlussdatumBis` | Decision date to |
| `--ministry` | | enum | * | `EinbringendeStelle` | Submitting ministry |
| `--since` | | enum | * | `ImRisSeit` | Time filter |
| `--sort-dir` | | enum | | `Sortierung.SortDirection` | `asc` / `desc` |
| `--sort-by` | | enum | | `Sortierung.SortedByColumn` | `kurztitel`, `stelle`, `datum` |

**API Endpoint:** `Bundesrecht`
**Applikation:** `RegV` (fixed)

**Ministry values (--ministry):**
| CLI value | API value |
|-----------|-----------|
| `bka` | `BKA (Bundeskanzleramt)` |
| `bmkoes` | `BMKOES (Bundesministerium für Kunst, Kultur, öffentlichen Dienst und Sport)` |
| `bmeia` | `BMEIA (Bundesministerium für europäische und internationale Angelegenheiten)` |
| `bmaw` | `BMAW (Bundesministerium für Arbeit und Wirtschaft)` |
| `bmbwf` | `BMBWF (Bundesministerium für Bildung, Wissenschaft und Forschung)` |
| `bmf` | `BMF (Bundesministerium für Finanzen)` |
| `bmi` | `BMI (Bundesministerium für Inneres)` |
| `bmj` | `BMJ (Bundesministerium für Justiz)` |
| `bmk` | `BMK (Bundesministerium für Klimaschutz, Umwelt, Energie, Mobilität, Innovation und Technologie)` |
| `bmlv` | `BMLV (Bundesministerium für Landesverteidigung)` |
| `bml` | `BML (Bundesministerium für Land- und Forstwirtschaft, Regionen und Wasserwirtschaft)` |
| `bmsgpk` | `BMSGPK (Bundesministerium für Soziales, Gesundheit, Pflege und Konsumentenschutz)` |
| `bmffim` | `BMFFIM (Bundesministerin für Frauen, Familie, Integration und Medien im Bundeskanzleramt)` |
| `bmeuv` | `BMEUV (Bundesministerin für EU und Verfassung im Bundeskanzleramt)` |

### 3.7 `ris dokument`

Retrieve full text of a legal document.

```
ris dokument NOR40052761
ris dokument NOR40052761 --json
ris dokument --url "https://ris.bka.gv.at/Dokumente/Bundesnormen/NOR40052761/NOR40052761.html"
```

| Flag | Short | Type | Required | Description |
|------|-------|------|----------|-------------|
| (positional) | | string | * | Document number (e.g., `NOR40052761`) |
| `--url` | | string | * | Direct URL to document content |

*Either positional document number or `--url` required.*

**Strategy (same as MCP server):**
1. Try direct URL construction from document number prefix
2. Fallback to search API to find document URL
3. Fetch HTML content from URL
4. Convert HTML to text, format, output

**SSRF protection:** Only allow URLs to `data.bka.gv.at`, `www.ris.bka.gv.at`, `ris.bka.gv.at` (HTTPS only).

**Document number validation:**
- 5-50 characters
- Start with uppercase letter
- Only uppercase letters, digits, underscores: `^[A-Z][A-Z0-9_]+$`

### 3.8 `ris bezirke`

Search district administrative authority announcements.

```
ris bezirke --state niederoesterreich --search "Bauordnung"
ris bezirke --authority "Bezirkshauptmannschaft Innsbruck"
```

| Flag | Short | Type | Required | API Param | Description |
|------|-------|------|----------|-----------|-------------|
| `--search` | `-s` | string | * | `Suchworte` | Full-text search |
| `--title` | `-t` | string | * | `Titel` | Title search |
| `--state` | | enum | * | `Bundesland` | Federal state |
| `--authority` | | string | * | `Bezirksverwaltungsbehoerde` | District authority name |
| `--number` | | string | * | `Kundmachungsnummer` | Announcement number |
| `--from` | | date | | `Kundmachungsdatum.Von` | Date from |
| `--to` | | date | | `Kundmachungsdatum.Bis` | Date to |
| `--since` | | enum | | `ImRisSeit` | Time filter |

**API Endpoint:** `Bezirke`
**Applikation:** `Bvb` (fixed)

**Note:** State values here use display names with Umlauts (`Niederösterreich`, `Oberösterreich`, `Kärnten`) — different from Landesrecht which uses ASCII versions. CLI normalizes lowercase input to correct API format.

### 3.9 `ris gemeinden`

Search Austrian municipal law.

```
ris gemeinden --municipality "Graz" --search "Parkgebuehren"
ris gemeinden --state tirol --title "Gebuehrenordnung"
```

| Flag | Short | Type | Required | API Param | Description |
|------|-------|------|----------|-----------|-------------|
| `--search` | `-s` | string | * | `Suchworte` | Full-text search |
| `--title` | `-t` | string | * | `Titel` | Title search |
| `--state` | | string | * | `Bundesland` | Federal state |
| `--municipality` | | string | * | `Gemeinde` | Municipality name |
| `--file-number` | | string | * | `Geschaeftszahl` | File number (Gr only) |
| `--index` | | enum | * | `Index` | Subject area (Gr only) |
| `--district` | | string | * | `Bezirk` | District (GrA only) |
| `--announcement-nr` | | string | * | `Kundmachungsnummer` | Announcement nr (GrA only) |
| `--app` | | enum | | `Applikation` | `gr` (default), `gra` |
| `--date` | | date | | `FassungVom` | Historical version (Gr only) |
| `--from` | | date | | `Kundmachungsdatum.Von` | Date from (GrA only) |
| `--to` | | date | | `Kundmachungsdatum.Bis` | Date to (GrA only) |
| `--since` | | enum | | `ImRisSeit` | Time filter |
| `--sort-dir` | | enum | | `Sortierung.SortDirection` | `asc` / `desc` |

**Index values (Gr):**
`Undefined`, `VertretungskoerperUndAllgemeineVerwaltung`, `OeffentlicheOrdnungUndSicherheit`, `UnterrichtErziehungSportUndWissenschaft`, `KunstKulturUndKultus`, `SozialeWohlfahrtUndWohnbaufoerderung`, `Gesundheit`, `StraßenUndWasserbauVerkehr`, `Wirtschaftsfoerderung`, `Dienstleistungen`, `Finanzwirtschaft`

### 3.10 `ris sonstige`

Search miscellaneous legal collections (8 sub-applications).

```
ris sonstige mrp --search "Budget"
ris sonstige erlaesse --ministry bmf
ris sonstige upts --party spoe
ris sonstige kmger --type geschaeftsordnung
```

**Design:** `ris sonstige` has a required sub-subcommand for the application type.

```
ris sonstige
├── mrp          Council of Ministers protocols
├── erlaesse     Ministerial decrees
├── upts         Party transparency decisions
├── kmger        Court announcements
├── avsv         Social insurance announcements
├── avn          Veterinary notices
├── spg          Health structure plans
└── pruefgewo    Trade licensing examinations
```

**Common flags (all sub-apps):**
| Flag | Short | Type | API Param |
|------|-------|------|-----------|
| `--search` | `-s` | string | `Suchworte` |
| `--title` | `-t` | string | `Titel` |
| `--from` | | date | (app-specific date param) |
| `--to` | | date | (app-specific date param) |
| `--since` | | enum | `ImRisSeit` |
| `--sort-dir` | | enum | `Sortierung.SortDirection` |

**App-specific flags:**

| Sub-app | Flag | API Param | Description |
|---------|------|-----------|-------------|
| `mrp` | `--submitter` | `Einbringer` | Submitter/ministry |
| `mrp` | `--session` | `Sitzungsnummer` | Session number |
| `mrp` | `--period` | `Gesetzgebungsperiode` | Legislative period |
| `mrp` | `--file-number` | `Geschaeftszahl` | File number |
| `erlaesse` | `--ministry` | `Bundesministerium` | Federal ministry |
| `erlaesse` | `--department` | `Abteilung` | Department |
| `erlaesse` | `--source` | `Fundstelle` | Source reference |
| `erlaesse` | `--norm` | `Norm` | Legal norm |
| `erlaesse` | `--date` | `FassungVom` | Version date |
| `upts` | `--party` | `Partei` | Political party |
| `upts` | `--file-number` | `Geschaeftszahl` | File number |
| `upts` | `--norm` | `Norm` | Legal norm |
| `kmger` | `--type` | `Typ` | `geschaeftsordnung` / `geschaeftsverteilung` |
| `kmger` | `--court-name` | `Gericht` | Court name |
| `kmger` | `--file-number` | `Geschaeftszahl` | File number |
| `avsv` | `--doc-type` | `Dokumentart` | Document type |
| `avsv` | `--author` | `Urheber` | Author/institution |
| `avsv` | `--avsv-number` | `Avsvnummer` | AVSV number |
| `avn` | `--avn-number` | `Avnnummer` | AVN number |
| `avn` | `--type` | `Typ` | `kundmachung` / `verordnung` / `erlass` |
| `spg` | `--spg-number` | `Spgnummer` | SPG number |
| `spg` | `--osg-type` | `OsgTyp` | `oesg` / `oesg-grossgeraete` |
| `spg` | `--rsg-type` | `RsgTyp` | `rsg` / `rsg-grossgeraete` |
| `spg` | `--rsg-state` | `RsgLand` | Federal state for RSG |
| `pruefgewo` | `--type` | `Typ` | `befaehigung` / `eignung` / `meister` |

**Date parameter mapping per app:**
| App | `--from` API param | `--to` API param |
|-----|--------------------|------------------|
| `mrp` | `Sitzungsdatum.Von` | `Sitzungsdatum.Bis` |
| `upts` | `Entscheidungsdatum.Von` | `Entscheidungsdatum.Bis` |
| `erlaesse` | `VonInkrafttretensdatum` | `BisInkrafttretensdatum` |
| `pruefgewo`, `spg`, `kmger` | `Kundmachungsdatum.Von` | `Kundmachungsdatum.Bis` |
| `avsv`, `avn` | `Kundmachung.Von` | `Kundmachung.Bis` |

**UPTS Party values:**
| CLI value | API value |
|-----------|-----------|
| `spoe` | `SPÖ - Sozialdemokratische Partei Österreichs` |
| `oevp` | `ÖVP - Österreichische Volkspartei` |
| `fpoe` | `FPÖ - Freiheitliche Partei Österreichs` |
| `gruene` | `GRÜNE - Die Grünen - Die Grüne Alternative` |
| `neos` | `NEOS - NEOS – Das Neue Österreich und Liberales Forum` |
| `bzoe` | `BZÖ - Bündnis Zukunft Österreich` |

**Erlaesse Ministry values:** Same enum as `ris regvorl --ministry`, mapped to full names from `BUNDESMINISTERIEN` constant.

**AVSV Author values:**
| CLI value | API value |
|-----------|-----------|
| `dvsv` | `Dachverband der Sozialversicherungsträger (DVSV)` |
| `pva` | `Pensionsversicherungsanstalt (PVA)` |
| `oegk` | `Österreichische Gesundheitskasse (ÖGK)` |
| `auva` | `Allgemeine Unfallversicherungsanstalt (AUVA)` |
| `svs` | `Sozialversicherungsanstalt der Selbständigen (SVS)` |
| `bvaeb` | `Versicherungsanstalt öffentlich Bediensteter, Eisenbahnen und Bergbau (BVAEB)` |

### 3.11 `ris history`

Search document change history.

```
ris history --app bundesnormen --from 2024-01-01 --to 2024-01-31
ris history --app justiz --from 2024-06-01 --include-deleted
```

| Flag | Short | Type | Required | API Param | Description |
|------|-------|------|----------|-----------|-------------|
| `--app` | `-a` | enum | yes | `Anwendung` | Application to search (required) |
| `--from` | | date | * | `AenderungenVon` | Changes from date |
| `--to` | | date | * | `AenderungenBis` | Changes to date |
| `--include-deleted` | | bool | | `IncludeDeletedDocuments` | Include deleted docs |

*At least one of `--from` or `--to` required.*

**API Endpoint:** `History`
**Note:** Uses `Anwendung` param, NOT `Applikation`.

**Application values (30 total):**
`bundesnormen`, `landesnormen`, `justiz`, `vfgh`, `vwgh`, `bvwg`, `lvwg`, `bgblauth`, `bgblalt`, `bgblpdf`, `lgblauth`, `lgbl`, `lgblno`, `gemeinderecht`, `gemeinderechtauth`, `bvb`, `vbl`, `regv`, `mrp`, `erlaesse`, `pruefgewo`, `avsv`, `spg`, `kmger`, `dsk`, `gbk`, `dok`, `pvak`, `normenliste`, `asylgh`

### 3.12 `ris verordnungen`

Search state ordinance gazettes (Verordnungsblätter).

```
ris verordnungen --search "Wolf" --state tirol
ris verordnungen --number 25 --from 2024-01-01
```

| Flag | Short | Type | Required | API Param | Description |
|------|-------|------|----------|-----------|-------------|
| `--search` | `-s` | string | * | `Suchworte` | Full-text search |
| `--title` | `-t` | string | * | `Titel` | Title search |
| `--state` | | enum | * | `Bundesland` | Federal state |
| `--number` | | string | * | `Kundmachungsnummer` | Publication number |
| `--from` | | date | * | `Kundmachungsdatum.Von` | Date from |
| `--to` | | date | | `Kundmachungsdatum.Bis` | Date to |

**API Endpoint:** `Landesrecht`
**Applikation:** `Vbl` (fixed)
**Note:** Currently only Tirol data available (since 2022-01-01). Uses direct `Bundesland` values, NOT `SucheIn*` format.

## 4. API Reference

### Base URL

```
https://data.bka.gv.at/ris/api/v2.6/
```

### Endpoints

| Endpoint | Used by CLI commands |
|----------|---------------------|
| `Bundesrecht` | `bundesrecht`, `bgbl`, `regvorl` |
| `Landesrecht` | `landesrecht`, `lgbl`, `verordnungen` |
| `Judikatur` | `judikatur` |
| `Bezirke` | `bezirke` |
| `Gemeinden` | `gemeinden` |
| `Sonstige` | `sonstige` |
| `History` | `history` |

### Common API Parameters

| API Param | Description | Notes |
|-----------|-------------|-------|
| `Applikation` | Sub-application within endpoint | Varies per endpoint |
| `DokumenteProSeite` | Results per page | `Ten`, `Twenty`, `Fifty`, `OneHundred` |
| `Seitennummer` | Page number | 1-based |
| `Suchworte` | Full-text search | Max 1000 chars |
| `Titel` | Title search | Max 500 chars |
| `ImRisSeit` | Time filter | `EinerWoche`, `ZweiWochen`, `EinemMonat`, `DreiMonaten`, `SechsMonaten`, `EinemJahr` |
| `Sortierung.SortDirection` | Sort direction | `Ascending`, `Descending` |

### Response Structure

```json
{
  "OgdSearchResult": {
    "OgdDocumentResults": {
      "Hits": {
        "#text": "42",
        "@pageNumber": "1",
        "@pageSize": "20"
      },
      "OgdDocumentReference": [
        {
          "Data": {
            "Metadaten": {
              "Technisch": {
                "ID": "NOR40052761",
                "Applikation": "BrKons"
              },
              "Allgemein": {
                "DokumentUrl": "https://..."
              },
              "Bundesrecht": {
                "Kurztitel": "ABGB",
                "Langtitel": "Allgemeines buergerliches Gesetzbuch",
                "Titel": "...",
                "Eli": "...",
                "BrKons": {
                  "Kundmachungsorgan": "...",
                  "ArtikelParagraphAnlage": "...",
                  "Inkrafttretensdatum": "...",
                  "Ausserkrafttretensdatum": "...",
                  "GesamteRechtsvorschriftUrl": "..."
                }
              }
            },
            "Dokumentliste": {
              "ContentReference": {
                "ContentType": "MainDocument",
                "Urls": {
                  "ContentUrl": [
                    { "DataType": "Html", "Url": "https://..." },
                    { "DataType": "Xml", "Url": "https://..." }
                  ]
                }
              }
            }
          }
        }
      ]
    }
  }
}
```

**Parsing notes:**
- `Hits` can be an object with `#text`, `@pageNumber`, `@pageSize` or a plain number
- `OgdDocumentReference` can be a single object or an array
- `ContentReference` can be a single object or an array (use `MainDocument` type)
- `ContentUrl` can be a single object or an array
- Metadata varies: `Bundesrecht`, `Landesrecht`, or `Judikatur` section present
- `Landesrecht` nests under `LrKons`, `Bundesrecht` under `BrKons`
- `Judikatur` has `Geschaeftszahl` (can be string or `{item: string|string[]}`)
- `Name` field in `ContentReference` can be string or `{"#text": string}`

## 5. Document Routing Table

### Prefix → Direct URL Pattern

For `ris dokument`, construct direct URLs from document number prefix:

| Prefix | URL Pattern (replace `{nr}` with document number) |
|--------|---------------------------------------------------|
| `NOR` | `https://ris.bka.gv.at/Dokumente/Bundesnormen/{nr}/{nr}.html` |
| `LBG` | `https://ris.bka.gv.at/Dokumente/LrBgld/{nr}/{nr}.html` |
| `LKT` | `https://ris.bka.gv.at/Dokumente/LrK/{nr}/{nr}.html` |
| `LNO` | `https://ris.bka.gv.at/Dokumente/LrNO/{nr}/{nr}.html` |
| `LOO` | `https://ris.bka.gv.at/Dokumente/LrOO/{nr}/{nr}.html` |
| `LSB` | `https://ris.bka.gv.at/Dokumente/LrSbg/{nr}/{nr}.html` |
| `LST` | `https://ris.bka.gv.at/Dokumente/LrStmk/{nr}/{nr}.html` |
| `LTI` | `https://ris.bka.gv.at/Dokumente/LrT/{nr}/{nr}.html` |
| `LVB` | `https://ris.bka.gv.at/Dokumente/LrVbg/{nr}/{nr}.html` |
| `LWI` | `https://ris.bka.gv.at/Dokumente/LrW/{nr}/{nr}.html` |
| `JWR` | `https://ris.bka.gv.at/Dokumente/Vwgh/{nr}/{nr}.html` |
| `JFR` | `https://ris.bka.gv.at/Dokumente/Vfgh/{nr}/{nr}.html` |
| `JFT` | `https://ris.bka.gv.at/Dokumente/Vfgh/{nr}/{nr}.html` |
| `JWT` | `https://ris.bka.gv.at/Dokumente/Justiz/{nr}/{nr}.html` |
| `JJR` | `https://ris.bka.gv.at/Dokumente/Justiz/{nr}/{nr}.html` |
| `BVWG` | `https://ris.bka.gv.at/Dokumente/Bvwg/{nr}/{nr}.html` |
| `LVWG` | `https://ris.bka.gv.at/Dokumente/Lvwg/{nr}/{nr}.html` |
| `DSB` | `https://ris.bka.gv.at/Dokumente/Dsk/{nr}/{nr}.html` |
| `GBK` | `https://ris.bka.gv.at/Dokumente/Gbk/{nr}/{nr}.html` |
| `PVAK` | `https://ris.bka.gv.at/Dokumente/Pvak/{nr}/{nr}.html` |
| `ASYLGH` | `https://ris.bka.gv.at/Dokumente/AsylGH/{nr}/{nr}.html` |
| `BGBLA` | `https://ris.bka.gv.at/Dokumente/BgblAuth/{nr}/{nr}.html` |
| `BGBL` | `https://ris.bka.gv.at/Dokumente/BgblAlt/{nr}/{nr}.html` |
| `BGBLPDF` | `https://ris.bka.gv.at/Dokumente/BgblPdf/{nr}/{nr}.html` |
| `REGV` | `https://ris.bka.gv.at/Dokumente/RegV/{nr}/{nr}.html` |
| `BVB` | `https://ris.bka.gv.at/Dokumente/Bvb/{nr}/{nr}.html` |
| `VBL` | `https://ris.bka.gv.at/Dokumente/Vbl/{nr}/{nr}.html` |
| `MRP` | `https://ris.bka.gv.at/Dokumente/Mrp/{nr}/{nr}.html` |
| `ERL` | `https://ris.bka.gv.at/Dokumente/Erlaesse/{nr}/{nr}.html` |
| `PRUEF` | `https://ris.bka.gv.at/Dokumente/PruefGewO/{nr}/{nr}.html` |
| `AVSV` | `https://ris.bka.gv.at/Dokumente/Avsv/{nr}/{nr}.html` |
| `SPG` | `https://ris.bka.gv.at/Dokumente/Spg/{nr}/{nr}.html` |
| `KMGER` | `https://ris.bka.gv.at/Dokumente/KmGer/{nr}/{nr}.html` |

**Important:** Check longer prefixes first (e.g., `BGBLA` before `BGBL`, `ASYLGH` before any shorter match).

### Prefix → Search API Fallback

When direct URL fails, route to search API:

| Prefix | Endpoint | Applikation |
|--------|----------|-------------|
| `NOR` | `Bundesrecht` | `BrKons` |
| `LBG`, `LNO`, `LST`, `LTI`, `LVO`, `LWI`, `LSB`, `LOO`, `LKT` | `Landesrecht` | `LrKons` |
| `JFR`, `JFT` | `Judikatur` | `Vfgh` |
| `JWR`, `JWT` | `Judikatur` | `Vwgh` |
| `BVWG` | `Judikatur` | `Bvwg` |
| `LVWG` | `Judikatur` | `Lvwg` |
| `DSB` | `Judikatur` | `Dsk` |
| `GBK` | `Judikatur` | `Gbk` |
| `PVAK` | `Judikatur` | `Pvak` |
| `ASYLGH` | `Judikatur` | `AsylGH` |
| `BGBLA` | `Bundesrecht` | `BgblAuth` |
| `BGBL` | `Bundesrecht` | `BgblAlt` |
| `REGV` | `Bundesrecht` | `RegV` |
| `MRP` | `Sonstige` | `Mrp` |
| `ERL` | `Sonstige` | `Erlaesse` |
| (unknown) | `Judikatur` | `Justiz` (fallback) |

Search params: `Dokumentnummer={nr}`, `DokumenteProSeite=Ten`

## 6. JSON Output Schema

When `--json` is used, output stable JSON to stdout.

### Search Results

```json
{
  "total_hits": 42,
  "page": 1,
  "page_size": 20,
  "has_more": true,
  "documents": [
    {
      "dokumentnummer": "NOR40052761",
      "applikation": "BrKons",
      "titel": "§ 1295 ABGB",
      "kurztitel": "ABGB",
      "citation": {
        "kurztitel": "ABGB",
        "langtitel": "Allgemeines buergerliches Gesetzbuch",
        "kundmachungsorgan": "JGS Nr. 946/1811",
        "paragraph": "§ 1295",
        "eli": "...",
        "inkrafttreten": "1812-01-01",
        "ausserkrafttreten": null
      },
      "content_urls": {
        "html": "https://...",
        "xml": "https://...",
        "pdf": null,
        "rtf": null
      },
      "dokument_url": "https://...",
      "gesamte_rechtsvorschrift_url": "https://..."
    }
  ]
}
```

### Document Content

```json
{
  "metadata": {
    "dokumentnummer": "NOR40052761",
    "applikation": "BrKons",
    "titel": "§ 1295 ABGB",
    "kurztitel": "ABGB",
    "citation": { "..." },
    "dokument_url": "...",
    "gesamte_rechtsvorschrift_url": "..."
  },
  "content": "Full text of the document..."
}
```

## 7. Go Project Structure

```
ris-cli/
├── PLAN.md                    # This file
├── README.md
├── LICENSE
├── go.mod
├── go.sum
├── main.go                    # Entry point, root cobra command
├── cmd/
│   ├── root.go                # Root command, global flags
│   ├── bundesrecht.go
│   ├── landesrecht.go
│   ├── judikatur.go
│   ├── bgbl.go
│   ├── lgbl.go
│   ├── regvorl.go
│   ├── dokument.go
│   ├── bezirke.go
│   ├── gemeinden.go
│   ├── sonstige.go            # Parent + sub-subcommands (mrp, erlaesse, etc.)
│   ├── history.go
│   ├── verordnungen.go
│   └── version.go
├── internal/
│   ├── api/
│   │   ├── client.go          # HTTP client, base URL, timeout, error types
│   │   ├── params.go          # Query parameter builder
│   │   └── endpoints.go       # Endpoint constants
│   ├── parser/
│   │   ├── response.go        # Raw API response types
│   │   ├── document.go        # Document parsing from raw response
│   │   └── search.go          # SearchResult parsing
│   ├── model/
│   │   ├── document.go        # Document, Citation, ContentUrl types
│   │   ├── search.go          # SearchResult type
│   │   └── routing.go         # Document prefix routing tables
│   ├── format/
│   │   ├── text.go            # Human-readable terminal output
│   │   ├── json.go            # JSON output
│   │   ├── html.go            # HTML-to-text conversion
│   │   └── citation.go        # Austrian legal citation formatting
│   └── constants/
│       ├── states.go           # Bundesland mappings
│       ├── courts.go           # Court type mappings
│       ├── ministries.go       # Ministry mappings
│       └── applications.go     # App-specific enum values
└── test/
    ├── api_test.go
    ├── parser_test.go
    ├── format_test.go
    └── integration/
        └── smoke_test.go
```

### Key Go Dependencies

| Package | Purpose |
|---------|---------|
| `github.com/spf13/cobra` | CLI framework (subcommands, flags, help generation) |
| `github.com/fatih/color` | Terminal colors (respects `NO_COLOR`) |
| `github.com/olekukonez/tablewriter` | Table formatting for search results |
| `golang.org/x/net/html` | HTML parsing (stdlib-adjacent, for html-to-text) |

No external HTTP client needed — Go's `net/http` stdlib is sufficient.

## 8. Phased Implementation Plan

### Phase 1 — Skeleton + Core (MVP)

**Goal:** `ris bundesrecht --search "Mietrecht" --json` works end-to-end.

1. Initialize Go module: `go mod init github.com/philrox/ris-cli`
2. Set up cobra root command with global flags
3. Implement `internal/api/client.go` — HTTP client with timeout, error types
4. Implement `internal/parser/` — Raw response parsing
5. Implement `internal/model/` — Document and SearchResult types
6. Implement `internal/format/json.go` — JSON output
7. Implement `cmd/bundesrecht.go` — First subcommand
8. Implement `cmd/dokument.go` — Document retrieval with prefix routing
9. Basic tests for parser and formatting

**Deliverable:** Working binary with `bundesrecht` + `dokument` + `--json`.

### Phase 2 — All Search Commands

1. Implement `cmd/judikatur.go`
2. Implement `cmd/landesrecht.go`
3. Implement `cmd/bgbl.go`
4. Implement `cmd/lgbl.go`
5. Implement `cmd/regvorl.go`
6. Implement `cmd/bezirke.go`
7. Implement `cmd/gemeinden.go`
8. Implement `cmd/sonstige.go` (with 8 sub-subcommands)
9. Implement `cmd/history.go`
10. Implement `cmd/verordnungen.go`

**Deliverable:** All 12 commands functional.

### Phase 3 — Human-Friendly Output

1. Implement `internal/format/text.go` — Colored terminal output
2. Implement `internal/format/citation.go` — Austrian legal citations
3. Implement `internal/format/html.go` — HTML-to-text for documents
4. Add progress spinner on stderr for network requests
5. Add `--plain` output mode (stable, no colors, line-based)
6. TTY detection: auto-select text vs plain based on `isatty(stdout)`
7. Pager support for long document output (`--pager` / `PAGER` env)

**Deliverable:** Beautiful terminal UX for humans, stable output for scripts.

### Phase 4 — Polish & Distribution

1. Shell completions (bash, zsh, fish, powershell — cobra generates these)
2. `ris version` command
3. GitHub Actions CI/CD (test matrix, cross-compile, release)
4. GoReleaser config for multi-platform binaries
5. Homebrew formula
6. `go install github.com/philrox/ris-cli@latest` support
7. Man page generation (cobra-based)
8. README with examples

**Deliverable:** Production-ready, installable via `brew install ris` or single binary download.

## 9. Environment & Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `RIS_TIMEOUT` | HTTP timeout | `30s` |
| `RIS_BASE_URL` | API base URL override | `https://data.bka.gv.at/ris/api/v2.6/` |
| `NO_COLOR` | Disable colors (standard) | — |
| `PAGER` | Pager for long output | `less -FIRX` |

### Precedence

```
flags > env vars > defaults
```

No config file needed — the CLI is stateless. All configuration via flags and env vars.

## 10. Example Invocations

```bash
# Search federal law for rent law, get first page as JSON
ris bundesrecht --search "Mietrecht" --json

# Pipe to jq for specific fields
ris bundesrecht --search "Mietrecht" --json | jq '.documents[].dokumentnummer'

# Get a specific ABGB paragraph
ris bundesrecht --title "ABGB" --paragraph 1295

# Retrieve full document text
ris dokument NOR40052761

# Retrieve document as JSON (for agents)
ris dokument NOR40052761 --json | jq '.content'

# Search Constitutional Court decisions about fundamental rights
ris judikatur --search "Grundrecht" --court vfgh --from 2020-01-01

# Search Salzburg state law
ris landesrecht --search "Bauordnung" --state salzburg

# Search Federal Law Gazette
ris bgbl --number 120 --year 2023 --part 1

# Search government bills from Finance Ministry
ris regvorl --ministry bmf --from 2024-01-01

# Search Council of Ministers protocols
ris sonstige mrp --search "Budget" --session 42

# Check document change history
ris history --app bundesnormen --from 2024-01-01 --to 2024-01-31

# Agent workflow: search, pick first, get full text
DOC=$(ris bundesrecht --search "Datenschutz" --json | jq -r '.documents[0].dokumentnummer')
ris dokument "$DOC" --json | jq '.content'

# Pagination
ris judikatur --search "Schadenersatz" --page 2 --limit 50
```

## 11. Safety & Robustness

- **SSRF Protection:** Only fetch document content from allowed hosts: `data.bka.gv.at`, `www.ris.bka.gv.at`, `ris.bka.gv.at` (HTTPS only)
- **Input Validation:** Document numbers validated against `^[A-Z][A-Z0-9_]+$` (5-50 chars) before any URL construction
- **Timeout:** Default 30s, configurable via `--timeout` / `RIS_TIMEOUT`
- **No secrets:** CLI has no authentication — RIS API is public
- **Crash-only:** No state to corrupt. Safe to Ctrl-C at any point
- **Idempotent:** All operations are read-only GET requests
