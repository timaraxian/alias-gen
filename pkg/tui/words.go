package tui

import (
	"github.com/rivo/tview"
)

func ShowNewWordModal(ui *Ui) (form *tview.Form) {
	if ui.NextState != "addWord" {
		panic("Invalid State")
	}
	form = tview.NewForm().
		AddInputField("word", "", 20, nil, nil).
		AddInputField("language", "", 20, nil, nil).
		AddInputField("part", "", 20, nil, nil).
		AddButton("Add Word", func() {
			ui.NextState = "stop"
			ui.App.Stop()
		}).
		AddButton("Cancel", func() {
			ui.NextState = "list"
			ui.App.Stop()
		})

	ui.PrevState = "addWord"
	ui.Update = true
	form.SetBorder(true).SetTitle("Add Word").SetTitleAlign(tview.AlignLeft)

	return form

}

func ShowMenuModal(ui *Ui) (list *tview.List) {
	if ui.NextState != "list" {
		panic("Invalid State")
	}
	list = tview.NewList().
		AddItem("Add a word", "Add a new word to the database", 'a', func() {
			ui.NextState = "addWord"
			ui.App.Stop()
		}).
		AddItem("List words", "List all of the words in the database", 'b', nil).
		AddItem("Add a pattern", "Add a new pattern to the database", 'c', nil).
		AddItem("List patterns", "List all of the patterns in the database", 'd', nil).
		AddItem("Quit", "Press to exit", 'q', func() {
			ui.NextState = "stop"
			ui.App.Stop()
		})

	ui.PrevState = "list"
	ui.Update = true
	return list

}
