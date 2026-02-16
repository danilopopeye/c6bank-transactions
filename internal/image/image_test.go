package image_test

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"testing"

	subject "git.home/c6bank-transactions/internal/image"
	"git.home/c6bank-transactions/internal/mobile"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCrop(t *testing.T) {
	t.Parallel()

	t.Skip("need refactor")

	tests := []struct {
		name   string
		width  int
		height int
	}{
		{"iPhone 16 Pro", 1206, 2622},
		{"iPhone 13 Pro Max", 1284, 2778},
		{"iPhone 13", 1170, 2532},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			buf := bytes.NewBuffer(nil)
			img := image.NewRGBA(image.Rect(0, 0, test.width, test.height))

			err := png.Encode(buf, img)
			require.NoError(t, err)

			bufSeeker := bytes.NewReader(buf.Bytes())

			croppedImg, _, err := subject.Crop(bufSeeker)
			require.NoError(t, err)

			croppedPNG, err := png.Decode(croppedImg)
			require.NoError(t, err)

			height := test.height

			assert.Equal(t, 0, croppedPNG.Bounds().Max.X)
			assert.Equal(t, height, croppedPNG.Bounds().Max.Y)
		})
	}
}

func TestCropImage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		phone mobile.Phone
	}{
		{"iPhone 16 Pro", mobile.IPhone16Pro},
		{"iPhone 13 Pro Max", mobile.IPhone13ProMax},
		{"iPhone 13", mobile.IPhone13},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testImage := buildImage(t, test.phone)
			cropped := subject.CropImage(testImage, test.phone)
			bounds := cropped.Bounds().Max

			assert.Equal(t, test.phone.Width, bounds.X)
			assert.Equal(t, test.phone.Height-test.phone.Header-test.phone.Footer, bounds.Y)
		})
	}
}

func TestCropMonth(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		phone mobile.Phone
	}{
		{"iPhone 16 Pro", mobile.IPhone16Pro},
		{"iPhone 13 Pro Max", mobile.IPhone13ProMax},
		{"iPhone 13", mobile.IPhone13},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testImage := buildImage(t, test.phone)
			cropped := subject.CropMonth(testImage, test.phone)
			bounds := cropped.Bounds().Max

			assert.Equal(t, test.phone.Width, bounds.X)
			assert.Equal(t, mobile.MonthSize, bounds.Y)
		})
	}
}

func TestGetPhone(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		phone mobile.Phone
		err   error
	}{
		{"iPhone 16 Pro", mobile.IPhone16Pro, nil},
		{"iPhone 13 Pro Max", mobile.IPhone13ProMax, nil},
		{"iPhone 13", mobile.IPhone13, nil},
		{"other", mobile.Phone{}, subject.ErrUnsupportedPhone},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testImage := buildImage(t, test.phone)
			phone, err := subject.GetPhone(testImage)

			assert.ErrorIs(t, err, test.err)
			assert.Equal(t, test.phone, phone)
		})
	}
}

func buildImage(t *testing.T, phone mobile.Phone) image.Image {
	t.Helper()

	return image.Rect(0, 0, phone.Width, phone.Height)
}

func TestHasTransparency(t *testing.T) {
	t.Parallel()

	t.Run("all transparent pixels (alpha == 0)", func(t *testing.T) {
		t.Parallel()

		img := createTransparentImage(t, 100, 100, true)
		assert.True(t, subject.HasTransparency(img))
	})

	t.Run("all opaque pixels (alpha > 0)", func(t *testing.T) {
		t.Parallel()

		img := createTransparentImage(t, 100, 100, false)
		assert.False(t, subject.HasTransparency(img))
	})

	t.Run("mixed alpha values (1-254)", func(t *testing.T) {
		t.Parallel()

		img := createMixedAlphaImage(t, 100, 100)
		assert.False(t, subject.HasTransparency(img))
	})

	t.Run("image too small (less than 10 pixels wide)", func(t *testing.T) {
		t.Parallel()

		img := createTransparentImage(t, 5, 100, true)
		assert.False(t, subject.HasTransparency(img))
	})

	t.Run("image too small (less than 1 pixel tall)", func(t *testing.T) {
		t.Parallel()

		img := createTransparentImage(t, 100, 0, true)
		assert.False(t, subject.HasTransparency(img))
	})
}

// createTransparentImage creates a synthetic PNG image with transparent or opaque pixels.
// Uses Go stdlib (image.NewRGBA, png.Encode) as per design decision.
func createTransparentImage(t *testing.T, width, height int, transparent bool) image.Image {
	t.Helper()

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	var c color.Color
	if transparent {
		c = color.RGBA{0, 0, 0, 0} // Fully transparent (alpha == 0)
	} else {
		c = color.RGBA{255, 255, 255, 255} // Fully opaque (alpha == 255)
	}

	// Fill entire image with the same color
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, c)
		}
	}

	return img
}

// createMixedAlphaImage creates an image with varying alpha values (1-254).
// Used to test strict zero check - should return false since not ALL pixels have alpha == 0.
func createMixedAlphaImage(t *testing.T, width, height int) image.Image {
	t.Helper()

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Set first 10 pixels to various non-zero alpha values
	for x := 0; x < 10 && x < width; x++ {
		alpha := uint8(x + 1) // Alpha values 1-10
		img.Set(x, 0, color.RGBA{255, 255, 255, alpha})
	}

	return img
}

func TestGetPhone_IPhoneMirror(t *testing.T) {
	t.Parallel()

	t.Run("detects IPhoneMirror by transparency + dimensions", func(t *testing.T) {
		t.Parallel()

		// Create 836×1840 image with transparent first row
		img := image.NewRGBA(image.Rect(0, 0, 836, 1840))
		// Set first 10 pixels to transparent
		for x := 0; x < 10; x++ {
			img.Set(x, 0, color.RGBA{0, 0, 0, 0})
		}

		phone, err := subject.GetPhone(img)
		require.NoError(t, err)
		assert.Equal(t, mobile.IPhoneMirror, phone)
	})

	t.Run("rejects 836×1840 without transparency", func(t *testing.T) {
		t.Parallel()

		// Create 836×1840 image WITHOUT transparency
		img := image.NewRGBA(image.Rect(0, 0, 836, 1840))
		// Set first 10 pixels to opaque
		for x := 0; x < 10; x++ {
			img.Set(x, 0, color.RGBA{255, 255, 255, 255})
		}

		phone, err := subject.GetPhone(img)
		assert.ErrorIs(t, err, subject.ErrUnsupportedPhone)
		assert.Equal(t, mobile.Phone{}, phone)
	})

	t.Run("rejects transparent image with wrong dimensions", func(t *testing.T) {
		t.Parallel()

		// Create 100×100 image with transparent first row (wrong dimensions)
		img := createTransparentImage(t, 100, 100, true)

		phone, err := subject.GetPhone(img)
		assert.ErrorIs(t, err, subject.ErrUnsupportedPhone)
		assert.Equal(t, mobile.Phone{}, phone)
	})
}

func TestCropMonth_MonthSize(t *testing.T) {
	t.Parallel()

	t.Run("produces 100px for IPhoneMirror", func(t *testing.T) {
		t.Parallel()

		img := image.NewRGBA(image.Rect(0, 0, 836, 1840))
		cropped := subject.CropMonth(img, mobile.IPhoneMirror)
		bounds := cropped.Bounds().Max

		assert.Equal(t, 836, bounds.X)
		assert.Equal(t, 100, bounds.Y) // IPhoneMirror has MonthSize=100
	})

	t.Run("produces 150px for other models", func(t *testing.T) {
		t.Parallel()

		img := image.NewRGBA(image.Rect(0, 0, 1206, 2622))
		cropped := subject.CropMonth(img, mobile.IPhone16Pro)
		bounds := cropped.Bounds().Max

		assert.Equal(t, 1206, bounds.X)
		assert.Equal(t, 150, bounds.Y) // IPhone16Pro has MonthSize=150
	})

	t.Run("fallback to 150px when MonthSize is 0", func(t *testing.T) {
		t.Parallel()

		// Create a phone with MonthSize=0 to test fallback
		phone := mobile.Phone{Width: 100, Height: 200, MonthSize: 0}
		img := image.NewRGBA(image.Rect(0, 0, 100, 200))
		cropped := subject.CropMonth(img, phone)
		bounds := cropped.Bounds().Max

		assert.Equal(t, 100, bounds.X)
		assert.Equal(t, 150, bounds.Y) // Should fallback to 150
	})
}

func TestCrop_FullFlow(t *testing.T) {
	t.Parallel()

	t.Run("iPhone Mirror full flow", func(t *testing.T) {
		t.Parallel()

		// Create iPhone Mirror PNG with transparency
		img := image.NewRGBA(image.Rect(0, 0, 836, 1840))
		// Set first row transparent (iPhone Mirror signature)
		for x := 0; x < 836; x++ {
			img.Set(x, 0, color.RGBA{0, 0, 0, 0})
		}

		// Verify GetPhone detects IPhoneMirror
		phone, err := subject.GetPhone(img)
		require.NoError(t, err)
		assert.Equal(t, mobile.IPhoneMirror, phone)

		// Verify Crop produces correct regions
		croppedImg := subject.CropImage(img, phone)
		croppedMonth := subject.CropMonth(img, phone)

		// Transaction area: 1840 - 600 - 180 = 1060px tall
		assert.Equal(t, 836, croppedImg.Bounds().Max.X)
		assert.Equal(t, 1060, croppedImg.Bounds().Max.Y)

		// Month area: 100px tall (MonthSize)
		assert.Equal(t, 836, croppedMonth.Bounds().Max.X)
		assert.Equal(t, 100, croppedMonth.Bounds().Max.Y)
	})
}
