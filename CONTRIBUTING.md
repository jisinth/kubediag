# Contributing to kubediag

Thanks for considering a contribution!

## Getting started

```
git clone https://github.com/jisinth/kubediag.git
cd kubediag
go mod tidy
go build ./...
```

## Development workflow

1. Open an issue describing the bug or feature before starting significant work.
2. Create a branch off `main`.
3. Make your change, following the existing package layout under `cmd/` and `internal/`
   (see [docs/architecture.md](docs/architecture.md)).
4. Run the checks CI will run:
   ```
   gofmt -l .
   go vet ./...
   go test ./...
   go build ./...
   ```
5. Open a pull request describing the change and, for new diagnostic checks, the failure
   scenario it detects and the recommendation it produces.

## Adding a diagnostic check

Diagnostic checks live under `internal/<module>` and return `[]reporter.Issue`. Root causes and
their recommended fixes are centralized in `internal/recommendations`, so add a new
`recommendations.Code` there rather than hardcoding advice inline.

## Code style

- Run `gofmt` before committing.
- Keep exported functions documented with a doc comment.
- Prefer table-driven tests for new checks.

## Reporting issues

Use GitHub Issues. Include your Kubernetes version, the `kubediag` command you ran, and (if
possible) anonymized output.
