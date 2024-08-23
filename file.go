package main

import (
	"image"
	"log"
	"os"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

func loadFontFromFile(path string, size float64) (font.Face, error) {
	fontBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Parse the font file
	f, err := opentype.Parse(fontBytes)
	if err != nil {
		return nil, err
	}

	// Create a font.Face with a dynamically determined size
	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size: size,
		DPI:  72,
	})
	if err != nil {
		return nil, err
	}

	return face, nil
}

func getImg(name string) image.Image {
	file, err := os.Open(name)
	if err != nil {
		log.Fatal("Error opening ", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatal("Error decoding ", err)
	}
	return img

}
