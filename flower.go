package main

import (
	"fmt"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

type Flower struct {
	name   string
	assets []*canvas.Image
	stage  int
}

func NewFlower(name string) Flower {
	assets := make([]*canvas.Image, 5)

	for idx := 1; idx <= 5; idx++ {
		img := canvas.NewImageFromFile(
			filepath.Join("assets", name, fmt.Sprintf("%d.png", idx)),
		)

		img.ScaleMode = canvas.ImageScalePixels
		img.SetMinSize(fyne.NewSize(144, 144))
		assets[idx-1] = img
	}

	return Flower{
		name:   name,
		assets: assets,
		stage:  1,
	}
}

func (f *Flower) Grow() {
	if f.stage < 5 {
		f.stage++
	} else {
		f.stage = 1
	}
}
