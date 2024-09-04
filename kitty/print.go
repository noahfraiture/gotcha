// Package kittyimg provides utilities to show image in a graphic terminal emulator supporting kitty's "terminal graphics protocol".
//
// See https://sw.kovidgoyal.net/kitty/graphics-protocol.html.
package kitty

import (
	"fmt"
	"gotcha/image"
	"gotcha/term"
	"io"
)

type KittyPrinter struct {
	image.ImageFile
}

func NewPrinter(img image.ImageFile) KittyPrinter {
	return KittyPrinter{img}
}

// Print the image to the provided writer.
// We use local client implementation only since it's easier.
// If I add image generation AI we could add remote implementation.
func (k *KittyPrinter) Fprint(w io.Writer) error {
	termWidth, termHeight, err := term.GetTerminalSize()
	if err != nil {
		return err
	}

	k.ScaleImage(termWidth, termHeight)

	// Header sent to the terminal emulator giving information about the image
	// https://sw.kovidgoyal.net/kitty/graphics-protocol/#transferring-pixel-data
	// We don't use compression
	// Kitty give the possibility to directly handle png, it complexifie the code,
	// If it's too slow, it could be implemented
	// TODO unicode placeholders
	bounds := k.Edited.Bounds()
	_, err = fmt.Fprintf(w, "\033_Gf=32,s=%d,v=%d,t=d;", bounds.Dx(), bounds.Dy())

	if err != nil {
		return err
	}

	buf := make([]byte, 0, bounds.Dx()*bounds.Dy()) // Multiple of 4 (RGBA)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := k.Edited.At(x, y).RGBA()
			buf = append(buf, byte(r>>8), byte(g>>8), byte(b>>8), byte(a>>8))
		}
	}

	if _, err = w.Write(buf); err != nil {
		return err
	}
	_, err = fmt.Fprintf(w, "\033")
	if err != nil {
		return err
	}
	return nil
}

func (k *KittyPrinter) Fprintln(w io.Writer) error {
	err := k.Fprint(w)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte{'\n'})
	return err
}
