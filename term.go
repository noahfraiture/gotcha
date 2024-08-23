package main

import (
	"image"
	"log"
	"os"
	"os/exec"
	"runtime"

	"gotcha/kittyimg"
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

// TODO : add auto scale
func show(img image.Image) {
	err := kittyimg.Fprintln(os.Stdout, img)
	if err != nil {
		log.Fatal("Error printing image ", err)
	}
}
