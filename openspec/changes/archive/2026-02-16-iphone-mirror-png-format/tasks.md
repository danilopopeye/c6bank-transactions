## 1. Phone Model Definition

- [x] 1.1 Add `MonthSize int` field to `mobile.Phone` struct
- [x] 1.2 Add `IPhoneMirror = Phone{836, 1840, 600, 180, 0, 100}` to `mobile/phones.go`
- [x] 1.3 Add `IPhoneMirror` to END of `mobile.Phones` array (fallback position)
- [x] 1.4 Update existing phone models to include `MonthSize: 150` explicitly
- [x] 1.5 Add inline comments for IPhoneMirror following existing patterns
- [x] 1.6 Create `internal/mobile/README.md` explaining models and MonthSize

## 2. Transparency Validation

- [x] 2.1 Implement `func HasTransparency(img image.Image) bool` in `internal/image/image.go`
- [x] 2.2 Check alpha channel of first 10 pixels in first row (x=0 to x=9, y=0)
- [x] 2.3 Return true only if ALL 10 pixels have alpha == 0 (strict zero)
- [x] 2.4 Add unit test for fully transparent image (returns true)
- [x] 2.5 Add unit test for partially transparent image (returns false)
- [x] 2.6 Add unit test for opaque image (returns false)
- [x] 2.7 Create synthetic PNG test fixture using Go stdlib (image.NewRGBA, png.Encode)

## 3. Update Image Processing

- [x] 3.1 Modify `GetPhone()` to call `HasTransparency()` at the start
- [x] 3.2 Add early return for IPhoneMirror (transparency + dimensions 836×1840)
- [x] 3.3 Update `CropMonth()` signature to use `phone.MonthSize` instead of global constant
- [x] 3.4 Add fallback logic: if `phone.MonthSize == 0`, use 150px
- [x] 3.5 Wrap `ErrUnsupportedPhone` with custom message for transparency validation failure
- [x] 3.6 Ensure error messages follow existing format patterns

## 4. Tests

- [x] 4.1 Unit test: `HasTransparency()` with alpha == 0 (all transparent)
- [x] 4.2 Unit test: `HasTransparency()` with alpha > 0 (opaque)
- [x] 4.3 Unit test: `HasTransparency()` with mixed alpha values (1-254)
- [x] 4.4 Unit test: `GetPhone()` detects IPhoneMirror by transparency + dimensions
- [x] 4.5 Unit test: `GetPhone()` rejects 836×1840 without transparency
- [x] 4.6 Unit test: `GetPhone()` rejects transparent image with wrong dimensions
- [x] 4.7 Unit test: `CropMonth()` produces 100px for IPhoneMirror
- [x] 4.8 Unit test: `CropMonth()` produces 150px for other models
- [x] 4.9 Integration test: Full flow PNG → HasTransparency → GetPhone → Crop → validate outputs

## 5. Documentation

- [x] 5.1 Add godoc comment to `HasTransparency()` function
- [x] 5.2 Add godoc comment to `MonthSize` field in Phone struct
- [x] 5.3 Write package README at `internal/mobile/README.md`
- [x] 5.4 Update inline comments in `phones.go` for IPhoneMirror definition
