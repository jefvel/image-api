package model

import (
	"bytes"
	"errors"
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
)

var (
	ErrInvalidFileFormat = errors.New("Invalid file format")
)

type Image struct {
	Width  int
	Height int
	Size   int
	Format string
	Bytes  []byte

	img image.Image
}

func ImageFromBytes(imageData []byte) (*Image, error) {
	img, format, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		if errors.Is(err, image.ErrFormat) {
			return nil, ErrInvalidFileFormat
		}
		return nil, err
	}

	bounds := img.Bounds()

	sizeInKb := len(imageData) / (1_024)

	metadata := &Image{
		Width:  bounds.Dx(),
		Height: bounds.Dy(),
		Size:   sizeInKb,
		Format: format,
		Bytes:  imageData,
		img:    img,
	}

	return metadata, nil
}

func (i *Image) Crop(rect image.Rectangle) (*Image, error) {
	newImg := image.NewRGBA(i.img.Bounds())
	draw.Draw(newImg, newImg.Bounds(), i.img, image.Point{}, draw.Over)

	croppedImage := newImg.SubImage(rect)

	buf := new(bytes.Buffer)
	switch i.Format {
	case "jpeg":
		jpeg.Encode(buf, croppedImage, nil)
	case "gif":
		gif.Encode(buf, croppedImage, nil)
	case "png":
		png.Encode(buf, croppedImage)
	default:
		return nil, ErrInvalidFileFormat
	}

	return ImageFromBytes(buf.Bytes())
}
