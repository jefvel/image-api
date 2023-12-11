package util

import (
	"bytes"
	"errors"
	"image"
	"image-api/pkg/repository"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"time"
)

var (
	ErrInvalidFileFormat = errors.New("Invalid file format")
)

func ExtractImageMetadata(imageData []byte) (*repository.Metadata, error) {
	img, format, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		if errors.Is(err, image.ErrFormat) {
			return nil, ErrInvalidFileFormat
		}
		return nil, err
	}

	bounds := img.Bounds()

	sizeInKb := len(imageData) / (1_024)

	metadata := &repository.Metadata{
		Width:     bounds.Dx(),
		Height:    bounds.Dy(),
		Size:      sizeInKb,
		Format:    format,
		CreatedAt: time.Now(),
	}

	return metadata, nil
}

func CropImage(imageData []byte, rect image.Rectangle) ([]byte, error) {
	img, format, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		if errors.Is(err, image.ErrFormat) {
			return nil, ErrInvalidFileFormat
		}

		return nil, err
	}

	newImg := image.NewRGBA(img.Bounds())
	draw.Draw(newImg, newImg.Bounds(), img, image.Point{}, draw.Over)

	croppedImage := newImg.SubImage(rect)

	buf := new(bytes.Buffer)
	switch format {
	case "jpeg":
		jpeg.Encode(buf, croppedImage, nil)
	case "gif":
		gif.Encode(buf, croppedImage, nil)
	case "png":
		png.Encode(buf, croppedImage)
	default:
		return nil, ErrInvalidFileFormat
	}

	return buf.Bytes(), nil
}
