# Project Context

## Purpose
C6 Bank Transactions is a Go web service that processes C6 Bank transaction statements and converts them to QIF/CSV formats for personal finance software. The application supports three input types: PDF statements, CSV files, and mobile screenshots (which require OCR processing).

## Tech Stack
- Go 1.21+
- Web Framework: Standard library `net/http`
- PDF Processing: `github.com/ledongthuc/pdf`
- OCR: Tesseract OCR with Portuguese+English language packs
- Hash Generation: `github.com/segmentio/fasthash/fnv1a`
- Container: Docker with multi-stage builds
- Testing: Go's built-in testing framework with race condition detection

## Project Conventions

### Code Style
- Standard Go formatting (`go fmt ./...`)
- Go vet for static analysis (`go vet ./...`)
- Interface-based design for parser extensibility
- Error handling with explicit returns and error wrapping
- Package organization by domain (`internal/` structure)

### Architecture Patterns
- **Interface-based Parser Architecture**: Extensible parser system supporting PDF, CSV, and image formats
- **Smart Image Processing**: iPhone model detection for precise transaction area cropping
- **OCR Integration Pipeline**: Image → Crop → OCR → Parse → Format
- **Graceful Shutdown**: HTTP server with signal handling
- **Multi-stage Docker Build**: Optimized container with Tesseract OCR

### Testing Strategy
- Unit tests with comprehensive coverage
- Integration tests via HTTP upload endpoint
- Race condition testing (`go test -race`)
- Fixture-based testing with sample images and data
- Coverage reporting (`go tool cover --func coverage.txt`)

### Git Workflow
- Main branch: `master`
- Feature branches for new functionality
- Pull requests for code review
- Conventional commit messages
- Automated builds with GitHub Actions

## Domain Context

### Brazilian Banking System
- C6 Bank is a Brazilian digital bank
- Statements follow specific format patterns for dates and amounts
- Portuguese language support required for OCR
- Currency format: R$ (Brazilian Real)
- Installment payments tracking (parcelas)

### Transaction Processing
- Regex patterns optimized for C6 Bank statement format
- Transaction categories and descriptions
- Unique ID generation using FNV-1a hash
- QIF format compatibility with personal finance software
- CSV export for spreadsheet applications

### Mobile Processing
- iPhone model detection for screenshot processing
- Screen resolution-based cropping logic
- OCR with Portuguese and English language models

## Important Constraints

### Technical Constraints
- OCR requires Tesseract system dependency
- Mobile screenshots need precise iPhone model detection
- PDF password handling must be secure and robust
- HTTP server must handle file uploads safely

### Business Constraints
- Brazilian Portuguese language support mandatory
- Must handle various iPhone screen resolutions
- Real-time processing within web request timeout
- Fallback for unparseable content

### Security Constraints
- File upload validation and sanitization
- PDF password handling in memory only
- Input validation for all file formats
- No persistent storage of user data

## External Dependencies

### System Dependencies
- Tesseract OCR engine with Portuguese and English language packs
- Docker container runtime for deployment

### Go Libraries
- `github.com/ledongthuc/pdf` - PDF text extraction
- `github.com/segmentio/fasthash/fnv1a` - Fast hash generation
- Standard library packages for HTTP server and file handling

### Integration Points
- C6 Bank statement formats (PDF, CSV, mobile app screenshots)
- Personal finance software via QIF export format
- Spreadsheet applications via CSV export format
