package image_test

import (
	"bytes"
	"image"
	"image/png"
	"testing"

	subject "git.home/c6bank-transactions/internal/image"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCropImage(t *testing.T) {
	t.Parallel()

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

			croppedImg, err := subject.Crop(bufSeeker)
			require.NoError(t, err)

			croppedPNG, err := png.Decode(croppedImg)
			require.NoError(t, err)

			height := test.height - subject.HeaderHeight - subject.FooterHeight

			assert.Equal(t, subject.Width, croppedPNG.Bounds().Max.X)
			assert.Equal(t, height, croppedPNG.Bounds().Max.Y)
		})
	}
}
