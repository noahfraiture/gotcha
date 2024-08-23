package main

import (
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/image/font"

	"github.com/charmbracelet/huh"
)

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func getFileOptions(dir string, extensions []string) ([]huh.Option[string], error) {

	var options []huh.Option[string]

	err := filepath.WalkDir(dir, func(path string, info os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && contains(extensions, filepath.Ext(path)) {
			options = append(options, huh.NewOption(info.Name(), path))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return options, nil
}

func main() {
	clear()

	var imagePath string
	var topLabel string
	var bottomLabel string
	var fontPath string
	var fontFace font.Face

	// Get font options from fonts directory
	fontOptions, err := getFileOptions("fonts", []string{".ttf", ".otf"})
	if err != nil {
		log.Fatal("Error reading fonts directory: ", err)
	}
	imgOptions, err := getFileOptions("imgs", []string{".jpg", ".jpeg", ".png"})
	if err != nil {
		log.Fatal("Error reading imgs directory: ", err)
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose a font").
				Options(fontOptions...).
				Value(&fontPath),
			huh.NewSelect[string]().
				Title("Choose an image").
				Options(imgOptions...).
				Value(&imagePath),
		),
	)

	err = form.Run()
	if err != nil {
		log.Fatal(err)
	}

	origin := getImg(imagePath)
	img := origin

	// Calculate font size based on image dimensions
	imageHeight := img.Bounds().Dy()
	fontSize := float64(imageHeight) * 0.05 // Font size is 5% of image height

	// Load the selected font with the calculated size
	fontFace, err = loadFontFromFile(fontPath, fontSize)
	if err != nil {
		log.Fatal("Error loading font: ", err)
	}

	clear()
	for {
		show(img)

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewText().
					Title("Top Caption").
					Value(&topLabel),

				huh.NewText().
					Title("Bottom Caption").
					Value(&bottomLabel),
			),
		)

		err := form.Run()
		if err != nil {
			log.Fatal(err)
		}

		img = addLabel(img, TOP, topLabel, fontFace)
		img = addLabel(img, BOTTOM, bottomLabel, fontFace)

		show(img)

		clear()
	}
}
