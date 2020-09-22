package tui

import (
	"github.com/rivo/tview"
	"github.com/timaraxian/alias-gen/pkg/database"
)

type App struct {
	Config Config

	Ui   *tview.Application
	DBAL *database.DBAL

	PrevState string
	NextState string
	Update    bool

	Word    Word
	Pattern Pattern

	WordListArgs    database.WordListArgs
	PatternListArgs database.PatternListArgs
}

type Config struct {
	DB database.Config
}

type Word struct {
	GetWordID   string
	SetWord     string
	SetLanguage string
	SetPart     string
	Archive     bool
}

type Pattern struct {
	GetPatternID string
	SetPattern   string
	SetLanguage  string
	Archive      bool
}
