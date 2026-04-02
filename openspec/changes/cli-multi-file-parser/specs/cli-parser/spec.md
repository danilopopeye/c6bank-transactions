## ADDED Requirements

### Requirement: CLI accepts multiple file paths as arguments

The CLI SHALL accept one or more file paths as positional arguments. Each file SHALL be parsed using the appropriate format handler based on file extension (.pdf, .csv, .jpg, .png).

#### Scenario: Single file argument
- **WHEN** user runs `cli transactions.pdf`
- **THEN** the CLI parses the PDF and outputs CSV with those transactions

#### Scenario: Multiple file arguments
- **WHEN** user runs `cli statement.pdf Fatura_2026-01-15.csv screenshot.png`
- **THEN** the CLI parses all three files and outputs a single CSV with accumulated transactions

#### Scenario: No file arguments
- **WHEN** user runs `cli` without any arguments
- **THEN** the CLI prints usage help and exits with code 1

#### Scenario: Unsupported file extension
- **WHEN** user runs `cli data.xlsx`
- **THEN** the CLI prints an error message indicating unsupported format and exits with code 1

### Requirement: CLI outputs CSV to stdout by default

The CLI SHALL write CSV output to stdout. The CSV SHALL have a header row with columns: Date, Payee, Memo, Value.

#### Scenario: Output to stdout
- **WHEN** user runs `cli transactions.pdf`
- **THEN** CSV content is written to stdout with header row and transaction data

#### Scenario: Pipe to file
- **WHEN** user runs `cli transactions.pdf > output.csv`
- **THEN** CSV content is written to the redirected file

### Requirement: CLI supports output file via flag

The CLI SHALL accept an `-o <path>` flag that writes CSV output to the specified file instead of stdout.

#### Scenario: Output to file
- **WHEN** user runs `cli -o result.csv statement.pdf Fatura_2026-01-15.csv`
- **THEN** CSV content is written to `result.csv`

#### Scenario: Output file path is directory
- **WHEN** user runs `cli -o /tmp/ statement.pdf`
- **THEN** the CLI prints an error and exits with code 1

### Requirement: Transactions are deduplicated across files

The CLI SHALL deduplicate transactions across all input files. Two transactions are considered duplicates if they have the same Date, Payee, and Amount.

#### Scenario: Duplicate across two files
- **WHEN** two files contain the same transaction (same date, payee, amount)
- **THEN** only one instance appears in the output CSV

#### Scenario: No duplicates
- **WHEN** no transactions share the same date, payee, and amount
- **THEN** all transactions appear in the output CSV

### Requirement: CLI reuses existing parser logic

The CLI SHALL use the existing format-specific scanners from `internal/parser` for PDF, CSV, and image parsing. No new parsing logic SHALL be introduced.

#### Scenario: PDF parsing
- **WHEN** user passes a .pdf file
- **THEN** the CLI uses the existing PDF scanner to extract transactions

#### Scenario: CSV parsing
- **WHEN** user passes a .csv file matching the `Fatura_YYYY-MM-DD.csv` naming convention
- **THEN** the CLI uses the existing CSV scanner to extract transactions

#### Scenario: CSV with wrong filename
- **WHEN** user passes a .csv file not matching the naming convention
- **THEN** the CLI reports a filename validation error

#### Scenario: Image parsing
- **WHEN** user passes a .png or .jpg file
- **THEN** the CLI uses the existing image scanner to extract transactions
