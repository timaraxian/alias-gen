package tui

import (
	"strings"

	"github.com/rivo/tview"
)

func (app *App) SelectLanguage() (form *tview.Form) {
	if app.NextState != "selectLanguage" {
		panic("Invalid State")
	}

	languages, err := app.DBAL.GetDistinctLanguages()
	if err != nil {
		panic(err)
	}
	if len(languages) < 1 {
		panic("No languages")
	}

	form = tview.NewForm().
		AddDropDown("Language", languages, 0, func(option string, idx int) {
			app.Random.language = option
		}).
		AddButton("Generate Alias", func() {
			app.NextState = "showRandomAlias"
			app.Ui.Stop()
		}).
		AddButton("Cancel", func() {
			app.NextState = "menu"
			app.Ui.Stop()
		})

	app.PrevState = "selectLanguage"
	app.Update = true

	form.SetBorder(true).SetTitle("Select Language").SetTitleAlign(tview.AlignLeft)

	return form
}

func (app *App) ShowRandomAlias() (modal *tview.Modal) {
	if app.NextState != "showRandomAlias" {
		panic("Invalid State")
	}

	for {
		pattern, err := app.DBAL.PatternRandom(app.Random.language)
		if err != nil {
			panic(err)
		}
		patternSlice := strings.Split(pattern.Pattern, ",")
		if len(patternSlice) < 1 {
			panic("No Patterns")
		}

		for i := 0; i < len(patternSlice); i++ {
			word, err := app.DBAL.WordRandom(app.Random.language, patternSlice[i])
			if err != nil {
				app.Random.empty = true
				continue
			}
			app.Random.alias = append(app.Random.alias, word.Word)
		}
		if !app.Random.empty {
			break
		}
	}

	alias := strings.Join(app.Random.alias, " ")
	app.Random.alias = []string{}

	modal = tview.NewModal().
		SetText("Random alias").
		SetText(alias).
		AddButtons([]string{"Change Language", "Menu", "Quit"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Change Language" {
				app.NextState = "selectLanguage"
				app.Ui.Stop()
			}
			if buttonLabel == "Menu" {
				app.NextState = "menu"
				app.Ui.Stop()
			}
			if buttonLabel == "Quit" {
				app.NextState = "stop"
				app.Ui.Stop()
			}
		})

	app.PrevState = "showRandomAlias"
	app.Update = true

	return modal
}
