## Why

iPhone Mirror on macOS generates PNG screenshots with smaller dimensions (836×1840) and transparent regions at the top. The existing image processing pipeline supports multiple iPhone models but doesn't recognize this specific format.

## What Changes

- Add IPhoneMirror model to `internal/mobile/phones.go` with dimensions 836×1840
- Add transparency validation in `internal/image/image.go` to detect iPhone Mirror format
- Update `mobile.Phone` struct to include `MonthSize` field (100px for Mirror vs 150px for others)
- No breaking changes - extends existing PNG/JPG image processing

## Capabilities

### New Capabilities

- `iphone-mirror-support`: Add iPhone Mirror as a recognized phone model with transparency validation and custom month region size

### Modified Capabilities

- `phone-detection`: Extend to support transparency check and variable MonthSize per model

## Impact

- **Modified code**: `internal/mobile/phones.go` (add IPhoneMirror), `internal/image/image.go` (transparency check), `mobile.Phone` struct (add MonthSize)
- **No new dependencies**: Uses existing Go `image` package
- **Affected systems**: Image processing pipeline (automatic - no API changes)
