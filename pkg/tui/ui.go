package tui

import "github.com/rivo/tview"

type Ui struct {
	App       *tview.Application
	PrevState string
	NextState string
	Update    bool
}
