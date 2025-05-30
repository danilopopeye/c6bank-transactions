package image_test

import (
	"bytes"
	"image"
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
