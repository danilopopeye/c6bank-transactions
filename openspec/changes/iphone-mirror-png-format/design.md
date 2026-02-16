## Context

iPhone Mirror on macOS creates PNG screenshots when users capture their iPhone screen. These images have a specific format with transparent regions at the top (notification bar) and need to be processed to extract transaction data from banking apps.

Current state: The system already processes images via `internal/parser/image.go` with:
- Phone model detection by dimensions (`internal/mobile/phones.go`)
- Image cropping (Header/Footer for transactions, Month region for reference)
- OCR processing on cropped images

Existing iPhone models: IPhone13, IPhone13ProMax, IPhone15Pro, IPhone16Pro

Constraints:
- Images are always 836×1840 pixels (smaller than physical iPhone screens)
- Header region: 600px (larger than existing models)
- Footer region: 180px (smaller than existing models)
- Month region: 100px height (smaller than existing 150px)
- Must validate transparency in first row (iPhone Mirror characteristic)

## Goals / Non-Goals

**Goals:**
- Add IPhoneMirror model to existing `mobile.Phones` array
- Add transparency validation before phone detection
- Support MonthSize=100px for iPhone Mirror (vs 150px for others)
- Reuse existing `image.Crop()` and `parser.Parse()` pipeline

**Non-Goals:**
- Reimplement OCR or transaction parsing (already exists)
- Create new package structure (use `internal/image/` and `internal/mobile/`)
- Support other image dimensions or formats

## Decisions

**Reuse existing `mobile.Phone` struct**
- Rationale: Already has Width, Height, Header, Footer, Month fields
- Add `MonthSize int` field to support variable month region heights
- Add `IPhoneMirror = Phone{836, 1840, 600, 180, 0, 100}` to `phones.go` (at end of array)

**Add transparency validation as dedicated function**
- Function: `func HasTransparency(img image.Image) bool`
- Check first 10 pixels of first row for alpha == 0 (strict zero, adjust empirically after testing)
- Called at start of `GetPhone()` for early detection
- Wrap `ErrUnsupportedPhone` with custom message on failure

**Both transparency AND dimensions required**
- Rationale: iPhone Mirror always has transparency AND 836×1840 dimensions
- Both conditions must be true to confirm iPhone Mirror format
- Prevents false positives from edited images

**MonthSize with default fallback**
- Add `MonthSize int` field to `Phone` struct
- Default to 150px if MonthSize == 0 (backward compatibility)
- IPhoneMirror uses 100px, all other models use 150px
- Update `CropMonth()` to use `phone.MonthSize` instead of global constant

**Error handling consistency**
- Reuse existing error message patterns
- No new error types - wrap `ErrUnsupportedPhone`
- Same format as other format validations

**Test fixtures**
- Generate synthetic PNG using Go stdlib (`image.NewRGBA`, `png.Encode`)
- Set alpha = 0 in first 10 pixels for transparency test
- Unit tests for each function (HasTransparency, GetPhone, CropMonth)

**Code organization**
- Add `IPhoneMirror` to end of `mobile.Phones` array (fallback position)
- Inline comments following existing patterns in `phones.go`
- Package README at `internal/mobile/README.md` explaining models and MonthSize

## Risks / Trade-offs

**[Risk]** MonthSize currently is global constant (150px)
- **Mitigation**: Add `MonthSize int` field to `Phone` struct, use default 150 if zero
- Update existing phone models to include `MonthSize: 150` explicitly

**[Risk]** Alpha == 0 may be too strict (anti-aliasing or compression artifacts)
- **Mitigation**: Test with real iPhone Mirror screenshots, adjust threshold empirically
- Consider alpha <= 10 buffer if needed after testing

**[Risk]** Decode cost for transparency check (scanline partial)
- **Mitigation**: `HasTransparency()` receives `image.Image` (already decoded)
- GetPhone already decodes image, no additional cost for transparency check

**[Trade-off]** Strict validation (both transparency AND dimensions)
- Benefit: Prevents false positives, clear format detection
- Cost: May reject valid images in edge cases (mitigated by empirical testing)

**[Trade-off]** IPhoneMirror at end of Phones array (fallback position)
- Benefit: Doesn't interfere with existing model detection
- Cost: May be slightly slower (iterates through all models first)
- Rationale: 836×1840 is unique dimension, collision unlikely

## Implementation Notes

**Thread safety**: All functions are stateless and pure, no shared state

**Performance**: Transparency check adds minimal overhead (first 10 pixels only)

**Error messages**: Reuse existing patterns - maintain consistency with PDF/CSV validation

**Documentation**: Package README at `internal/mobile/README.md` explaining:
- Phone model detection by dimensions
- MonthSize field and its purpose
- iPhone Mirror specific behavior (transparency check)

## Open Questions

All resolved - architecture decisions made and documented
