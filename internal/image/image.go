package image

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg" // Import JPEG format
	"image/png"  // Import PNG format
	"io"
)

const (
	Width        = 1100
	HeaderHeight = 800
	FooterHeight = 250
)

func Crop(file io.ReadSeeker) (io.Reader, error) {
	var img image.Image
	var err error

	_, format, err := image.DecodeConfig(file)
	if err != nil {
		return nil, err
	}

	if _, err = file.Seek(0, 0); err != nil {
		return nil, err
	}

	switch format {
	case "jpg":
		img, err = jpeg.Decode(file)
	case "png":
		img, err = png.Decode(file)
	default:
		return nil, fmt.Errorf("unsupported image format: %s", format)
	}

	if err != nil {
		return nil, err
	}

	croppedImg := cropImage(img)

	var buf bytes.Buffer

	switch format {
	case "jpg":
		err = jpeg.Encode(&buf, croppedImg, nil)
	case "png":
		err = png.Encode(&buf, croppedImg)
	}

	if err != nil {
		return nil, err
	}

	return &buf, nil
}

func cropImage(img image.Image) *image.RGBA {
	bounds := img.Bounds()
	height := bounds.Max.Y

	newHeight := height - HeaderHeight - FooterHeight

	rect := image.Rect(0, HeaderHeight, Width, height-FooterHeight)

	croppedImg := image.NewRGBA(image.Rect(0, 0, Width, newHeight))

	draw.Draw(croppedImg, croppedImg.Bounds(), img, rect.Min, draw.Src)

	return croppedImg
}
