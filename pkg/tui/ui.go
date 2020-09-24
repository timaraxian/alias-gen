package tui

import (
	"github.com/rivo/tview"
	"github.com/timaraxian/alias-gen/pkg/database"
)

type App struct {
	Config Config

	Ui   *tview.Application
	DBAL *database.DBAL

	Err error

	PrevState string
	NextState string
	Update    bool

	Word    Word
	Pattern Pattern

	WordListArgs    WordListArgs
	PatternListArgs PatternListArgs

	Random Random
}

type Service func(app *App) (err error)

func NewApp(config Config, services []Service) (a *App, err error) {
	a = &App{
		Ui:              tview.NewApplication(),
		PrevState:       "init",
		NextState:       "menu",
		Update:          false,
		Word:            Word{},
		Pattern:         Pattern{},
		WordListArgs:    WordListArgs{},
		PatternListArgs: PatternListArgs{},
		Random:          Random{},
	}

	a.Config = config

	for _, s := range services {
		if err = s(a); err != nil {
			return a, err
		}
	}

	return a, nil
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

type WordListArgs struct {
	Limit            int
	Offset           int
	OrderByWord      bool
	DescWord         bool
	OrderByLanguage  bool
	DescLanguage     bool
	OrderByPart      bool
	DescPart         bool
	OrderByUpdatedAt bool
	DescUpdatedAt    bool
	OrderByCreatedAt bool
	DescCreatedAt    bool
	ShowArchived     bool
}

type PatternListArgs struct {
	Limit            int
	Offset           int
	OrderByPattern   bool
	DescPattern      bool
	OrderByLanguage  bool
	DescLanguage     bool
	OrderByUpdatedAt bool
	DescUpdatedAt    bool
	OrderByCreatedAt bool
	DescCreatedAt    bool
	ShowArchived     bool
}

type Random struct {
	language string
	pattern  []string
	alias    []string
	empty    bool
}
