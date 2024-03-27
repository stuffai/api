package img

import (
	"bytes"
	"image"
	"image/jpeg"
	"io"

	"github.com/disintegration/imaging"
)

const sz = 512

func ProcessImage(src io.Reader) (*bytes.Buffer, error) {
	// Decode and check the image
	img, format, err := image.Decode(src)
	if err != nil {
		return nil, err
	}

	// Only allow JPEG and PNG formats
	if format != "jpeg" && format != "png" {
		return nil, err
	}

	// Use the resize package to check dimensions without resizing
	m := imaging.Resize(img, sz, sz, imaging.Lanczos)

	// Store into buffer and return
	buf := new(bytes.Buffer)
	if err := jpeg.Encode(buf, m, nil); err != nil {
		return nil, err
	}

	return buf, nil
}
