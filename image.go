package main

import (
	"image"
	"image/color"
	"image/draw"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type Position int

const (
	_ Position = iota
	TOP
	BOTTOM
)

func addLabel(img image.Image, position Position, label string, face font.Face) image.Image {
	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, img.Bounds(), img, image.Point{}, draw.Src)

	col := color.RGBA{200, 100, 0, 255}
	point := fixed.Point26_6{}

	d := &font.Drawer{
		Dst:  rgba,
		Src:  image.NewUniform(col),
		Face: face,
		Dot:  point,
	}

	textWidth := d.MeasureString(label).Ceil()
	textHeight := (face.Metrics().Ascent + face.Metrics().Descent).Ceil()

	x := (img.Bounds().Dx() - textWidth) / 2
	point.X = fixed.I(x)

	var y int
	switch position {
	case TOP:
		y = textHeight // Already have some marge
	case BOTTOM:
		y = img.Bounds().Dy() - 20
	}
	point.Y = fixed.I(y)

	d.Dot = point
	d.DrawString(label)

	return rgba
}
