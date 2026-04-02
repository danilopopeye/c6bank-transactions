## 1. Parser Package API

- [x] 1.1 Add exported `ParseFile(path string) ([]Transaction, error)` function to `internal/parser/` that opens a file, detects format by extension, calls the appropriate scanner, and returns `[]Transaction`
- [x] 1.2 Add exported `Deduplicate(transactions []Transaction) []Transaction` function that removes duplicates using FNV-1a hash on Date+Payee+Amount+Memo
- [ ] 1.3 Test `ParseFile` with each supported format (PDF, CSV, PNG) and error cases
- [ ] 1.4 Test `Deduplicate` with overlapping and non-overlapping transactions

## 2. CLI Entry Point

- [ ] 2.1 Create `cmd/cli/main.go` with argument parsing using `flag` package (support `-o` flag for output file)
- [ ] 2.2 Implement main logic: iterate over file args, call `ParseFile` for each (fail-fast on any error), accumulate results, deduplicate, sort chronologically, write CSV
- [ ] 2.3 Handle edge cases: no args (usage + exit 1), missing files, unsupported formats, invalid `-o` path, password-protected PDFs
- [ ] 2.4 Test CLI main with no args, single file, and multiple files
