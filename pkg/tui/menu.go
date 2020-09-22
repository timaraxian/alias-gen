package tui

import "github.com/rivo/tview"

func (app *App) ShowMenu() (list *tview.List) {
	if app.NextState != "menu" {
		panic("Invalid State")
	}
	list = tview.NewList().
		AddItem("Add a word", "Add a new word to the database", 'a', func() {
			app.NextState = "addWord"
			app.Ui.Stop()
		}).
		AddItem("List words", "List all of the words in the database", 'b', func() {
			app.NextState = "listWords"
			app.Ui.Stop()
		}).
		AddItem("Add a pattern", "Add a new pattern to the database", 'c', func() {
			app.NextState = "addPattern"
			app.Ui.Stop()
		}).
		AddItem("List patterns", "List all of the patterns in the database", 'd', func() {
			app.NextState = "listPatterns"
			app.Ui.Stop()
		}).
		AddItem("Quit", "Press to exit", 'q', func() {
			app.NextState = "stop"
			app.Ui.Stop()
		})

	app.PrevState = "menu"
	app.Update = true
	return list

}
