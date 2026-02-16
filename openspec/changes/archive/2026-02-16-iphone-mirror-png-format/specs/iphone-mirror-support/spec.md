## ADDED Requirements

### Requirement: Detect iPhone Mirror by transparency and dimensions
The system SHALL detect iPhone Mirror format when BOTH transparency validation passes AND dimensions match 836×1840 pixels.

#### Scenario: Valid iPhone Mirror PNG
- **WHEN** a PNG image with 836×1840 dimensions AND transparent pixels (alpha == 0) in first 10 pixels of first row is provided
- **THEN** the system SHALL identify it as IPhoneMirror model
- **AND** processing continues with Mirror-specific settings (Header=600, Footer=180, MonthSize=100)

#### Scenario: Wrong dimensions (Mirror-like but wrong size)
- **WHEN** a PNG image has transparency but dimensions are NOT 836×1840
- **THEN** the system SHALL NOT identify it as IPhoneMirror
- **AND** SHALL return `ErrUnsupportedPhone` with standard error message

#### Scenario: Missing transparency (Mirror dimensions but opaque)
- **WHEN** a PNG image has 836×1840 dimensions but NO transparency in first row
- **THEN** the system SHALL return `ErrUnsupportedPhone` wrapped with custom message
- **AND** error message shall follow existing format patterns

### Requirement: Transparency validation uses strict zero threshold
The system SHALL validate transparency by checking if alpha channel equals zero in first 10 pixels.

#### Scenario: All pixels transparent (alpha == 0)
- **WHEN** first 10 pixels of first row have alpha == 0
- **THEN** `HasTransparency()` SHALL return true

#### Scenario: Some pixels opaque (alpha > 0)
- **WHEN** ANY of first 10 pixels has alpha > 0
- **THEN** `HasTransparency()` SHALL return false

#### Scenario: Partial transparency (alpha between 1-254)
- **WHEN** pixels have alpha values between 1-254 (not zero, not fully opaque)
- **THEN** `HasTransparency()` SHALL return false (strict zero check)
- **AND** threshold may be adjusted empirically after testing with real screenshots

### Requirement: IPhoneMirror model definition
The system SHALL define IPhoneMirror with specific dimensions and region settings.

#### Scenario: IPhoneMirror structure
- **WHEN** IPhoneMirror is defined in `mobile/phones.go`
- **THEN** it SHALL have `Width = 836, Height = 1840, Header = 600, Footer = 180, Month = 500, MonthSize = 100`
- **AND** SHALL be added at the END of `mobile.Phones` array (fallback position)

#### Scenario: MonthSize for IPhoneMirror
- **WHEN** IPhoneMirror is used for cropping
- **THEN** month region SHALL be 100 pixels tall
- **AND** SHALL start at Y position 500

### Requirement: Phone struct includes MonthSize field
The `mobile.Phone` struct SHALL include `MonthSize int` field to support variable month region heights.

#### Scenario: Existing models without MonthSize
- **WHEN** existing Phone models are instantiated without MonthSize field
- **THEN** MonthSize SHALL default to 150px
- **AND** CropMonth SHALL use phone.MonthSize if > 0, else 150

#### Scenario: Explicit MonthSize definition
- **WHEN** Phone model defines MonthSize explicitly (like IPhoneMirror with 100)
- **THEN** CropMonth SHALL use the explicit MonthSize value
- **AND** SHALL NOT use the global constant

#### Scenario: All models updated
- **WHEN** all existing iPhone models are updated
- **THEN** each SHALL have `MonthSize: 150` explicitly defined
- **AND** no model shall rely on default behavior

### Requirement: CropMonth uses phone-specific MonthSize
The `CropMonth()` function SHALL use `phone.MonthSize` instead of global `mobile.MonthSize` constant.

#### Scenario: Backward compatibility
- **WHEN** CropMonth is called with phone where MonthSize == 0
- **THEN** it SHALL fallback to 150px (existing behavior)
- **AND** SHALL NOT break existing functionality

#### Scenario: IPhoneMirror produces 100px month region
- **WHEN** CropMonth is called with IPhoneMirror
- **THEN** it SHALL produce a 100px tall month region
- **AND** other models produce 150px month region

### Requirement: GetPhone validates transparency
The `GetPhone()` function SHALL call `HasTransparency()` at the start to detect iPhone Mirror format.

#### Scenario: IPhoneMirror detection flow
- **WHEN** GetPhone is called with a PNG image
- **THEN** it SHALL first call HasTransparency()
- **AND** if transparency exists AND dimensions match 836×1840, return IPhoneMirror
- **AND** if no transparency, continue with dimension-based detection

#### Scenario: Error handling
- **WHEN** HasTransparency returns true but dimensions don't match any model
- **THEN** GetPhone SHALL return `ErrUnsupportedPhone` with standard message
- **AND** SHALL NOT mention transparency specifically (implementation detail)
