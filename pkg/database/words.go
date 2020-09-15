package database

import (
	"database/sql"
	"time"

	"github.com/timaraxian/hotel-gen/pkg/errors"
	"github.com/timaraxian/hotel-gen/pkg/helpers/crypto"
	"github.com/timaraxian/hotel-gen/pkg/helpers/validators"
)

type Word struct {
	WordID     string     `json:"wordID"`
	Word       string     `json:"word"`
	Language   string     `json:language`
	Part       string     `json:part`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
	ArchivedAt *time.Time `json:"archivedAt"`
}

func (dbal *DBAL) WordCreate(word_in, language, part string) (word Word, err error) {
	word.WordID = crypto.NewUUID()

	// todo: validate word, language, part
	word.Word = word_in
	word.Language = language
	word.Part = part

	word.CreatedAt = time.Now()
	word.UpdatedAt = word.CreatedAt

	stmt := `INSERT INTO words (
		word_id,
		word,
		language,
                part,
		created_at,
		updated_at,
		archived_at
	) VALUES ($1, $2, $3, $4, $5, $6, NULL);`

	_, err = dbal.Exec(stmt,
		word.WordID,
		word.Word,
		word.Language,
		word.Part,
		word.CreatedAt,
		word.UpdatedAt,
	)

	if err == nil {
		return word, nil
	}

	if dbIsDuplicateErr(err, "words_language_part_word") {
		return word, errors.WordDuplicate
	}

	return word, errors.UnexpectedError(err, "Failed creating word")
}

func (dbal *DBAL) WordGet(wordID string) (word Word, err error) {
	if err := validators.UUID(wordID); err != nil {
		return word, errors.WordNotFound
	}

	stmt := `SELECT
                word_id,
                word,
                language,
                part,
                created_at,
                updated_at,
                archived_at FROM words WHERE word_id=$1;`

	err = dbal.QueryRow(stmt, wordID).Scan(
		&word.WordID,
		&word.Word,
		&word.Language,
		&word.Part,
		&word.CreatedAt,
		&word.UpdatedAt,
		&word.ArchivedAt,
	)

	if err == nil {
		return word, nil
	}

	if err == sql.ErrNoRows {
		return word, errors.WordNotFound
	}

	return word, errors.UnexpectedError(err, "Failed getting word")
}

func (dbal DBAL) WordSetWord(wordID, word string) (err error) {
	if err := validators.UUID(wordID); err != nil {
		return errors.WordNotFound
	}

	// todo: validate word

	stmt := `UPDATE words SET word=$1, updated_at=$2 WHERE word_id=$3 AND archived_at IS NULL;`

	_, n, err := dbal.ExecOne(stmt, word, time.Now(), wordID)
	if dbIsDuplicateErr(err, "words_language_part_word") {
		return errors.WordDuplicate
	}
	if err != nil {
		return errors.UnexpectedError(err, "Failed to set word")
	} else if n == 0 {
		return errors.WordNotFound
	}

	return nil
}

func (dbal DBAL) WordSetLanguage(wordID, language string) (err error) {
	if err := validators.UUID(wordID); err != nil {
		return errors.WordNotFound
	}

	// todo: validate language
	stmt := `UPDATE words SET language=$1, updated_at=$2 WHERE word_id=$3 AND archived_at IS NULL;`

	_, n, err := dbal.ExecOne(stmt, language, time.Now(), wordID)
	if dbIsDuplicateErr(err, "words_language_part_word") {
		return errors.WordDuplicate
	}
	if err != nil {
		return errors.UnexpectedError(err, "Failed to set language")
	} else if n == 0 {
		return errors.WordNotFound
	}

	return nil
}

func (dbal DBAL) WordSetPart(wordID, part string) (err error) {
	if err := validators.UUID(wordID); err != nil {
		return errors.WordNotFound
	}

	// todo: validate part
	stmt := `UPDATE words SET part=$1, updated_at=$2 WHERE word_id=$3 AND archived_at IS NULL;`

	_, n, err := dbal.ExecOne(stmt, part, time.Now(), wordID)
	if dbIsDuplicateErr(err, "words_language_part_word") {
		return errors.WordDuplicate
	}
	if err != nil {
		return errors.UnexpectedError(err, "Failed to set part")
	} else if n == 0 {
		return errors.WordNotFound
	}

	return nil
}

func (dbal DBAL) WordSetArchive(wordID string) (err error) {
	if err := validators.UUID(wordID); err != nil {
		return errors.WordNotFound
	}

	stmt := `UPDATE words SET archived_at=COALESCE(archived_at, NOW()) WHERE word_id=$1;`

	_, n, err := dbal.ExecOne(stmt, wordID)
	if err != nil {
		return errors.UnexpectedError(err, "Failed to archive word")
	} else if n == 0 {
		return errors.WordNotFound
	}

	return nil
}

func (dbal DBAL) WordSetUnArchive(wordID string) (err error) {
	if err := validators.UUID(wordID); err != nil {
		return errors.WordNotFound
	}

	stmt := `UPDATE words SET archived_at=NULL WHERE word_id=$1;`

	_, n, err := dbal.ExecOne(stmt, wordID)
	if err != nil {
		return errors.UnexpectedError(err, "Failed to archive word")
	} else if n == 0 {
		return errors.WordNotFound
	}

	return nil
}
