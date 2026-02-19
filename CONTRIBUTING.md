# Contributing to ris-cli

Thanks for your interest in contributing to **ris-cli**, a Go CLI tool for searching Austrian legal documents via the [RIS API](https://data.bka.gv.at/ris/api/v2.6/).

## Prerequisites

- [Go](https://go.dev/dl/) 1.24 or later
- [git](https://git-scm.com/)
- A GitHub account

## Getting Started

1. Fork this repository on GitHub.
2. Clone your fork:
   ```bash
   git clone https://github.com/<your-user>/ris-cli.git
   cd ris-cli
   ```
3. Build the binary:
   ```bash
   go build -o ris .
   ```
4. Run the tests:
   ```bash
   go test ./...
   ```

If everything passes you are ready to go.

## Development

The project ships a `Makefile` with common shortcuts:

| Command       | Description                           |
|---------------|---------------------------------------|
| `make build`  | Compile the `ris` binary              |
| `make test`   | Run tests with race detector          |
| `make lint`   | Run `go vet`                          |
| `make fmt`    | Format code with `go fmt`             |
| `make check`  | Run fmt, lint, and test in sequence   |
| `make clean`  | Remove the compiled binary            |
| `make install`| Install `ris` into your `$GOPATH/bin` |

Before opening a pull request, run `make check` to make sure everything is in order.

## Project Structure

```
.
├── main.go                 # Entry point
├── cmd/                    # CLI commands (cobra)
├── internal/
│   ├── api/                # HTTP client for the RIS API
│   ├── parser/             # Response parsing
│   ├── model/              # Shared types and structs
│   ├── format/             # Output formatting (table, detail views)
│   ├── constants/          # Enum mappings and named constants
│   └── ui/                 # Terminal UI helpers (spinner, colors)
├── Makefile                # Developer shortcuts
├── go.mod / go.sum         # Go module files
└── ...
```

## Pull Request Process

1. Create a feature branch off `main`:
   ```bash
   git checkout -b feat/my-change
   ```
2. Make your changes. Keep commits small and focused.
3. Use **conventional commit** style messages:
   - `feat: add --output json flag`
   - `fix: handle empty API response`
   - `docs: update README examples`
   - `refactor: extract search helper`
   - `test: add parser unit tests`
4. Run `make check` and make sure everything passes.
5. Push your branch and open a pull request against `main`.

## Code Style

- Run `go fmt` and `go vet` before committing.
- Keep functions short and names descriptive.
- Prefer returning errors over panicking.
- Write tests for new functionality (see `*_test.go` files for examples).

## Language Convention

The CLI output and help text are in **German** because the tool targets Austrian legal professionals. Code, comments, commit messages, and PR discussions should be in **English**.

## Reporting Bugs

If you find a bug, please [open a GitHub Issue](https://github.com/philrox/ris-cli/issues/new) with:

- A clear description of the problem
- Steps to reproduce
- Expected vs. actual behavior
- Your Go version (`go version`) and OS

## Questions?

Feel free to open an issue if something is unclear. We appreciate every contribution, no matter how small.
