package main

import (
	"fmt"
	"math/rand/v2"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func createTimerWindow(a fyne.App, name string, duration time.Duration) {
	win := a.NewWindow(name)
	tw := NewTimerWindow(name, duration)

	img := tw.flower.assets[0]

	remaining := binding.NewString()
	if err := remaining.Set(tw.timer.remaining.String()); err != nil {
		fmt.Printf("ERR: Failed to set remaining timer: %+v\n", err)
		panic(err)
	}

	label := widget.NewLabelWithData(remaining)

	content := container.NewVBox(
		widget.NewLabel(name),
		img,
		container.NewHBox(
			layout.NewSpacer(),
			label,
			layout.NewSpacer(),
		),
	)

	go func() {
		tw.timer.OnTick(func() error {
			chance := rand.IntN(11)
			if chance == 10 {
				tw.flower.Grow()
				*img = *tw.flower.assets[tw.flower.stage-1]
			}

			err := remaining.Set(tw.timer.remaining.String())
			fyne.Do(content.Refresh)

			return err
		})
	}()

	win.SetContent(content)
	win.Show()
}

func createMainWindow(a fyne.App) fyne.Window {
	win := a.NewWindow("Main Window")

	timerNameInput := widget.NewEntry()
	timerNameInput.SetPlaceHolder("Timer Name...")

	hourInput := NewTimeInput()
	minsInput := NewTimeInput()
	secsInput := NewTimeInput()

	startButton := widget.NewButton("Create Timer", func() {
		hours := hourInput.Get()
		mins := minsInput.Get()
		secs := secsInput.Get()

		duration := time.Hour*time.Duration(hours) +
			time.Minute*time.Duration(mins) +
			time.Second*time.Duration(secs)

		createTimerWindow(a, timerNameInput.Text, duration)
	})

	win.SetContent(container.NewVBox(
		widget.NewLabelWithStyle(
			"Create a new timer",
			fyne.TextAlignCenter,
			fyne.TextStyle{
				Bold: true,
			},
		),
		timerNameInput,
		container.NewVBox(
			widget.NewLabel("Timer Duration"),
			container.NewHBox(
				container.NewVBox(
					widget.NewLabel("Hours"),
					hourInput,
				),
				container.NewVBox(
					widget.NewLabel("Minutes"),
					minsInput,
				),
				container.NewVBox(
					widget.NewLabel("Seconds"),
					secsInput,
				),
			),
		),
		startButton,
	))

	return win
}
