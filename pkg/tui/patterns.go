package tui

import (
	"strings"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/timaraxian/alias-gen/pkg/database"
)

func (app *App) ShowNewPattern() (form *tview.Form) {
	if app.NextState != "addPattern" {
		panic("Invalid State")
	}
	form = tview.NewForm().
		AddInputField("pattern", "", 20, nil, func(text string) {
			app.processPatternPattern(text)
		}).
		AddInputField("language", "", 20, nil, func(text string) {
			app.processPatternLanguage(text)
		}).
		AddButton("Add Pattern", func() {
			app.NextState = "submitPattern"
			app.Ui.Stop()
		}).
		AddButton("Cancel", func() {
			app.NextState = "menu"
			app.Ui.Stop()
		})

	app.PrevState = "addPattern"
	app.Update = true
	form.SetBorder(true).SetTitle("Add Pattern").SetTitleAlign(tview.AlignLeft)

	return form

}

func (app *App) processPatternPattern(text string) {
	s := strings.Split(text, ",")
	for i := 0; i < len(s); i++ {
		s[i] = strings.TrimSpace(s[i])
	}
	app.Pattern.SetPattern = strings.Join(s, ",")
}

func (app *App) processPatternLanguage(text string) {
	app.Pattern.SetLanguage = strings.TrimSpace(text)
}

func (app *App) SubmitNewPattern() (err error) {
	if app.NextState != "submitPattern" {
		panic("Invalid State")
	}

	_, err = app.DBAL.PatternCreate(app.Pattern.SetPattern, app.Pattern.SetLanguage)
	app.PrevState = "submitPattern"
	app.NextState = "menu"
	app.Update = true
	return err

}

func (app *App) ListPatterns() (table *tview.Table) {
	if app.NextState != "listPatterns" {
		panic("Invalid State")
	}

	// get patterns from DBAL todo: list args
	patterns, err := app.DBAL.PatternList(app.PatternListArgs)
	if err != nil {
		panic(err)
	}

	table = tview.NewTable().
		SetBorders(true)

	cols, rows := 6, len(patterns)+1

	// build header
	header := []string{"PatternID", "Pattern", "Language", "CreatedAt", "UpdatedAt", "ArchivedAt"}

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
				tview.NewTableCell(getPatternRowValue(patterns[r-1], c)).
					SetTextColor(tcell.ColorWhite).
					SetAlign(tview.AlignCenter))
		}
	}

	// table navigation
	table.Select(1, 0).SetFixed(1, 0).SetSelectable(true, false).SetSelectedFunc(func(row, col int) {
		app.Pattern.GetPatternID = table.GetCell(row, col).Text
		app.NextState = "viewPattern"
		app.Update = true
		app.Ui.Stop()
	}).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyESC {
			app.NextState = "menu"
			app.Update = true
			app.Ui.Stop()
		}
	})

	app.PrevState = "listPatterns"
	app.Update = true

	table.SetBorder(true).SetTitle("Patterns").SetTitleAlign(tview.AlignLeft)

	return table
}

func getPatternRowValue(row database.Pattern, c int) string {
	switch c {
	case 0:
		return row.PatternID
	case 1:
		return row.Pattern
	case 2:
		return row.Language
	case 3:
		return row.CreatedAt.Format("2006-01-02 15:04:05")
	case 4:
		return row.UpdatedAt.Format("2006-01-02 15:04:05")
	case 5:
		if row.ArchivedAt != nil {
			return row.ArchivedAt.Format("2006-01-02 15:04:05")
		}
	default:
		return ""
	}
	return ""
}

func (app *App) ViewPattern() (list *tview.List) {
	if app.NextState != "viewPattern" {
		panic("Invalid State")
	}

	pattern, err := app.DBAL.PatternGet(app.Pattern.GetPatternID)
	if err != nil {
		panic(err)
	}

	app.Pattern.SetPattern = pattern.Pattern
	app.Pattern.SetLanguage = pattern.Language

	CreatedAt := pattern.CreatedAt.Format("2006-01-02 15:04:05")
	UpdatedAt := pattern.UpdatedAt.Format("2006-01-02 15:04:05")
	ArchivedAt := ""
	if pattern.ArchivedAt != nil {
		app.Pattern.Archive = true
		ArchivedAt = pattern.ArchivedAt.Format("2006-01-02 15:04:05")
	} else {
		app.Pattern.Archive = false
	}

	list = tview.NewList().
		AddItem("PatternID", pattern.PatternID, 'a', nil).
		AddItem("Pattern", pattern.Pattern, 'b', func() {
			app.NextState = "editPatternPattern"
			app.Ui.Stop()
		}).
		AddItem("Language", pattern.Language, 'c', func() {
			app.NextState = "editPatternLanguage"
			app.Ui.Stop()
		}).
		AddItem("CreatedAt", CreatedAt, 'd', nil).
		AddItem("UpdatedAt", UpdatedAt, 'e', nil).
		AddItem("ArchivedAt", ArchivedAt, 'f', func() {
			app.NextState = "editPatternArchive"
			app.Ui.Stop()
		}).
		AddItem("Back to list", "", 'g', func() {
			app.NextState = "listPatterns"
			app.Ui.Stop()
		}).
		AddItem("Back to menu", "", 'h', func() {
			app.NextState = "menu"
			app.Ui.Stop()
		}).
		AddItem("Quit", "", 'q', func() {
			app.NextState = "stop"
			app.Ui.Stop()
		})

	app.PrevState = "viewPattern"
	app.Update = true

	list.SetBorder(true).SetTitle("View Pattern").SetTitleAlign(tview.AlignLeft)

	return list
}

func (app *App) ShowEditPatternPattern() (form *tview.Form) {
	if app.NextState != "editPatternPattern" {
		panic("Invalid State")
	}
	form = tview.NewForm().
		AddInputField("pattern", app.Pattern.SetPattern, 20, nil, func(text string) {
			app.processPatternPattern(text)
		}).
		AddButton("Edit Pattern", func() {
			app.NextState = "submitPatternPattern"
			app.Ui.Stop()
		}).
		AddButton("Cancel", func() {
			app.NextState = "viewPattern"
			app.Ui.Stop()
		})

	app.PrevState = "editPatternPattern"
	app.Update = true
	form.SetBorder(true).SetTitle("Edit Pattern").SetTitleAlign(tview.AlignLeft)

	return form
}

func (app *App) SubmitPatternPattern() (err error) {
	if app.NextState != "submitPatternPattern" {
		panic("Invalid State")
	}

	err = app.DBAL.PatternSetPattern(app.Pattern.GetPatternID, app.Pattern.SetPattern)
	app.PrevState = "submitPatternPattern"
	app.NextState = "viewPattern"
	app.Update = true
	return err
}

func (app *App) ShowEditPatternLanguage() (form *tview.Form) {
	if app.NextState != "editPatternLanguage" {
		panic("Invalid State")
	}
	form = tview.NewForm().
		AddInputField("language", app.Pattern.SetLanguage, 20, nil, func(text string) {
			app.processPatternLanguage(text)
		}).
		AddButton("Edit Language", func() {
			app.NextState = "submitPatternLanguage"
			app.Ui.Stop()
		}).
		AddButton("Cancel", func() {
			app.NextState = "viewPattern"
			app.Ui.Stop()
		})

	app.PrevState = "editPatternLanguage"
	app.Update = true
	form.SetBorder(true).SetTitle("Edit Pattern").SetTitleAlign(tview.AlignLeft)

	return form
}

func (app *App) SubmitPatternLanguage() (err error) {
	if app.NextState != "submitPatternLanguage" {
		panic("Invalid State")
	}

	err = app.DBAL.PatternSetLanguage(app.Pattern.GetPatternID, app.Pattern.SetLanguage)
	app.PrevState = "submitPatternLanguage"
	app.NextState = "viewPattern"
	app.Update = true
	return err
}

func (app *App) ShowEditPatternArchive() (form *tview.Form) {
	if app.NextState != "editPatternArchive" {
		panic("Invalid State")
	}

	form = tview.NewForm().
		AddCheckbox("archived", app.Pattern.Archive, func(checked bool) {
			app.processPatternArchived(checked)
		}).
		AddButton("Edit Archive", func() {
			app.NextState = "submitPatternArchive"
			app.Ui.Stop()
		}).
		AddButton("Cancel", func() {
			app.NextState = "viewPattern"
			app.Ui.Stop()
		})

	app.PrevState = "editPatternArchive"
	app.Update = true
	form.SetBorder(true).SetTitle("Edit Pattern").SetTitleAlign(tview.AlignLeft)

	return form
}

func (app *App) processPatternArchived(checked bool) {
	app.Pattern.Archive = checked
}

func (app *App) SubmitPatternArchive() (err error) {
	if app.NextState != "submitPatternArchive" {
		panic("Invalid State")
	}

	if app.Pattern.Archive {
		err = app.DBAL.PatternSetArchive(app.Pattern.GetPatternID)
	} else {
		err = app.DBAL.PatternSetUnArchive(app.Pattern.GetPatternID)
	}

	app.PrevState = "submitPatternArchive"
	app.NextState = "viewPattern"
	app.Update = true
	return err
}
