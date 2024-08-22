package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"

	"github.com/charmbracelet/huh"
	"github.com/dolmen-go/kittyimg"
)

func clear() {
	switch runtime.GOOS {
	case "linux", "darwin":
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func show(img image.Image) {
	err := kittyimg.Fprintln(os.Stdout, img)
	if err != nil {
		log.Fatal("Error printing image ", err)
	}
}

func addLabel(img image.Image, x, y int, label string, face font.Face) image.Image {
	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, img.Bounds(), img, image.Point{}, draw.Src)

	col := color.RGBA{200, 100, 0, 255}
	point := fixed.Point26_6{Y: fixed.I(y)}

	d := &font.Drawer{
		Dst:  rgba,
		Src:  image.NewUniform(col),
		Face: face,
		Dot:  point,
	}

	// Measure the text width
	textWidth := d.MeasureString(label).Ceil()

	// Center the text
	imageWidth := img.Bounds().Dx()
	x = (imageWidth - textWidth) / 2
	point.X = fixed.I(x)

	d.Dot = point
	d.DrawString(label)

	return rgba
}

func loadFontFromFile(path string) (font.Face, error) {
	fontBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Parse the font file
	f, err := opentype.Parse(fontBytes)
	if err != nil {
		return nil, err
	}

	// Create a font.Face
	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size: 12, // Adjust as needed
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

func getFontOptions(dir string) ([]huh.Option[string], error) {
	var options []huh.Option[string]

	err := filepath.WalkDir(dir, func(path string, info os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && (filepath.Ext(path) == ".ttf" || filepath.Ext(path) == ".otf") {
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
	fontOptions, err := getFontOptions("fonts")
	if err != nil {
		log.Fatal("Error reading fonts directory: ", err)
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose a font").
				Options(fontOptions...).
				Value(&fontPath),

			huh.NewInput().
				Title("Enter image file name").
				Value(&imagePath).
				Validate(func(str string) error {
					if str == "" {
						return fmt.Errorf("Image file name cannot be empty")
					}
					return nil
				}),
		),
	)

	err = form.Run()
	if err != nil {
		log.Fatal(err)
	}

	// Load the selected font
	fontFace, err = loadFontFromFile(fontPath)
	if err != nil {
		log.Fatal("Error loading font: ", err)
	}

	origin := getImg(imagePath)
	img := origin

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

		img = addLabel(img, img.Bounds().Dx()/2, 20, topLabel, fontFace)
		img = addLabel(img, img.Bounds().Dx()/2, img.Bounds().Dy()-20, bottomLabel, fontFace)

		show(img)

		clear()
	}
}
