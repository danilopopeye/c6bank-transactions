# Mobile - iPhone Model Detection

This package defines iPhone models and their screen dimensions for smart image cropping in transaction processing.

## Phone Model Structure

```go
type Phone struct {
    Width, Height  int  // Image dimensions in pixels
    Header, Footer int  // Top/bottom margins to exclude from transaction area
    Month          int  // Y-position where month region starts (0 = top)
    MonthSize      int  // Height of month region to extract
}
```

## Supported Models

| Model | Dimensions | Header | Footer | Month Y | MonthSize |
|-------|------------|--------|--------|---------|-----------|
| IPhone13 | 1170×2532 | 755px | 245px | 640px | 150px |
| IPhone13ProMax | 1284×2778 | 800px | 250px | 0 (disabled) | 150px |
| IPhone15Pro | 1179×2556 | 776px | 250px | 660px | 150px |
| IPhone16Pro | 1206×2622 | 800px | 250px | 660px | 150px |
| IPhoneMirror | 836×1840 | 600px | 180px | 0px | 100px |

## IPhoneMirror

iPhone Mirror is a special format from macOS screen mirroring with these characteristics:

- **Smaller dimensions**: 836×1840 (vs physical iPhone screens)
- **Transparent header**: First row has alpha=0 pixels (iPhone Mirror signature)
- **Custom regions**: Header=600px, Footer=180px, MonthSize=100px

Detection requires BOTH:
1. Exact dimensions: 836×1840 pixels
2. Transparency in first 10 pixels of first row

## Usage

Images are automatically detected by dimensions in `internal/image.GetPhone()`:

```go
phone, err := GetPhone(img)  // Returns Phone model or ErrUnsupportedPhone
cropped := CropImage(img, phone)    // Transaction area
month := CropMonth(img, phone)      // Month reference
```

## Adding New Models

To add a new iPhone model:

1. Measure dimensions and crop regions from sample screenshots
2. Add const: `ModelName = Phone{Width, Height, Header, Footer, MonthY, MonthSize}`
3. Append to `Phones` array
4. Add inline comment explaining each value
