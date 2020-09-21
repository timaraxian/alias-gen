package main

import (
	"github.com/rivo/tview"
	"github.com/timaraxian/alias-gen/pkg/tui"
)

func main() {
	app := tview.NewApplication()

	ui := tui.Ui{
		app, "init", "list", false,
	}

	var menu *tview.List
	var form *tview.Form

	for {
		if ui.NextState == "stop" {
			break
		}
		if ui.PrevState != ui.NextState {
			ui.Update = false
			switch ui.NextState {
			case "list":
				menu = tui.ShowMenuModal(&ui)
			case "addWord":
				form = tui.ShowNewWordModal(&ui)
			}
		}

		if ui.Update {
			switch ui.PrevState {
			case "list":
				err := ui.App.SetRoot(menu, true).SetFocus(menu).Run()
				if err != nil {
					panic(err)
				}
			case "addWord":
				err := ui.App.SetRoot(form, true).SetFocus(form).Run()
				if err != nil {
					panic(err)
				}
			}
		}
	}
}
