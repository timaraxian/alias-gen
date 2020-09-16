package application

import (
	"net/http"
	"time"
)

type WordCreateArgs struct {
	Word     string `json:word`
	Language string `json:language`
	Part     string `json:part`
}

type WordCreateReply struct {
	WordID     string     `json:wordID`
	Word       string     `json:word`
	Language   string     `json:language`
	Part       string     `json:part`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
	ArchivedAt *time.Time `json:"archivedAt"`
}

func (app *App) WordCreate(w http.ResponseWriter, r *http.Request) {
	args := WordCreateArgs{}
	if err := app.decodeRequest(r, &args); err != nil {
		app.respondApi(w, r, nil, err)
	}

	word, err := app.DBAL.WordCreate(args.Word, args.Language, args.Part)
	if err != nil {
		app.respondApi(w, r, nil, err)
	}

	app.respondApi(w, r, WordCreateReply{
		WordID:     word.WordID,
		Language:   word.Language,
		Part:       word.Part,
		CreatedAt:  word.CreatedAt,
		UpdatedAt:  word.UpdatedAt,
		ArchivedAt: word.ArchivedAt,
	}, nil)

}
