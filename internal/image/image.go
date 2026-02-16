package image

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg" // Import JPEG format
	"image/png"  // Import PNG format
	"io"

	"git.home/c6bank-transactions/internal/mobile"
)

var ErrUnsupportedPhone = errors.New("unsupported phone")

// HasTransparency checks if the first 10 pixels of the first row have alpha == 0.
// Used to detect iPhone Mirror screenshots which have a transparent header.
// Returns true only if ALL 10 pixels are fully transparent (alpha == 0).
func HasTransparency(img image.Image) bool {
	bounds := img.Bounds()
	if bounds.Max.X < 10 || bounds.Max.Y < 1 {
		return false
	}

	// Check first 10 pixels of first row (y=0, x=0 to x=9)
	for x := 0; x < 10; x++ {
		c := img.At(x, 0)
		_, _, _, a := c.RGBA()
		// RGBA returns alpha in range 0..65535, where 0 is fully transparent
		if a != 0 {
			return false
		}
	}

	return true
}

func Crop(file io.ReadSeeker) (io.Reader, io.Reader, error) {
	var img image.Image
	var err error

	_, format, err := image.DecodeConfig(file)
	if err != nil {
		return nil, nil, err
	}

	if _, err = file.Seek(0, 0); err != nil {
		return nil, nil, err
	}

	switch format {
	case "jpg":
		img, err = jpeg.Decode(file)
	case "png":
		img, err = png.Decode(file)
	default:
		return nil, nil, fmt.Errorf("unsupported image format: %s", format)
	}

	if err != nil {
		return nil, nil, err
	}

	phone, err := GetPhone(img)
	if err != nil {
		return nil, nil, err
	}

	croppedImg := CropImage(img, phone)
	croppedMonth := CropMonth(img, phone)

	var (
		ierr, merr error
		imageBuf   bytes.Buffer
		monthBuf   bytes.Buffer
	)

	switch format {
	case "jpg":
		ierr = jpeg.Encode(&imageBuf, croppedImg, nil)
		merr = jpeg.Encode(&monthBuf, croppedMonth, nil)

	case "png":
		ierr = png.Encode(&imageBuf, croppedImg)
		merr = png.Encode(&monthBuf, croppedMonth)
	}

	err = errors.Join(ierr, merr)
	if err != nil {
		return nil, nil, err
	}

	return &imageBuf, &monthBuf, nil
}

func GetPhone(img image.Image) (mobile.Phone, error) {
	bounds := img.Bounds()

	// Check for iPhone Mirror: transparency + exact dimensions
	if HasTransparency(img) && bounds.Max.X == 836 && bounds.Max.Y == 1840 {
		return mobile.IPhoneMirror, nil
	}

	// Regular dimension-based detection for other models
	for _, phone := range mobile.Phones {
		if bounds.Max.X == phone.Width && bounds.Max.Y == phone.Height {
			return phone, nil
		}
	}

	return mobile.Phone{}, ErrUnsupportedPhone
}

func CropImage(img image.Image, phone mobile.Phone) *image.RGBA {
	size := img.Bounds().Max.Y - phone.Header - phone.Footer

	return cropImage(img, phone.Header, size)
}

func CropMonth(img image.Image, phone mobile.Phone) *image.RGBA {
	// Use phone.MonthSize if set, otherwise fallback to 150px
	monthSize := phone.MonthSize
	if monthSize == 0 {
		monthSize = 150
	}
	return cropImage(img, phone.Month, monthSize)
}

func cropImage(img image.Image, height, size int) *image.RGBA {
	bounds := img.Bounds()
	width := bounds.Max.X

	rect := image.Rect(0, height, width, height+size)
	croppedImg := image.NewRGBA(image.Rect(0, 0, width, size))
	draw.Draw(croppedImg, croppedImg.Bounds(), img, rect.Min, draw.Src)

	return croppedImg
}
