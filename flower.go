package main

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

type Flower struct {
	AssetPath  string `json:"asset_path"`
	Name       string `json:"name"`
	GrowthRate int    `json:"growth_rate"`

	stages       []*canvas.Image `json:"-"`
	currentStage int             `json:"-"`
}

func (f *Flower) LoadStages() {
	f.stages = make([]*canvas.Image, 0)
	f.currentStage = 0

	err := filepath.WalkDir(f.AssetPath, func(path string, d fs.DirEntry, err error) error {
		if filepath.Ext(path) != ".png" {
			return nil
		}

		img := canvas.NewImageFromFile(path)
		if img.Hidden {
			panic("image should not be hidden")
		}
		img.ScaleMode = canvas.ImageScalePixels
		img.SetMinSize(fyne.NewSize(144, 144))

		f.stages = append(f.stages, img)

		return nil
	})

	if err != nil {
		panic(err)
	}
}

func LoadFlowers() ([]*Flower, error) {
	bytes, err := os.ReadFile("flowers.json")
	if err != nil {
		return nil, err
	}

	var flowers []*Flower

	if err := json.Unmarshal(bytes, &flowers); err != nil {
		return nil, err
	}

	for _, flower := range flowers {
		flower.LoadStages()
	}

	return flowers, nil
}

func (f *Flower) Update(img *canvas.Image) *time.Ticker {
	ticker := time.NewTicker(time.Second * time.Duration(f.GrowthRate))

	go func() {
		for range ticker.C {
			f.currentStage++
			if f.currentStage >= len(f.stages) {
				f.currentStage = 0
			}

			*img = *f.stages[f.currentStage]
		}
	}()

	return ticker
}
