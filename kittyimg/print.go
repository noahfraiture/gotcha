// Package kittyimg provides utilities to show image in a graphic terminal emulator supporting kitty's "terminal graphics protocol".
//
// See https://sw.kovidgoyal.net/kitty/graphics-protocol.html.
package kittyimg

import (
	"fmt"
	"image"
	"io"
	"math"
	"os"

	"golang.org/x/sys/unix"
)

func scaleImage(img image.Image, maxWidth, maxHeight int) (image.Image, error) {
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
			srcX := int(float64(x) / factor)
			srcY := int(float64(y) / factor)
			newImage.Set(x, y, img.At(srcX, srcY))
		}
	}

	return newImage, nil
}

func getTerminalSize() (int, int, error) {
	f, err := os.OpenFile("/dev/tty", unix.O_NOCTTY|unix.O_CLOEXEC|unix.O_NDELAY|unix.O_RDWR, 0666)
	if err != nil {
		return 0, 0, err
	}
	defer f.Close()

	sz, err := unix.IoctlGetWinsize(int(f.Fd()), unix.TIOCGWINSZ)
	if err != nil {
		return 0, 0, err
	}

	return int(sz.Col), int(sz.Row), nil
}

func Fprint(w io.Writer, img image.Image) error {
	// Get terminal size
	termWidth, termHeight, err := getTerminalSize()
	if err != nil {
		return err
	}

	// Scale the image to fit the terminal
	scaledImg, err := scaleImage(img, termWidth, termHeight)
	if err != nil {
		return err
	}

	bounds := scaledImg.Bounds()

	// f=32 => RGBA
	_, err = fmt.Fprintf(w, "\033_Gq=1,a=T,f=32,s=%d,v=%d,t=d,", bounds.Dx(), bounds.Dy())
	if err != nil {
		return err
	}

	buf := make([]byte, 0, 16384) // Multiple of 4 (RGBA)

	var p zlibPayload
	p.Reset(w)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if len(buf) == cap(buf) {
				if _, err = p.Write(buf); err != nil {
					return err
				}
				buf = buf[:0]
			}
			r, g, b, a := scaledImg.At(x, y).RGBA()
			buf = append(buf, byte(r>>8), byte(g>>8), byte(b>>8), byte(a>>8))
		}
	}

	if _, err = p.Write(buf); err != nil {
		return err
	}
	return p.Close()
}

func Fprintln(w io.Writer, img image.Image) error {
	err := Fprint(w, img)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte{'\n'})
	return err
}
