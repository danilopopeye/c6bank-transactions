## Why

The current application is an HTTP server that processes C6 Bank transaction files one at a time via a web upload endpoint. Users who need to process multiple exported files (e.g., several months of PDF statements or CSV exports) must upload them individually through the browser. A CLI tool would allow batch processing of multiple files from the command line, producing a single consolidated output — a common workflow for personal finance management.

## What Changes

- Add a new CLI entry point (`cmd/cli`) that accepts one or more file paths as arguments
- Reuse the existing parser package (`internal/parser`) to parse each file into transactions
- Deduplicate transactions across files (same transaction may appear in multiple files)
- Output a single CSV file with all accumulated transactions to stdout or a specified output file
- Support the same input formats already handled by the parser: PDF, CSV, and mobile screenshots (PNG/JPG)

## Capabilities

### New Capabilities
- `cli-parser`: Command-line interface for parsing multiple C6 Bank transaction files into a single CSV output, including file argument handling, transaction accumulation, deduplication, and CSV generation

### Modified Capabilities

_(none — this is a new entry point that reuses existing parser internals)_

## Impact

- **New code**: `cmd/cli/` directory with CLI entry point
- **Dependencies**: May need a CLI argument parsing library or use `flag` stdlib
- **Existing packages**: `internal/parser` consumed as a library (no changes expected)
- **Build**: New binary target in addition to the existing web server
