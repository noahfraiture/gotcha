package printer

import (
	"image"
	"io"
	"math"
	"os"

	"golang.org/x/sys/unix"
)

type Printer interface {
	Fprint(w io.Writer, img image.Image) error
	Fprintln(w io.Writer, img image.Image) error
}

func ScaleImage(img image.Image, maxWidth, maxHeight int) (image.Image, error) {
	bounds := img.Bounds()
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
			newImage.Set(x, y, img.At(srcX, srcY))
		}
	}

	return newImage, nil
}

func GetTerminalSize() (int, int, error) {
	f, err := os.OpenFile("/dev/tty", unix.O_NOCTTY|unix.O_CLOEXEC|unix.O_NDELAY|unix.O_RDWR, 0666)
	if err != nil {
		return 0, 0, err
	}
	defer f.Close()

	sz, err := unix.IoctlGetWinsize(int(f.Fd()), unix.TIOCGWINSZ)
	if err != nil {
		return 0, 0, err
	}

	return int(sz.Xpixel), int(sz.Ypixel), nil
}
