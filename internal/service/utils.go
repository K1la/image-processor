package service

import (
	"fmt"
	res "github.com/nfnt/resize"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

var allowedFormats = map[string]struct{}{
	"jpeg": {},
	"jpg":  {},
	"png":  {},
	"gif":  {},
}

var tasks = map[string]struct{}{
	Resize:    struct{}{},
	Watermark: struct{}{},
	Thumbnail: struct{}{},
}

func isAllowedFormat(format string) bool {
	_, ok := allowedFormats[strings.ToLower(format)]
	return ok
}

func checkFormat(contentType string) (string, error) {
	parsed := strings.Split(contentType, "/")
	if len(parsed) < 2 {
		return "", ErrInvalidImageFormat
	}

	format := parsed[1]
	if format == "jpg" {
		format = "jpeg"
	}

	if !isAllowedFormat(format) {
		return "", ErrInvalidImageFormat
	}

	return format, nil
}

func isCorrectTask(task string) bool {
	_, ok := tasks[task]
	return ok
}

func resize(format string, r *os.File, w *os.File, width, height int) error {
	src, err := decode(format, r)
	if err != nil {
		return err
	}

	resized := res.Resize(uint(width), uint(height), src, res.Lanczos3)

	return encode(format, resized, w)
}

func decode(format string, r *os.File) (image.Image, error) {
	defer r.Close()

	switch format {
	case "jpeg":
		return jpeg.Decode(r)
	case "png":
		return png.Decode(r)
	case "gif":
		return gif.Decode(r)
	}

	return nil, fmt.Errorf("invalid file format")
}

func encode(format string, dst image.Image, w *os.File) error {
	defer w.Close()

	switch format {
	case "jpeg":
		return jpeg.Encode(w, dst, nil)
	case "png":
		return png.Encode(w, dst)
	case "gif":
		return gif.Encode(w, dst, nil)
	}

	return fmt.Errorf("invalid file format")
}

func addLabel(img draw.Image, x, y int, label string, fontSize float64) error {
	bytes, err := os.Open("web/font.ttf")
	if err != nil {
		return err
	}

	ftBytes, err := io.ReadAll(bytes)
	if err != nil {
		return err
	}

	ft, err := opentype.Parse(ftBytes)
	if err != nil {
		return err
	}

	face, err := opentype.NewFace(ft, &opentype.FaceOptions{
		Size:    fontSize,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return err
	}

	col := color.RGBA{255, 20, 100, 100} //red

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: face,
		Dot: fixed.Point26_6{
			X: fixed.I(x),
			Y: fixed.I(y),
		},
	}
	d.DrawString(label)

	return nil
}
