
# C6 Bank Transactions - Agent Guide

This document provides essential information for AI agents working with this Go web service that processes C6 Bank transaction statements and converts them to QIF/CSV formats.

## Project Overview

A Go web service that processes C6 Bank transaction statements from multiple input formats (PDF, CSV, mobile screenshots) and converts them to QIF/CSV formats for personal finance software. The application features OCR processing for mobile screenshots with iPhone model detection.

## Essential Commands

### Development
```bash
# Build the application
go build -o /dev/null ./...

# Run the server (listens on port 4500)
go run ./cmd/c6bank-transactions

# Run full test suite with coverage
go test -v -race -coverprofile=coverage.txt ./...
go tool cover --func coverage.txt

# Code quality checks
go fmt ./...
go vet ./...
go mod tidy

# Test upload endpoint with sample image
./run.sh [optional-fixture-file]

# Direct curl test
curl -X 'POST' -F "file=@test/fixtures/IMG_0352.PNG;type=image/png" \
  'http://localhost:4500/upload'
```

### Docker
```bash
# Build container image
docker build -t git.home/danilo/c6bank-transactions .

# Run container
docker run -p 4500:4500 git.home/danilo/c6bank-transactions
```

## Code Organization

### Directory Structure
```
cmd/c6bank-transactions/    # HTTP server and web interface
├── main.go                # Server setup with graceful shutdown
├── handler.go             # File upload and processing handlers
└── index.html             # Web interface

internal/                  # Application core packages
├── parser/               # Multi-format transaction parsing
│   ├── parser.go         # Main orchestration and format detection
│   ├── pdf.go            # PDF statement processing
│   ├── csv.go            # CSV statement processing
│   ├── image.go          # Image processing with OCR
│   ├── transaction.go    # Transaction data structures
│   ├── time.go           # Date parsing utilities
│   ├── validation.go     # Input validation
│   └── ocr/              # OCR processing
│       └── ocr.go        # Tesseract integration
├── image/                # Image processing utilities
│   ├── image.go          # Smart cropping based on iPhone models
│   └── image_test.go     # Image processing tests
├── mobile/               # iPhone model specifications
│   └── phones.go         # Device dimensions for cropping
└── qif/                  # QIF format generation
    ├── qif.go            # QIF template and transaction formatting
    └── qif_test.go       # QIF generation tests

test/fixtures/            # Sample images and transaction data
```

## Key Components

### Parser Engine (`internal/parser/`)
- **Multi-format Support**: PDF, CSV, and image processing
- **Format Detection**: Automatic file type identification based on extension
- **CSV Validation**: Strict filename validation (`Fatura_YYYY-MM-DD.csv`)
- **OCR Integration**: Portuguese+English language packs for Brazilian statements

### Image Processing (`internal/image/`)
- **Smart Cropping**: Auto-detects iPhone models from image dimensions
- **iPhone Models**: iPhone 13, 13 Pro Max, 15 Pro, 16 Pro supported
- **OCR Preparation**: Crops transaction areas for optimal OCR results

### HTTP Server (`cmd/c6bank-transactions/`)
- **Graceful Shutdown**: Proper signal handling with context cancellation
- **File Upload**: Multipart form handling with size limits
- **Error Handling**: Structured error responses
- **Timeouts**: Configurable read/write/idle timeouts

## Code Conventions

### Naming Patterns
- **Packages**: Short, lowercase (parser, image, qif, mobile)
- **Functions**: CamelCase with exported names starting with capital
- **Variables**: camelCase for local, uppercase for constants
- **Files**: Snake case (parser.go, image_test.go)

### Error Handling
- Use `fmt.Errorf` with `fmt.Errorf("%w: details", err)` for wrapping
- Define custom error variables (e.g., `ErrWrongCSVFilename`)
- Return errors as the last return value
- Use structured error types when additional context is needed

### Testing Patterns
- Test files: `*_test.go` in same package
- Use `testify/assert` and `testify/require` for assertions
- Run tests in parallel with `t.Parallel()`
- Use table-driven tests for multiple scenarios
- Test files in `test/fixtures/` for sample data

## Important Implementation Details

### OCR Configuration
- **Languages**: Portuguese (`por`) + English (`eng`)
- **Page Segmentation**: PSM 4 for single column text
- **Binary**: System `tesseract` command
- **Dependency**: Must be available in Docker image

### Transaction ID Generation
- Uses FNV-1a hash (`github.com/segmentio/fasthash/fnv1a`)
- Hashes `Date + Payee + Amount` for unique QIF transaction IDs

### iPhone Model Support
```go
// Dimensions: Width, Height, Header, Footer, Month
IPhone13       = Phone{1170, 2532, 755, 245, 640}
IPhone13ProMax = Phone{1284, 2778, 800, 250, 0}
IPhone15Pro    = Phone{1179, 2556, 776, 250, 660}
IPhone16Pro    = Phone{1206, 2622, 800, 250, 660}
```

### PDF Processing
- Supports password-protected PDFs
- Handles both scanned and native PDF text
- Uses `github.com/ledongthuc/pdf` for text extraction

### CSV Processing
- Strict filename validation: `Fatura_YYYY-MM-DD.csv`
- Date parsing from filename for transaction reference
- Credit card transaction format

## Configuration

### Environment Variables
- `HOST`: Server bind address (default: "0.0.0.0")
- `PORT`: Server port (default: "4500")

### Constants
- `MAX_UPLOAD_SIZE`: 10MB file upload limit
- Server timeouts: 5s read, 1s read header, 10s write, 1m idle

## Testing Strategy

### Test Categories
- **Unit Tests**: Comprehensive coverage for all parser components
- **Integration Tests**: End-to-end upload → processing → download flow
- **Fixtures**: Sample images in `test/fixtures/`
- **Race Condition Testing**: All tests run with `-race` flag

### Running Tests
```bash
# Run all tests with coverage
go test -v -race -coverprofile=coverage.txt ./...

# Run specific package tests
go test -v ./internal/parser/...

# Run coverage analysis
go tool cover --func coverage.txt
```

## CI/CD Configuration

### GitHub Actions (`.github/workflows/`)
- **Go Workflow**: Format, tidy, vet, build, test on push/PR
- **GolangCI Lint**: Code quality checks with latest golangci-lint
- **OCI Image**: Docker build and push on master branch

### Forgejo Actions (`.forgejo/workflows/`)
- Similar CI setup for self-hosted Forgejo
- Docker image building and registry pushing

## Gotchas & Non-Obvious Patterns

### File Processing
- Images use smart cropping based on iPhone detection
- CSV files must follow exact naming convention
- PDFs support password protection
- All files processed in memory (no streaming currently)

### Error Handling Nuances
- Mixed use of `fmt.Printf` vs `log.Printf` (technical debt)
- Some errors expose internal file paths
- Inconsistent upload size limits between main.go and handler.go

### Testing Quirks
- One test skipped due to refactoring needs (`internal/image/image_test.go:18`)
- HTTP handlers have 0% test coverage (critical gap)
- OCR tests depend on system Tesseract installation

### Performance Considerations
- Entire files loaded into memory (not streamed)
- OCR processing is CPU-intensive
- No connection pooling for HTTP clients
- No rate limiting on upload endpoint

## Development Workflow

1. **Setup**: Ensure Go 1.22+ and Tesseract OCR installed
2. **Local Development**: Run `go run ./cmd/c6bank-transactions`
3. **Testing**: Use `go test -v -race ./...` before committing
4. **Build**: Verify with `go build -o /dev/null ./...`
5. **Docker**: Test with `docker build -t test . && docker run -p 4500:4500 test`

## Security Considerations

### Current State
- No rate limiting on upload endpoint
- No input validation beyond file type/size
- Potential information disclosure in error messages
- No malware scanning for uploaded files

### Recommendations (from TODO.md)
- Add rate limiting and CORS configuration
- Implement structured logging
- Add health check endpoint
- Fix inconsistent upload size limits
- Add security headers

## Common Issues

### OCR Failures
- Ensure Tesseract is installed with Portuguese and English language packs
- Check image dimensions match supported iPhone models
- Verify image format (PNG/JPG)

### CSV Processing
- Filename must be exactly `Fatura_YYYY-MM-DD.csv`
- Date format in filename must be valid
- Reference date extracted from filename

### PDF Processing
- Password-protected PDFs supported
- Both scanned and native PDFs handled
- Text extraction depends on PDF quality

## Dependencies

### Key Libraries
- `github.com/ledongthuc/pdf`: PDF text extraction
- `github.com/segmentio/fasthash/fnv1a`: Fast hashing for transaction IDs
- `github.com/stretchr/testify`: Testing framework

### System Dependencies
- `tesseract`: OCR processing with Portuguese+English language packs

## Architecture Decisions

### Interface-Based Design
- Parser architecture uses interfaces for extensibility
- Easy to add new input formats
- Clean separation of concerns

### Smart Image Processing
- iPhone model detection enables precise cropping
- Optimizes OCR accuracy by focusing on transaction areas
- Reduces processing time and improves results

### Multi-Language OCR
- Portuguese+English language packs for Brazilian bank statements
- Handles bilingual content effectively
- PSM 4 mode optimized for single column transaction lists
