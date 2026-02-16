## 1. Phone Model Definition

- [ ] 1.1 Add `MonthSize int` field to `mobile.Phone` struct
- [ ] 1.2 Add `IPhoneMirror = Phone{836, 1840, 600, 180, 0, 100}` to `mobile/phones.go`
- [ ] 1.3 Add `IPhoneMirror` to END of `mobile.Phones` array (fallback position)
- [ ] 1.4 Update existing phone models to include `MonthSize: 150` explicitly
- [ ] 1.5 Add inline comments for IPhoneMirror following existing patterns
- [ ] 1.6 Create `internal/mobile/README.md` explaining models and MonthSize

## 2. Transparency Validation

- [ ] 2.1 Implement `func HasTransparency(img image.Image) bool` in `internal/image/image.go`
- [ ] 2.2 Check alpha channel of first 10 pixels in first row (x=0 to x=9, y=0)
- [ ] 2.3 Return true only if ALL 10 pixels have alpha == 0 (strict zero)
- [ ] 2.4 Add unit test for fully transparent image (returns true)
- [ ] 2.5 Add unit test for partially transparent image (returns false)
- [ ] 2.6 Add unit test for opaque image (returns false)
- [ ] 2.7 Create synthetic PNG test fixture using Go stdlib (image.NewRGBA, png.Encode)

## 3. Update Image Processing

- [ ] 3.1 Modify `GetPhone()` to call `HasTransparency()` at the start
- [ ] 3.2 Add early return for IPhoneMirror (transparency + dimensions 836×1840)
- [ ] 3.3 Update `CropMonth()` signature to use `phone.MonthSize` instead of global constant
- [ ] 3.4 Add fallback logic: if `phone.MonthSize == 0`, use 150px
- [ ] 3.5 Wrap `ErrUnsupportedPhone` with custom message for transparency validation failure
- [ ] 3.6 Ensure error messages follow existing format patterns

## 4. Tests

- [ ] 4.1 Unit test: `HasTransparency()` with alpha == 0 (all transparent)
- [ ] 4.2 Unit test: `HasTransparency()` with alpha > 0 (opaque)
- [ ] 4.3 Unit test: `HasTransparency()` with mixed alpha values (1-254)
- [ ] 4.4 Unit test: `GetPhone()` detects IPhoneMirror by transparency + dimensions
- [ ] 4.5 Unit test: `GetPhone()` rejects 836×1840 without transparency
- [ ] 4.6 Unit test: `GetPhone()` rejects transparent image with wrong dimensions
- [ ] 4.7 Unit test: `CropMonth()` produces 100px for IPhoneMirror
- [ ] 4.8 Unit test: `CropMonth()` produces 150px for other models
- [ ] 4.9 Integration test: Full flow PNG → HasTransparency → GetPhone → Crop → validate outputs

## 5. Documentation

- [ ] 5.1 Add godoc comment to `HasTransparency()` function
- [ ] 5.2 Add godoc comment to `MonthSize` field in Phone struct
- [ ] 5.3 Write package README at `internal/mobile/README.md`
- [ ] 5.4 Update inline comments in `phones.go` for IPhoneMirror definition
