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
	return cropImage(img, phone.Month, mobile.MonthSize)
}

func cropImage(img image.Image, height, size int) *image.RGBA {
	bounds := img.Bounds()
	width := bounds.Max.X

	rect := image.Rect(0, height, width, height+size)
	croppedImg := image.NewRGBA(image.Rect(0, 0, width, size))
	draw.Draw(croppedImg, croppedImg.Bounds(), img, rect.Min, draw.Src)

	return croppedImg
}
