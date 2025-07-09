//go:generate fyne bundle -o assets.go assets/hyacinth/purple/hp_0001.png

package main

import (
	"fmt"
	"math/rand/v2"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// Once again, fuck your errors
func getIntFromBoundString(bs binding.String) int {
	val, err := bs.Get()
	if err != nil {
		panic(err)
	}

	num, err := strconv.Atoi(val)
	if err != nil {
		panic(err)
	}

	return num
}

type TimerData struct {
	flowerName binding.String
	timerName  binding.String

	hours   binding.String
	minutes binding.String
	seconds binding.String
}

func newTimerData() TimerData {
	return TimerData{
		flowerName: binding.NewString(),
		timerName:  binding.NewString(),
		hours:      binding.NewString(),
		minutes:    binding.NewString(),
		seconds:    binding.NewString(),
	}
}

func (t *TimerData) getDuration() time.Duration {
	hours := getIntFromBoundString(t.hours)
	minutes := getIntFromBoundString(t.minutes)
	seconds := getIntFromBoundString(t.seconds)

	return time.Hour*time.Duration(hours) +
		time.Minute*time.Duration(minutes) +
		time.Second*time.Duration(seconds)
}

func spawnTimerWindow(a fyne.App, data TimerData, flower *Flower) {
	name, err := data.timerName.Get()
	if err != nil {
		fmt.Printf("ERR: Failed to get Timer Name: %+v\n", err)
		panic(err)
	}

	duration := data.getDuration()
	timer := NewTimer(duration)

	remaining := binding.NewString()
	if err := remaining.Set(timer.remaining.String()); err != nil {
		fmt.Printf("ERR: Failed to set remaining timer: %+v\n", err)
		panic(err)
	}

	remainingLabel := widget.NewLabelWithData(remaining)
	// This is a problem, we are over-writing the image every time we set it. need to figure out how to not do that
	// image := flower.stages[0]
	image := &canvas.Image{}
	*image = *flower.stages[0]

	content := container.NewVBox(
		widget.NewLabel(name),
		image,
		container.NewHBox(
			layout.NewSpacer(),
			remainingLabel,
			layout.NewSpacer(),
		),
	)
	win := a.NewWindow(name)

	flowerTicker := flower.Update(image)
	go func() {
		timer.Run(
			func() error {
				if err := remaining.Set(timer.remaining.String()); err != nil {
					panic(err)
				}

				fyne.Do(content.Refresh)
				return nil
			},
			func() error {
				println("Timer Completed")
				return nil
			},
		)
	}()

	win.SetOnClosed(func() {
		flowerTicker.Stop()
	})

	win.SetContent(content)
	win.Show()
}

type App struct {
	app     fyne.App
	flowers []*Flower

	data TimerData
}

func newApp() (*App, error) {
	flowers, err := LoadFlowers()
	if err != nil {
		return nil, err
	}

	a := app.New()

	return &App{
		app:     a,
		flowers: flowers,
		data:    newTimerData(),
	}, nil
}

func (a *App) createFlowerSelect() *widget.Select {
	keys := make([]string, len(a.flowers))
	for idx, flower := range a.flowers {
		keys[idx] = flower.Name
	}

	randomFlower := a.flowers[rand.IntN(len(a.flowers))].Name
	if err := a.data.flowerName.Set(randomFlower); err != nil {
		panic(err)
	}

	return widget.NewSelectWithData(keys, a.data.flowerName)
}

func (a *App) createTimerNameInput() *widget.Entry {
	entry := widget.NewEntryWithData(a.data.timerName)

	entry.SetPlaceHolder("Timer Name...")

	// TODO: Implement a random timer name feature
	//       (if the user doesn't add a timer name themselves)

	return entry
}

func (a *App) reset() {
	if err := a.data.timerName.Set(""); err != nil {
		fmt.Printf("ERR: Failed to reset TimerName: %+v\n", err)
	}
	if err := a.data.hours.Set("0"); err != nil {
		fmt.Printf("ERR: Failed to reset Hours: %+v\n", err)
	}
	if err := a.data.minutes.Set("0"); err != nil {
		fmt.Printf("ERR: Failed to reset Minutes: %+v\n", err)
	}
	if err := a.data.seconds.Set("0"); err != nil {
		fmt.Printf("ERR: Failed to reset Seconds: %+v\n", err)
	}
}

func (a *App) ShowMainWindow() {
	win := a.app.NewWindow("Create Timers")

	win.SetContent(container.NewVBox(
		widget.NewLabelWithStyle("Create Timer", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		container.NewGridWithColumns(
			2,
			container.NewVBox(
				widget.NewLabel("Timer Name"),
				a.createTimerNameInput(),
			),
			container.NewVBox(
				widget.NewLabel("Choose Flower"),
				a.createFlowerSelect(),
			),
		),
		container.NewGridWithColumns(
			3,
			container.NewVBox(
				widget.NewLabelWithStyle("Hours", fyne.TextAlignCenter, fyne.TextStyle{Italic: true}),
				NewTimeEntry(a.data.hours),
			),
			container.NewVBox(
				widget.NewLabelWithStyle("Minutes", fyne.TextAlignCenter, fyne.TextStyle{Italic: true}),
				NewTimeEntry(a.data.minutes),
			),
			container.NewVBox(
				widget.NewLabelWithStyle("Seconds", fyne.TextAlignCenter, fyne.TextStyle{Italic: true}),
				NewTimeEntry(a.data.seconds),
			),
		),
		widget.NewButton(
			"Start Timer",
			func() {
				flowerName, err := a.data.flowerName.Get()
				if err != nil {
					fmt.Printf("ERR: Failed to get flower name: %+v\n", err)
					panic(err)
				}

				var flower *Flower = nil

				for _, f := range a.flowers {
					if f.Name == flowerName {
						flower = f
					}
				}
				if flower == nil {
					panic("flower " + flowerName + " not found")
				}

				spawnTimerWindow(a.app, a.data, flower)

				a.reset()
			},
		),
	))

	win.Resize(fyne.NewSize(320, 0))

	win.Show()
}

func (a *App) Run() {
	a.ShowMainWindow()

	a.app.Run()
}

func main() {
	a, err := newApp()
	if err != nil {
		panic(err)
	}

	a.Run()
}
