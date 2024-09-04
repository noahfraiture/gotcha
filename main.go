package main

import (
	"gotcha/image"
	"gotcha/kitty"
	"gotcha/term"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"

	"golang.org/x/image/font"
)

func unwrap(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	writer := os.Stdout
	unwrap(term.Clear())

	var topLabel string
	var bottomLabel string
	var fontFace font.Face

	var fontPath string
	var imagePath string
	form := FormSelect(&fontPath, &imagePath)

	unwrap(form.Run())

	img, err := image.NewImage(imagePath)
	unwrap(err)
	var imgPrinter = kitty.NewPrinter(img)

	// Load the selected font with the calculated size
	fontFace, err = loadFontFromFile(fontPath, imgPrinter.FontSize())
	unwrap(err)

	unwrap(term.Clear())
	for {
		imgPrinter.Fprint(writer)

		form := FormCaption(&topLabel, &bottomLabel)

		unwrap(form.Run())

		imgPrinter.AddLabel(image.TOP, topLabel, fontFace)
		imgPrinter.AddLabel(image.BOTTOM, bottomLabel, fontFace)

		imgPrinter.Fprint(writer)

		unwrap(term.Clear())
	}
}
