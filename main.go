package main

import (
	"fyne.io/fyne/v2/app"
)

func main() {
	a := app.New()

	mainWindow := createMainWindow(a)
	mainWindow.Show()

	a.Run()
}
