# RIS CLI — Refactoring-Plan

Erstellt: 2026-02-18
Basierend auf automatisierter Analyse mit 3 parallelen Subagents (Struktur, Code-Qualität, Bug-Review).

---

## Zusammenfassung

| Kategorie | Anzahl Issues | Betroffene Dateien |
|-----------|---------------|-------------------|
| Bugs | 2 | 2 |
| Error-Handling-Refactoring | 1 (43 Callsites) | 12 |
| DRY-Refactoring | 2 | 11 |
| Dead Code | 6 | 5 |
| Hardcoded Values | 5 | 4 |
| Konsistenz-Fixes | 3 | 4 |

Geschätzter Gesamtumfang: ~250 geänderte Zeilen, ~120 gelöschte Zeilen.

---

## Phase 1: Bug-Fixes (Korrektheit)

### 1.1 JWR-Prefix falsch geroutet

**Problem:** VwGH-Dokumente mit Prefix `JWR` werden im Such-Fallback unter `Applikation=Justiz` statt `Applikation=Vwgh` gesucht. Die direkte URL funktioniert, aber der Fallback liefert falsche oder keine Ergebnisse.

**Datei:** `internal/model/routing.go:41`

**Vorher:**
```go
{"JWR", DocumentRoute{URLPath: "Vwgh", Endpoint: "Judikatur", Applikation: "Justiz"}},
```

**Nachher:**
```go
{"JWR", DocumentRoute{URLPath: "Vwgh", Endpoint: "Judikatur", Applikation: "Vwgh"}},
```

**Vergleich mit korrekten Einträgen:**
- `JFR` → `Applikation: "Vfgh"` (korrekt)
- `JFT` → `Applikation: "Vfgh"` (korrekt)
- `JWT` → `Applikation: "Justiz"` (korrekt — Justiz-Dokumente)
- `JWR` → `Applikation: "Justiz"` (BUG — VwGH-Dokumente)

**Validierung:** `ris dokument JWR_...` mit einem echten VwGH-Dokumentnummer testen. Ohne Fix: Fallback findet nichts. Mit Fix: Fallback findet das Dokument.

---

### 1.2 Duplizierte SSRF-Allowlist

**Problem:** Zwei unabhängige Allowlists für erlaubte Hosts. Bei Änderung muss man an zwei Stellen updaten — Risiko für Drift.

**Dateien:**
- `cmd/dokument.go:41-45` — lokale `allowedHosts` Map
- `internal/api/endpoints.go:18-22` — `AllowedHosts` (kanonische Quelle)

**Lösung:**

1. `cmd/dokument.go:41-45` löschen (lokale `allowedHosts` Map entfernen)

2. `cmd/dokument.go:57-70` `validateURL` auf `api.AllowedHosts` umstellen:

```go
func validateURL(rawURL string) error {
	u, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("Ungültige URL: %w", err)
	}
	if u.Scheme != "https" {
		return fmt.Errorf("Nur HTTPS-URLs erlaubt, erhalten: %q", u.Scheme)
	}
	host := strings.ToLower(u.Hostname())
	if !api.AllowedHosts[host] {
		return fmt.Errorf("Host %q nicht erlaubt", host)
	}
	return nil
}
```

3. `internal/api/client.go:135` Fehler-Nachricht dynamisch aus `AllowedHosts` generieren statt Hosts hardzucoden.

**Validierung:** `go build ./...` — Compile-Check reicht, da nur Import-Pfad sich ändert.

---

## Phase 2: Error-Handling (`os.Exit` → `return error`)

### 2.1 Problem

Alle `RunE`-Handler verwenden `os.Exit(2)` für Validierungsfehler statt `return error`. Das:
- Umgeht `defer`-Cleanup (konkretes Problem: Pager-Prozess in `dokument.go:155,189`)
- Macht Commands nicht als Library testbar
- Ignoriert Cobras `SilenceUsage`/`SilenceErrors`-System

**43 Callsites in 12 Dateien:**

| Datei | os.Exit-Aufrufe | Zeilen |
|-------|-----------------|--------|
| `cmd/bundesrecht.go` | 2 | 49, 56 |
| `cmd/landesrecht.go` | 2 | 43, 62 |
| `cmd/judikatur.go` | 2 | 50, 58 |
| `cmd/bgbl.go` | 3 | 49, 55, 78 |
| `cmd/lgbl.go` | 3 | 49, 55, 73 |
| `cmd/regvorl.go` | 5 | 53, 77, 85, 93, 101 |
| `cmd/bezirke.go` | 3 | 53, 72, 92 |
| `cmd/gemeinden.go` | 7 | 70, 76, 102, 128, 136, 144, 137* |
| `cmd/verordnungen.go` | 2 | 49, 68 |
| `cmd/history.go` | 3 | 45, 51, 59 |
| `cmd/dokument.go` | 4 | 81, 90, 98, 137 |
| `cmd/sonstige.go` | 8 | 278, 319, 355, 394, 428, 461, 469, 499 |

### 2.2 Lösung

Neuen Error-Typ für Validierungsfehler in `cmd/helpers.go` einführen:

```go
// validationError represents a user input validation error.
// main.go uses this to set exit code 2.
type validationError struct {
	msg string
}

func (e *validationError) Error() string { return e.msg }

// errValidation creates a validation error with fmt.Sprintf formatting.
func errValidation(format string, args ...any) error {
	return &validationError{msg: fmt.Sprintf(format, args...)}
}
```

`main.go` anpassen, um exit code 2 für Validierungsfehler zu setzen:

```go
func main() {
	if err := cmd.Execute(); err != nil {
		var ve *cmd.ValidationError
		if errors.As(err, &ve) {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
```

**Muster für jeden Command-Handler:**

Vorher:
```go
if search == "" && title == "" {
    fmt.Fprintln(os.Stderr, "Fehler: mindestens --search oder --title erforderlich")
    os.Exit(2)
}
```

Nachher:
```go
if search == "" && title == "" {
    return errValidation("Fehler: mindestens --search oder --title erforderlich")
}
```

### 2.3 Spezialfall `dokument.go:137`

Aktuell `os.Exit(3)` für "Dokument nicht gefunden". Diesen Fall ebenfalls als Error modellieren:

```go
if len(result.Documents) == 0 {
    return errValidation("Fehler: Dokument %q nicht gefunden", docNumber)
}
```

(Exit-Code 3 war nicht dokumentiert und inkonsistent — wir vereinheitlichen auf 2 für alle User-Fehler.)

### 2.4 Umsetzungsreihenfolge

1. `cmd/helpers.go` — `validationError` Typ + `errValidation()` hinzufügen
2. `main.go` — Exit-Code-Logik anpassen
3. Jeden Command-Handler systematisch umstellen (alphabetisch)
4. `go build ./...` nach jeder Datei
5. Am Ende: `grep -rn "os.Exit" cmd/` — muss 0 Treffer liefern

**Validierung:** `go vet ./...` und `go build ./...`. Manuell testen: `ris bundesrecht` (ohne Flags) → sollte Exit-Code 2 + Fehlermeldung auf stderr liefern.

---

## Phase 3: DRY — Search-Execute-Helper extrahieren

### 3.1 Problem

Der identische 16-Zeilen-Block (Spinner → API-Call → Parse → Output) ist in 9 Command-Dateien kopiert:

```go
s := startSpinner(cmd, "Suche in ...")
body, err := client.Search("<endpoint>", params)
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
```

`sonstige.go` hat dies bereits korrekt mit `executeSonstigeSearch` gelöst — das Muster wird verallgemeinert.

### 3.2 Lösung

In `cmd/helpers.go` eine generische Funktion hinzufügen:

```go
// executeSearch runs the common search pipeline: spinner → API call → parse → output.
func executeSearch(cmd *cobra.Command, endpoint, spinnerMsg string, params *api.Params) error {
	setPageParams(cmd, params)

	client := newClient(cmd)
	s := startSpinner(cmd, spinnerMsg)
	body, err := client.Search(endpoint, params)
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
```

### 3.3 Betroffene Dateien

Jede Datei wird auf `return executeSearch(cmd, endpoint, msg, params)` reduziert:

| Datei | Endpoint | Spinner-Nachricht |
|-------|----------|-------------------|
| `bundesrecht.go` | `"Bundesrecht"` | `"Suche in Bundesrecht..."` |
| `landesrecht.go` | `"Landesrecht"` | `"Suche in Landesrecht..."` |
| `judikatur.go` | `"Judikatur"` | `"Suche in Judikatur..."` |
| `bgbl.go` | `"Bundesrecht"` | `"Suche in Bundesgesetzblättern..."` |
| `lgbl.go` | `"Landesrecht"` | `"Suche in Landesgesetzblättern..."` |
| `regvorl.go` | `"Bundesrecht"` | `"Suche in Regierungsvorlagen..."` |
| `bezirke.go` | `"Bezirke"` | `"Suche in Bezirksverwaltung..."` |
| `gemeinden.go` | `"Gemeinden"` | `"Suche in Gemeinderecht..."` |
| `verordnungen.go` | `"Landesrecht"` | `"Suche in Verordnungsblättern..."` |
| `history.go` | `"History"` | `"Suche in Änderungshistorie..."` |

### 3.4 `executeSonstigeSearch` refactoren

`sonstige.go:executeSonstigeSearch` wird auf `executeSearch` umgestellt:

```go
func executeSonstigeSearch(cmd *cobra.Command, params *api.Params) error {
	return executeSearch(cmd, "Sonstige", "Suche in Sonstige Rechtsquellen...", params)
}
```

**Alternativ:** `executeSonstigeSearch` komplett entfernen und alle 8 sonstige-Subcommands direkt `executeSearch` aufrufen. Empfehlung: `executeSonstigeSearch` beibehalten als Wrapper, da es `setCommonSonstigeParams` nicht aufruft (das passiert in jedem run*).

### 3.5 Import-Cleanup

Nach dem Refactoring können in jedem Command-File die Imports für `parser`, `format` und `os` entfernt werden (da diese nur im duplizierten Block verwendet wurden):

Vorher (bundesrecht.go):
```go
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
```

Nachher (bundesrecht.go):
```go
import (
    "strings"
    "github.com/philrox/ris-cli/internal/api"
    "github.com/philrox/ris-cli/internal/constants"
    "github.com/spf13/cobra"
)
```

Hinweis: `"fmt"` und `"os"` werden nach Phase 2 möglicherweise noch für `errValidation` gebraucht — nach Phase 2 prüfen, ob `"fmt"` noch nötig ist (vermutlich nicht, da `errValidation` in helpers.go lebt). `"os"` wird definitiv nicht mehr gebraucht.

**Validierung:** `go build ./...` und `go vet ./...`. Manuell testen: `ris bundesrecht --search "ABGB"` — Output muss identisch sein.

---

## Phase 4: Dead Code entfernen

### 4.1 `ParseDocumentResponse` — nie aufgerufen

**Datei:** `internal/parser/document.go` (komplett)

**Aktion:** Gesamte Datei löschen. Die Funktion ist ein Wrapper um `ParseSearchResponse`, der nirgends importiert wird. `cmd/dokument.go` verwendet `ParseSearchResponse` direkt.

### 4.2 `FlexibleStringArray` — nie verwendet

**Datei:** `internal/parser/flexible.go:98-118`

**Aktion:** Typ und `UnmarshalJSON`-Methode löschen (Zeilen 98-118).

### 4.3 `validateLimit` — nie aufgerufen

**Datei:** `cmd/root.go:92-100`

**Aktion:** Funktion löschen. Die Limit-Validierung passiert implizit in `setPageParams` via `constants.PageSizes` Lookup.

### 4.4 `LeitsatzCourts` — nie referenziert

**Datei:** `internal/constants/courts.go:19-25`

**Aktion:** Map löschen. Die Leitsatz-Logik ist inline in `parser/search.go:180-195` via `judApp.hasLeitsatz` gelöst.

### 4.5 `ValidCourts` — nie aufgerufen

**Datei:** `internal/constants/courts.go:27-33`

**Aktion:** Funktion löschen. Fehlermeldungen in `cmd/judikatur.go:57` listen die gültigen Werte als Hardcoded-String.

### 4.6 `ValidRegvorlMinistries` und `ValidErlasseMinistries` — nie aufgerufen

**Datei:** `internal/constants/ministries.go:39-53`

**Aktion:** Beide Funktionen löschen.

### 4.7 Zusammenfassung Dead Code

| Was | Datei | Aktion | Gelöschte Zeilen |
|-----|-------|--------|-----------------|
| `ParseDocumentResponse` | `parser/document.go` | Datei löschen | 19 |
| `FlexibleStringArray` | `parser/flexible.go` | Zeilen 98-118 löschen | 21 |
| `validateLimit` | `cmd/root.go` | Zeilen 92-100 löschen | 9 |
| `LeitsatzCourts` | `constants/courts.go` | Zeilen 19-25 löschen | 7 |
| `ValidCourts` | `constants/courts.go` | Zeilen 27-33 löschen | 7 |
| `ValidRegvorlMinistries` | `constants/ministries.go` | Zeilen 39-45 löschen | 7 |
| `ValidErlasseMinistries` | `constants/ministries.go` | Zeilen 47-53 löschen | 7 |
| **Gesamt** | | | **77** |

**Validierung:** `go build ./...` nach jeder Löschung. Kein Test-Impact, da die Funktionen nie aufgerufen werden.

---

## Phase 5: Hardcoded Values → Named Constants

### 5.1 Magic String `"9999-12-31"`

**Datei:** `internal/parser/search.go:124,154`

**Lösung:** Konstante am Dateianfang:

```go
const noExpiryDate = "9999-12-31"
```

Ersetzen an beiden Stellen:
```go
// Zeile 124 und 154
if akt != "" && akt != noExpiryDate {
```

### 5.2 Magic String `"Ten"` in `dokument.go`

**Datei:** `cmd/dokument.go:121`

Vorher:
```go
params.Set("DokumenteProSeite", "Ten")
```

Nachher:
```go
params.Set("DokumenteProSeite", constants.PageSizes[10])
```

### 5.3 Magic String `"Twenty"` in `helpers.go`

**Datei:** `cmd/helpers.go:76`

Vorher:
```go
params.Set("DokumenteProSeite", "Twenty")
```

Nachher:
```go
params.Set("DokumenteProSeite", constants.PageSizes[20])
```

### 5.4 Magic Number `200` (Leitsatz-Truncation)

**Datei:** `internal/format/text.go:62`

**Lösung:** Konstante in `text.go`:

```go
const maxLeitsatzPreview = 200
```

```go
if len(leitsatz) > maxLeitsatzPreview {
    leitsatz = leitsatz[:maxLeitsatzPreview] + "..."
}
```

### 5.5 Magic Number `60` (Separator-Breite)

**Datei:** `internal/format/text.go:33,105`

**Lösung:** Konstante in `text.go`:

```go
const separatorWidth = 60
```

```go
fmt.Fprintln(w, dim(strings.Repeat("─", separatorWidth)))
```

**Validierung:** `go build ./...`. Unit-Tests: `go test ./internal/format/...` — bestehende Tests müssen weiter grün sein.

---

## Phase 6: Konsistenz-Fixes

### 6.1 `setCommonSonstigeParams` — stille Fehlerunterdrückung

**Datei:** `cmd/sonstige.go:184-195`

**Problem:** Ungültige `--since` und `--sort-dir` Werte werden still ignoriert (kein Fehler).
Alle anderen Commands (regvorl, bezirke, gemeinden) geben Fehler aus.

**Lösung:** `setCommonSonstigeParams` gibt `error` zurück:

```go
func setCommonSonstigeParams(cmd *cobra.Command, params *api.Params) error {
	search, _ := cmd.Flags().GetString("search")
	title, _ := cmd.Flags().GetString("title")
	since, _ := cmd.Flags().GetString("since")
	sortDir, _ := cmd.Flags().GetString("sort-dir")

	if search != "" {
		params.Set("Suchworte", search)
	}
	if title != "" {
		params.Set("Titel", title)
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
	return nil
}
```

Alle 8 Aufrufe in `sonstige.go` anpassen:

```go
if err := setCommonSonstigeParams(cmd, params); err != nil {
    return err
}
```

### 6.2 `rawBrKons` / `rawLrKons` zu einem Typ zusammenführen

**Datei:** `internal/parser/response.go:72-103`

**Problem:** Zwei strukturell identische Typen mit identischen JSON-Tags.

**Lösung:** Einen Typ `rawSubApp` definieren, der beide ersetzt:

```go
// rawSubApp is a sub-application section for both Bundesrecht and Landesrecht.
type rawSubApp struct {
	Kundmachungsorgan          string         `json:"Kundmachungsorgan"`
	ArtikelParagraphAnlage     FlexibleString `json:"ArtikelParagraphAnlage"`
	Inkrafttretensdatum        string         `json:"Inkrafttretensdatum"`
	Ausserkrafttretensdatum    string         `json:"Ausserkrafttretensdatum"`
	GesamteRechtsvorschriftURL string         `json:"GesamteRechtsvorschriftUrl"`
}
```

Referenzen aktualisieren:
- `response.go`: `rawBundesrecht` — alle `*rawBrKons` → `*rawSubApp`
- `response.go`: `rawLandesrecht` — alle `*rawLrKons` → `*rawSubApp`
- `search.go:240-256`: `firstNonNil` und `firstNonNilLR` → eine Funktion `firstNonNil(ptrs ...*rawSubApp) *rawSubApp`
- `search.go:118`: Aufruf von `firstNonNil` bleibt gleich
- `search.go:148`: Aufruf von `firstNonNilLR` → `firstNonNil`

### 6.3 Inkonsistenter Flag-Zugriff (Globals vs. Cobra-API)

**Dateien:** `cmd/helpers.go:16-17` vs. `cmd/root.go:13-27`

**Problem:** `newClient` und `isVerbose` lesen Flags via Cobra-API (`root.PersistentFlags().GetBool`), obwohl die gleichen Werte bereits als Package-Level-Variablen in `root.go` existieren (`verbose`, `timeout`).

**Lösung:** `newClient` und `isVerbose` direkt auf die Package-Vars zugreifen:

```go
func newClient(cmd *cobra.Command) *api.Client {
	return api.NewClient(api.ClientOptions{
		Timeout: timeout,
		Verbose: verbose,
	})
}

func isVerbose(cmd *cobra.Command) bool {
	return verbose
}
```

Dadurch wird der `cmd`-Parameter bei `isVerbose` überflüssig, kann aber für API-Kompatibilität beibehalten werden. Alternativ: Signatur zu `isVerbose() bool` ändern und alle Callsites anpassen (nur `dokument.go:111`).

**Validierung:** `go build ./...` und `go vet ./...`.

---

## Phasen-Übersicht & Abhängigkeiten

```
Phase 1: Bug-Fixes
  ├── 1.1 JWR-Routing Fix           (unabhängig)
  └── 1.2 SSRF-Allowlist deduplizieren (unabhängig)

Phase 2: Error-Handling              (unabhängig von Phase 1)
  └── os.Exit → return error in allen Commands

Phase 3: DRY-Refactoring            (ABHÄNGIG von Phase 2, da os.Exit weg sein muss)
  └── executeSearch Helper extrahieren

Phase 4: Dead Code entfernen         (unabhängig)

Phase 5: Named Constants             (unabhängig)

Phase 6: Konsistenz-Fixes            (6.1 ABHÄNGIG von Phase 2 für errValidation)
  ├── 6.1 sonstige Fehlerbehandlung
  ├── 6.2 rawSubApp Vereinigung     (unabhängig)
  └── 6.3 Flag-Zugriff              (unabhängig)
```

**Empfohlene Reihenfolge:**
1. Phase 1 (Bug-Fixes) — sofort, da Korrektheitsprobleme
2. Phase 4 (Dead Code) — einfach, reduziert Rauschen für folgende Phasen
3. Phase 2 (Error-Handling) — Grundlage für Phase 3 und 6.1
4. Phase 3 (DRY) — größte Verbesserung, braucht Phase 2
5. Phase 5 (Constants) — Cleanup
6. Phase 6 (Konsistenz) — Abschluss

---

## Validierungsstrategie pro Phase

| Phase | Validierung |
|-------|-------------|
| 1 | `go build ./...`, manueller Test mit VwGH-Dokument |
| 2 | `go build ./...`, `go vet ./...`, `grep -rn "os.Exit" cmd/` = 0 Treffer, `ris bundesrecht` → Exit 2 |
| 3 | `go build ./...`, `go test ./...`, `ris bundesrecht --search "ABGB"` Output identisch |
| 4 | `go build ./...`, `go vet ./...` |
| 5 | `go build ./...`, `go test ./internal/format/...` |
| 6 | `go build ./...`, `go test ./...` |

Nach jeder Phase: `go test ./...` — alle bestehenden Tests müssen grün bleiben.

---

## Nicht im Scope (bewusst ausgeklammert)

- **`PAGER` Environment Variable Sanitization** — Standard-Pattern bei CLI-Tools (git macht es genauso), User-kontrolliert, kein Sicherheitsproblem.
- **`FlexibleArray` Hard-Error bei unbekanntem JSON-Typ** — API ist stabil definiert, theoretisches Risiko zu niedrig.
- **`hasMore` Paginierungs-Berechnung** — Funktioniert korrekt für alle praktischen Fälle, Edge-Case nur bei letzter Seite mit weniger Ergebnissen als PageSize.
- **Typed API-Parameter-Namen** — Overengineering für die aktuelle Codebase-Größe.
