package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/charmbracelet/huh"
)

func FormCaption(topLabel, bottomLabel *string) *huh.Form {

	return huh.NewForm(
		huh.NewGroup(
			huh.NewText().
				Title("Top Caption").
				Value(topLabel),

			huh.NewText().
				Title("Bottom Caption").
				Value(bottomLabel),
		),
	)
}

func FormSelect(fontPath, imagePath *string) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			fontSelect(fontPath),
			imageSelect(imagePath),
		),
	)

}

func fontSelect(fontPath *string) huh.Field {
	fontOptions, err := getFileOptions("fonts", []string{".ttf", ".otf"})
	if err != nil {
		log.Fatal("Error reading fonts directory: ", err)
	}
	return huh.NewSelect[string]().
		Title("Choose a font").
		Options(fontOptions...).
		Value(fontPath)
}

func imageSelect(imagePath *string) huh.Field {
	imgOptions, err := getFileOptions("imgs", []string{".jpg", ".jpeg", ".png", ".gif"})
	if err != nil {
		log.Fatal("Error reading imgs directory: ", err)
	}

	return huh.NewSelect[string]().
		Title("Choose an image").
		Options(imgOptions...).
		Value(imagePath)
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
func contains(s []string, e string) bool {

	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
