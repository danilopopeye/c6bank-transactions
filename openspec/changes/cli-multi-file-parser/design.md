## Context

The project is a Go web service (`cmd/c6bank-transactions`) that processes C6 Bank transaction files (PDF, CSV, mobile screenshots) and converts them to QIF/CSV. The parser logic lives in `internal/parser/` with format-specific scanners (`scanPDFRows`, `scanCSVRows`, `ScanImage`) that extract transaction data.

Currently the `Parse()` function couples file reading (via `multipart.File`) with output generation (QIF or CSV). The CLI needs to call the lower-level scanners directly to accumulate transactions from multiple files before producing a single output.

## Goals / Non-Goals

**Goals:**
- Provide a CLI that accepts multiple file paths as positional arguments
- Parse each file using existing format-specific scanners
- Deduplicate transactions by hash (Date + Payee + Amount + Memo)
- Output a single CSV to stdout (or a file via `-o` flag)

**Non-Goals:**
- QIF output (the CLI targets CSV accumulation only)
- Modifying the existing web server or parser package API
- Watching directories or glob expansion (users pass explicit paths)
- Parallel file processing (sequential is sufficient for typical usage)

## Decisions

### 1. CLI entry point at `cmd/cli/`

**Decision**: Create a separate binary under `cmd/cli/` rather than adding subcommands to the existing server.

**Rationale**: Keeps the web server and CLI as independent binaries. The Go toolchain naturally supports multiple `cmd/` directories. No subcommand framework needed — this CLI does one thing.

**Alternative considered**: Adding a `cli` subcommand to the existing server binary — rejected to avoid coupling and reduce binary size.

### 2. Reuse scanners directly, not `Parse()`

**Decision**: Call `scanCSVRows` and `ScanImage` directly instead of the high-level `Parse()` function.

**Rationale**: `Parse()` couples input (multipart.File) with output (QIF generation). The CLI needs the raw `[]Line` or `[]Transaction` data to accumulate across files. The scanners are already unexported but in the same package — we will expose new exported functions that wrap them for CLI use. PDF support was excluded from the CLI to keep the binary focused on the most common batch processing formats (CSV and mobile screenshots).

### 3. Exported function `ParseFile(path string) ([]Transaction, error)`

**Decision**: Add a single exported function to the parser package that accepts a file path and returns `[]Transaction`.

**Rationale**: Provides a clean API for the CLI without exposing internal scanner details. The function handles opening the file, detecting format, calling the appropriate scanner, and converting `[]Line` to `[]Transaction`.

### 4. Deduplication via FNV-1a hash

**Decision**: Use the existing FNV-1a hash approach (Date + Payee + Amount + Memo) to identify and remove duplicate transactions.

**Rationale**: The same transaction can appear in overlapping statements. Including Memo avoids collapsing distinct transactions that share date/payee/amount. The project already depends on `github.com/segmentio/fasthash/fnv1a`.

### 5. Output to stdout by default, `-o` flag for file

**Decision**: Write CSV to stdout. Optional `-o <path>` flag writes to a file.

**Rationale**: Unix convention — stdout allows piping. `-o` flag for convenience when users want to save directly.

### 6. Fail-fast on any file error

**Decision**: If any file fails to parse (missing, wrong format, validation error), the CLI exits immediately with a non-zero code.

**Rationale**: Silent partial results could mislead users. Explicit failure lets them fix the problem and re-run.

### 7. Chronological output ordering

**Decision**: Sort accumulated transactions by date (chronological) before writing CSV.

**Rationale**: Multi-file input produces interleaved transaction order. Chronological output is the natural expectation for financial data.

## Risks / Trade-offs

- **Unexported scanners**: Need to either export them or create a new exported wrapper. Creating a wrapper is cleaner but adds a function to the parser package API. → *Mitigation*: Single `ParseFile()` function keeps the surface small.
- **CSV filename validation**: `scanCSVRows` requires filename in `Fatura_YYYY-MM-DD.csv` format. For CLI usage, the filename comes from the file path, which should work naturally if users pass properly named files. → *Mitigation*: Document the naming requirement in CLI help text.
- **No glob expansion**: Users must pass explicit file paths. Shell globbing (`*.csv`) handles common cases. → *Acceptable trade-off*: Keeps implementation simple.
