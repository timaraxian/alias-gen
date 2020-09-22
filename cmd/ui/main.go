package main

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/rivo/tview"
	"github.com/timaraxian/alias-gen/pkg/database"
	"github.com/timaraxian/alias-gen/pkg/tui"
)

func main() {
	config := tui.Config{}
	if _, err := toml.DecodeFile(os.Getenv("ALIASGEN_CONFIG"), &config); err != nil {
		log.Printf("Failed to open config file: %s\n", err)
		os.Exit(1)
	}
	ui := tview.NewApplication()

	app := tui.App{
		Config:          config,
		Ui:              ui,
		DBAL:            nil,
		PrevState:       "init",
		NextState:       "menu",
		Update:          false,
		Word:            tui.Word{},
		Pattern:         tui.Pattern{},
		WordListArgs:    database.WordListArgs{},
		PatternListArgs: database.PatternListArgs{},
	}

	var err error
	app.DBAL, err = database.Bootstrap(app.Config.DB)

	var list *tview.List
	var form *tview.Form
	var table *tview.Table

	for {
		if app.NextState == "stop" {
			break
		}
		if app.PrevState != app.NextState {
			app.Update = false
			switch app.NextState {
			case "menu":
				list = app.ShowMenu()
			case "addWord":
				form = app.ShowNewWord()
			case "submitWord":
				err = app.SubmitNewWord(form)
				if err != nil {
					panic(err)
				}
			case "listWords":
				table = app.ListWords()
			case "viewWord":
				list = app.ViewWord()
			case "editWordWord":
				form = app.ShowEditWordWord()
			case "submitWordWord":
				err = app.SubmitWordWord(form)
				if err != nil {
					panic(err)
				}
			case "editWordLanguage":
				form = app.ShowEditWordLanguage()
			case "submitWordLanguage":
				err = app.SubmitWordLanguage(form)
				if err != nil {
					panic(err)
				}
			case "editWordPart":
				form = app.ShowEditWordPart()
			case "submitWordPart":
				err = app.SubmitWordPart(form)
				if err != nil {
					panic(err)
				}
			case "addPattern":
				form = app.ShowNewPattern()
			case "submitPattern":
				err = app.SubmitNewPattern(form)
				if err != nil {
					panic(err)
				}
			case "listPatterns":
				table = app.ListPatterns()
			case "viewPattern":
				list = app.ViewPattern()
			case "editPatternPattern":
				form = app.ShowEditPatternPattern()
			case "submitPatternPattern":
				err = app.SubmitPatternPattern(form)
				if err != nil {
					panic(err)
				}
			case "editPatternLanguage":
				form = app.ShowEditPatternLanguage()
			case "submitPatternLanguage":
				err = app.SubmitPatternLanguage(form)
				if err != nil {
					panic(err)
				}
			}

		}

		if app.Update {
			switch app.PrevState {
			case "menu", "viewWord", "viewPattern":
				err := app.Ui.SetRoot(list, true).SetFocus(list).Run()
				if err != nil {
					panic(err)
				}
			case "addWord", "editWordWord", "editWordLanguage", "editWordPart", "addPattern", "editPatternPattern", "editPatternLanguage":
				err := app.Ui.SetRoot(form, true).SetFocus(form).Run()
				if err != nil {
					panic(err)
				}
			case "listWords", "listPatterns":
				err := app.Ui.SetRoot(table, true).SetFocus(table).Run()
				if err != nil {
					panic(err)
				}
			}

		}
	}
}
