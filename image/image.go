package image

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"
	"os"
	"path/filepath"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type ImageFile struct {
	origin    image.Image
	Edited    image.Image // public to be used in printer
	path      string
	extension string
}

func NewImage(path string) (ImageFile, error) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening image:", err)
		return ImageFile{}, err
	}
	defer file.Close()

	// Decode the image using image.Decode, which automatically handles different formats
	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("Error decoding image:", err)
		return ImageFile{}, err
	}

	ext := filepath.Ext(path)
	return ImageFile{origin: img, Edited: img, path: path, extension: ext}, nil
}

// Calculate font size based on image dimensions
func (imgFile *ImageFile) FontSize() float64 {
	imageHeight := imgFile.origin.Bounds().Dy()
	return float64(imageHeight) * 0.05 // Font size is 5% of image height
}

// Return a new image with the label. Does NOT mutate the imgFile.
// Since we don't save the image we will have to recompute the scale every time.
// We consider it is negligeable performance for now.
func (imgFile *ImageFile) ScaleImage(maxWidth, maxHeight int) {
	bounds := imgFile.Edited.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	widthFactor := float64(width) / float64(maxWidth)
	heightFactor := float64(height) / float64(maxHeight)
	factor := math.Max(widthFactor, heightFactor)

	var resizedWidth, resizedHeight int
	resizedWidth = width / int(factor)
	resizedHeight = height / int(factor)

	newImage := image.NewRGBA(image.Rect(0, 0, resizedWidth, resizedHeight))
	for y := 0; y < resizedHeight; y++ {
		for x := 0; x < resizedWidth; x++ {
			srcX := int(float64(x) * factor)
			srcY := int(float64(y) * factor)
			newImage.Set(x, y, imgFile.Edited.At(srcX, srcY))
		}
	}

	imgFile.Edited = newImage
}

type Position int

const (
	_ Position = iota
	TOP
	BOTTOM
)

// Mutate 'edited' field from the 'origin'
func (imgFile *ImageFile) AddLabel(position Position, label string, face font.Face) {
	rgba := image.NewRGBA(imgFile.origin.Bounds())
	draw.Draw(rgba, imgFile.origin.Bounds(), imgFile.origin, image.Point{}, draw.Src)

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

	x := (imgFile.origin.Bounds().Dx() - textWidth) / 2
	point.X = fixed.I(x)

	var y int
	switch position {
	case TOP:
		y = textHeight // Already have some marge
	case BOTTOM:
		y = imgFile.origin.Bounds().Dy() - 20
	}
	point.Y = fixed.I(y)

	d.Dot = point
	d.DrawString(label)

	imgFile.Edited = rgba
}
