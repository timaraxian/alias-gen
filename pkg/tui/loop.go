package tui

import "github.com/rivo/tview"

func (app *App) Loop() (err error) {

	var list *tview.List
	var form *tview.Form
	var table *tview.Table
	var modal *tview.Modal

	for {
		if app.NextState == "stop" {
			break
		}
		if app.PrevState != app.NextState {
			app.Update = false
			switch app.NextState {
			case "menu":
				list = app.ShowMenu()

				//Words
			case "addWord":
				form = app.ShowNewWord()
			case "submitWord":
				err = app.SubmitNewWord()
				if err != nil {
					return err
				}
			case "listWords":
				table = app.ListWords()
			case "viewWordListArgs":
				form = app.ShowWordListArgs()
			case "viewWord":
				list = app.ViewWord()
			case "editWordWord":
				form = app.ShowEditWordWord()
			case "submitWordWord":
				err = app.SubmitWordWord()
				if err != nil {
					return err
				}
			case "editWordLanguage":
				form = app.ShowEditWordLanguage()
			case "submitWordLanguage":
				err = app.SubmitWordLanguage()
				if err != nil {
					return err
				}
			case "editWordPart":
				form = app.ShowEditWordPart()
			case "submitWordPart":
				err = app.SubmitWordPart()
				if err != nil {
					return err
				}
			case "editWordArchive":
				form = app.ShowEditWordArchive()
			case "submitWordArchive":
				err = app.SubmitWordArchive()
				if err != nil {
					return err
				}

				//Patterns
			case "addPattern":
				form = app.ShowNewPattern()
			case "submitPattern":
				err = app.SubmitNewPattern()
				if err != nil {
					return err
				}
			case "listPatterns":
				table = app.ListPatterns()
			case "viewPatternListArgs":
				form = app.ShowPatternListArgs()
			case "viewPattern":
				list = app.ViewPattern()
			case "editPatternPattern":
				form = app.ShowEditPatternPattern()
			case "submitPatternPattern":
				err = app.SubmitPatternPattern()
				if err != nil {
					return err
				}
			case "editPatternLanguage":
				form = app.ShowEditPatternLanguage()
			case "submitPatternLanguage":
				err = app.SubmitPatternLanguage()
				if err != nil {
					return err
				}
			case "editPatternArchive":
				form = app.ShowEditPatternArchive()
			case "submitPatternArchive":
				err = app.SubmitPatternArchive()
				if err != nil {
					return err
				}

				//random
			case "selectLanguage":
				form = app.SelectLanguage()
			case "showRandomAlias":
				modal = app.ShowRandomAlias()
			}
		}

		if app.Update {
			switch app.PrevState {
			case "menu", "viewWord", "viewPattern":
				err := app.Ui.SetRoot(list, true).SetFocus(list).Run()
				if err != nil {
					return err
				}
			case "addWord", "editWordWord", "editWordLanguage", "editWordPart", "editWordArchive", "viewWordListArgs", "addPattern", "editPatternPattern", "editPatternLanguage", "editPatternArchive", "viewPatternListArgs", "selectLanguage":
				err := app.Ui.SetRoot(form, true).SetFocus(form).Run()
				if err != nil {
					return err
				}
			case "listWords", "listPatterns":
				err := app.Ui.SetRoot(table, true).SetFocus(table).Run()
				if err != nil {
					return err
				}
			case "showRandomAlias", "err":
				err := app.Ui.SetRoot(modal, true).SetFocus(modal).Run()
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
