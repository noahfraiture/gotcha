// Package kittyimg provides utilities to show image in a graphic terminal emulator supporting kitty's "terminal graphics protocol".
//
// See https://sw.kovidgoyal.net/kitty/graphics-protocol.html.
package kitty

import (
	"fmt"
	"image"
	"io"

	"gotcha/printer"
)

type KittyPrinter struct{}

func (_ KittyPrinter) Fprint(w io.Writer, img image.Image) error {
	// Get terminal size
	termWidth, termHeight, err := printer.GetTerminalSize()
	if err != nil {
		return err
	}

	// Scale the image to fit the terminal
	scaledImg, err := printer.ScaleImage(img, termWidth, termHeight)
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

func (k KittyPrinter) Fprintln(w io.Writer, img image.Image) error {
	err := k.Fprint(w, img)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte{'\n'})
	return err
}
