package tui

import (
	"strconv"
	"strings"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/timaraxian/alias-gen/pkg/database"
)

func (app *App) ShowNewWord() (form *tview.Form) {
	if app.NextState != "addWord" {
		panic("Invalid State")
	}
	form = tview.NewForm().
		AddInputField("word", "", 20, nil, func(text string) {
			app.processWordWord(text)
		}).
		AddInputField("language", "", 20, nil, func(text string) {
			app.processWordLanguage(text)
		}).
		AddInputField("part", "", 20, nil, func(text string) {
			app.processWordPart(text)
		}).
		AddButton("Add Word", func() {
			app.NextState = "submitWord"
			app.Ui.Stop()
		}).
		AddButton("Cancel", func() {
			app.NextState = "menu"
			app.Ui.Stop()
		})

	app.PrevState = "addWord"
	app.Update = true
	form.SetBorder(true).SetTitle("Add Word").SetTitleAlign(tview.AlignLeft)

	return form

}

func (app *App) processWordWord(text string) {
	app.Word.SetWord = strings.TrimSpace(text)
}

func (app *App) processWordLanguage(text string) {
	app.Word.SetLanguage = strings.TrimSpace(text)
}

func (app *App) processWordPart(text string) {
	app.Word.SetPart = strings.TrimSpace(text)
}

func (app *App) SubmitNewWord() (err error) {
	if app.NextState != "submitWord" {
		panic("Invalid State")
	}

	_, err = app.DBAL.WordCreate(app.Word.SetWord, app.Word.SetLanguage, app.Word.SetPart)
	app.PrevState = "submitWord"
	app.NextState = "menu"
	app.Update = true
	return err

}

func (app *App) ListWords() (table *tview.Table) {
	if app.NextState != "listWords" {
		panic("Invalid State")
	}

	// get words from DBAL
	listArgs := database.WordListArgs{}
	listArgs.Limit = &app.WordListArgs.Limit
	listArgs.Offset = &app.WordListArgs.Offset
	listArgs.OrderByWord = &app.WordListArgs.OrderByWord
	listArgs.DescWord = &app.WordListArgs.DescWord
	listArgs.OrderByLanguage = &app.WordListArgs.OrderByLanguage
	listArgs.DescLanguage = &app.WordListArgs.DescLanguage
	listArgs.OrderByPart = &app.WordListArgs.OrderByPart
	listArgs.DescPart = &app.WordListArgs.DescPart
	listArgs.OrderByUpdatedAt = &app.WordListArgs.OrderByUpdatedAt
	listArgs.DescUpdatedAt = &app.WordListArgs.DescUpdatedAt
	listArgs.OrderByCreatedAt = &app.WordListArgs.OrderByCreatedAt
	listArgs.DescCreatedAt = &app.WordListArgs.DescCreatedAt
	listArgs.ShowArchived = &app.WordListArgs.ShowArchived
	words, err := app.DBAL.WordList(listArgs)
	if err != nil {
		panic(err)
	}

	table = tview.NewTable().
		SetBorders(true)

	cols, rows := 7, len(words)+1

	// build header
	header := []string{"WordID", "Word", "Language", "Part", "CreatedAt", "UpdatedAt", "ArchivedAt"}

	for c := 0; c < cols; c++ {
		table.SetCell(0, c,
			tview.NewTableCell(header[c]).
				SetTextColor(tcell.ColorYellow).
				SetAlign(tview.AlignCenter))
	}

	// build content
	for r := 1; r < rows; r++ {
		for c := 0; c < cols; c++ {
			table.SetCell(r, c,
				tview.NewTableCell(getWordRowValue(words[r-1], c)).
					SetTextColor(tcell.ColorWhite).
					SetAlign(tview.AlignCenter))
		}
	}

	// table navigation
	table.Select(1, 0).SetFixed(1, 0).SetSelectable(true, false).SetSelectedFunc(func(row, col int) {
		app.Word.GetWordID = table.GetCell(row, col).Text
		app.NextState = "viewWord"
		app.Update = true
		app.Ui.Stop()
	}).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyESC {
			app.NextState = "menu"
			app.Update = true
			app.Ui.Stop()
		}
		if key == tcell.KeyTAB {
			app.NextState = "viewWordListArgs"
			app.Update = true
			app.Ui.Stop()
		}
	})

	app.PrevState = "listWords"
	app.Update = true

	table.SetBorder(true).SetTitle("Words (TAB for options || ESC for menu)").SetTitleAlign(tview.AlignLeft)

	return table
}

func getWordRowValue(row database.Word, c int) string {
	switch c {
	case 0:
		return row.WordID
	case 1:
		return row.Word
	case 2:
		return row.Language
	case 3:
		return row.Part
	case 4:
		return row.CreatedAt.Format("2006-01-02 15:04:05")
	case 5:
		return row.UpdatedAt.Format("2006-01-02 15:04:05")
	case 6:
		if row.ArchivedAt != nil {
			return row.ArchivedAt.Format("2006-01-02 15:04:05")
		}
	default:
		return ""
	}
	return ""
}

func (app *App) ShowWordListArgs() (form *tview.Form) {
	if app.NextState != "viewWordListArgs" {
		panic("Invalid State")
	}

	limit := strconv.Itoa(app.WordListArgs.Limit)
	offset := strconv.Itoa(app.WordListArgs.Offset)

	form = tview.NewForm().
		AddInputField("Limit", limit, 5, nil, func(text string) {
			app.processWordListArgsLimit(text)
		}).
		AddInputField("Offset", offset, 5, nil, func(text string) {
			app.processWordListArgsOffset(text)
		}).
		AddCheckbox("Order By Word", app.WordListArgs.OrderByWord, func(checked bool) {
			app.WordListArgs.OrderByWord = checked
		}).
		AddCheckbox("Descending Word", app.WordListArgs.DescWord, func(checked bool) {
			app.WordListArgs.DescWord = checked
		}).
		AddCheckbox("Order By Language", app.WordListArgs.OrderByLanguage, func(checked bool) {
			app.WordListArgs.OrderByLanguage = checked
		}).
		AddCheckbox("Descending Language", app.WordListArgs.DescLanguage, func(checked bool) {
			app.WordListArgs.DescLanguage = checked
		}).
		AddCheckbox("Order By Part", app.WordListArgs.OrderByPart, func(checked bool) {
			app.WordListArgs.OrderByPart = checked
		}).
		AddCheckbox("Descending Part", app.WordListArgs.DescPart, func(checked bool) {
			app.WordListArgs.DescPart = checked
		}).
		AddCheckbox("Order By UpdatedAt", app.WordListArgs.OrderByUpdatedAt, func(checked bool) {
			app.WordListArgs.OrderByUpdatedAt = checked
		}).
		AddCheckbox("Descending UpdatedAt", app.WordListArgs.DescUpdatedAt, func(checked bool) {
			app.WordListArgs.DescUpdatedAt = checked
		}).
		AddCheckbox("Order By CreatedAt", app.WordListArgs.OrderByCreatedAt, func(checked bool) {
			app.WordListArgs.OrderByCreatedAt = checked
		}).
		AddCheckbox("Descending CreatedAt", app.WordListArgs.DescCreatedAt, func(checked bool) {
			app.WordListArgs.DescCreatedAt = checked
		}).
		AddCheckbox("Show Archived", app.WordListArgs.ShowArchived, func(checked bool) {
			app.WordListArgs.ShowArchived = checked
		}).
		AddButton("Back to list", func() {
			app.updateWordListArgsLimitOffset()
			app.NextState = "listWords"
			app.Ui.Stop()
		})

	app.PrevState = "viewWordListArgs"
	app.Update = true
	form.SetBorder(true).SetTitle("Word List Args").SetTitleAlign(tview.AlignLeft)

	return form
}

func (app *App) processWordListArgsLimit(text string) {
	i, err := strconv.Atoi(strings.TrimSpace(text))
	if err != nil {
		app.Err = err
	} else {
		app.Err = nil
	}
	app.WordListArgs.Limit = i
}

func (app *App) processWordListArgsOffset(text string) {
	i, err := strconv.Atoi(strings.TrimSpace(text))
	if err != nil {
		app.Err = err
	} else {
		app.Err = nil
	}
	app.WordListArgs.Offset = i
}

func (app *App) updateWordListArgsLimitOffset() {
	if app.Err != nil {
		panic(app.Err)
	}
}

func (app *App) ViewWord() (list *tview.List) {
	if app.NextState != "viewWord" {
		panic("Invalid State")
	}

	word, err := app.DBAL.WordGet(app.Word.GetWordID)
	if err != nil {
		panic(err)
	}

	app.Word.SetWord = word.Word
	app.Word.SetLanguage = word.Language
	app.Word.SetPart = word.Part

	CreatedAt := word.CreatedAt.Format("2006-01-02 15:04:05")
	UpdatedAt := word.UpdatedAt.Format("2006-01-02 15:04:05")
	ArchivedAt := ""
	if word.ArchivedAt != nil {
		app.Word.Archive = true
		ArchivedAt = word.ArchivedAt.Format("2006-01-02 15:04:05")
	} else {
		app.Word.Archive = false
	}

	list = tview.NewList().
		AddItem("WordID", word.WordID, 'a', nil).
		AddItem("Word", word.Word, 'b', func() {
			app.NextState = "editWordWord"
			app.Ui.Stop()
		}).
		AddItem("Language", word.Language, 'c', func() {
			app.NextState = "editWordLanguage"
			app.Ui.Stop()
		}).
		AddItem("Part", word.Part, 'd', func() {
			app.NextState = "editWordPart"
			app.Ui.Stop()
		}).
		AddItem("CreatedAt", CreatedAt, 'e', nil).
		AddItem("UpdatedAt", UpdatedAt, 'f', nil).
		AddItem("ArchivedAt", ArchivedAt, 'g', func() {
			app.NextState = "editWordArchive"
			app.Ui.Stop()
		}).
		AddItem("Back to list", "", 'h', func() {
			app.NextState = "listWords"
			app.Ui.Stop()
		}).
		AddItem("Back to menu", "", 'i', func() {
			app.NextState = "menu"
			app.Ui.Stop()
		}).
		AddItem("Quit", "", 'q', func() {
			app.NextState = "stop"
			app.Ui.Stop()
		})

	app.PrevState = "viewWord"
	app.Update = true

	list.SetBorder(true).SetTitle("View Word").SetTitleAlign(tview.AlignLeft)

	return list
}

func (app *App) ShowEditWordWord() (form *tview.Form) {
	if app.NextState != "editWordWord" {
		panic("Invalid State")
	}
	form = tview.NewForm().
		AddInputField("word", app.Word.SetWord, 20, nil, func(text string) {
			app.processWordWord(text)
		}).
		AddButton("Edit Word", func() {
			app.NextState = "submitWordWord"
			app.Ui.Stop()
		}).
		AddButton("Cancel", func() {
			app.NextState = "viewWord"
			app.Ui.Stop()
		})

	app.PrevState = "editWordWord"
	app.Update = true
	form.SetBorder(true).SetTitle("Edit Word").SetTitleAlign(tview.AlignLeft)

	return form
}

func (app *App) SubmitWordWord() (err error) {
	if app.NextState != "submitWordWord" {
		panic("Invalid State")
	}

	err = app.DBAL.WordSetWord(app.Word.GetWordID, app.Word.SetWord)
	app.PrevState = "submitWordWord"
	app.NextState = "viewWord"
	app.Update = true
	return err
}

func (app *App) ShowEditWordLanguage() (form *tview.Form) {
	if app.NextState != "editWordLanguage" {
		panic("Invalid State")
	}
	form = tview.NewForm().
		AddInputField("language", app.Word.SetLanguage, 20, nil, func(text string) {
			app.processWordLanguage(text)
		}).
		AddButton("Edit Language", func() {
			app.NextState = "submitWordLanguage"
			app.Ui.Stop()
		}).
		AddButton("Cancel", func() {
			app.NextState = "viewWord"
			app.Ui.Stop()
		})

	app.PrevState = "editWordLanguage"
	app.Update = true
	form.SetBorder(true).SetTitle("Edit Word").SetTitleAlign(tview.AlignLeft)

	return form
}

func (app *App) SubmitWordLanguage() (err error) {
	if app.NextState != "submitWordLanguage" {
		panic("Invalid State")
	}

	err = app.DBAL.WordSetLanguage(app.Word.GetWordID, app.Word.SetLanguage)
	app.PrevState = "submitWordLanguage"
	app.NextState = "viewWord"
	app.Update = true
	return err
}

func (app *App) ShowEditWordPart() (form *tview.Form) {
	if app.NextState != "editWordPart" {
		panic("Invalid State")
	}
	form = tview.NewForm().
		AddInputField("part", app.Word.SetPart, 20, nil, func(text string) {
			app.processWordPart(text)
		}).
		AddButton("Edit Part", func() {
			app.NextState = "submitWordPart"
			app.Ui.Stop()
		}).
		AddButton("Cancel", func() {
			app.NextState = "viewWord"
			app.Ui.Stop()
		})

	app.PrevState = "editWordPart"
	app.Update = true
	form.SetBorder(true).SetTitle("Edit Word").SetTitleAlign(tview.AlignLeft)

	return form
}

func (app *App) SubmitWordPart() (err error) {
	if app.NextState != "submitWordPart" {
		panic("Invalid State")
	}

	err = app.DBAL.WordSetPart(app.Word.GetWordID, app.Word.SetPart)
	app.PrevState = "submitWordPart"
	app.NextState = "viewWord"
	app.Update = true
	return err
}

func (app *App) ShowEditWordArchive() (form *tview.Form) {
	if app.NextState != "editWordArchive" {
		panic("Invalid State")
	}

	form = tview.NewForm().
		AddCheckbox("archived", app.Word.Archive, func(checked bool) {
			app.processWordArchived(checked)
		}).
		AddButton("Edit Archive", func() {
			app.NextState = "submitWordArchive"
			app.Ui.Stop()
		}).
		AddButton("Cancel", func() {
			app.NextState = "viewWord"
			app.Ui.Stop()
		})

	app.PrevState = "editWordArchive"
	app.Update = true
	form.SetBorder(true).SetTitle("Edit Word").SetTitleAlign(tview.AlignLeft)

	return form
}

func (app *App) processWordArchived(checked bool) {
	app.Word.Archive = checked
}

func (app *App) SubmitWordArchive() (err error) {
	if app.NextState != "submitWordArchive" {
		panic("Invalid State")
	}

	if app.Word.Archive {
		err = app.DBAL.WordSetArchive(app.Word.GetWordID)
	} else {
		err = app.DBAL.WordSetUnArchive(app.Word.GetWordID)
	}

	app.PrevState = "submitWordArchive"
	app.NextState = "viewWord"
	app.Update = true
	return err
}
