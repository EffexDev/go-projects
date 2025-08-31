package main

import (
	"fmt"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	// Create a new Fyne app
	myApp := app.New()
	myWindow := myApp.NewWindow("Fyne Test")

	// Add a label and a button
	label := widget.NewLabel("Hello, Fyne!")
	button := widget.NewButton("Quit", func() {
		myApp.Quit()
	})

	// Create a dropdown with options
	options := []string{"Option 1", "Option 2", "Option 3"}
	dropdown := widget.NewSelect(options, func(selected string) {
		// This function is called when the user selects an option
		fmt.Println("Selected:", selected)
	})

	// Arrange widgets vertically
	content := container.NewVBox(label, button, dropdown)
	myWindow.SetContent(content)

	// Show the window and run the app

	myWindow.ShowAndRun()
}
