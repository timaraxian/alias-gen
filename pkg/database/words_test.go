package database

import (
	"testing"
	"time"

	"github.com/timaraxian/hotel-gen/pkg/errors"
	"github.com/timaraxian/hotel-gen/pkg/helpers/crypto"
	"github.com/timaraxian/hotel-gen/pkg/helpers/validators"
)

// -----------------------------------------------------------------------------
// DBAL.WordCreate
// -----------------------------------------------------------------------------
func TestDBAL_WordCreate(t *testing.T) {
	t.Parallel()
	dbal, close := NewTestDBAL()
	defer close()

	before := time.Now().Round(time.Microsecond)
	word, err := dbal.WordCreate("Grand", "en", "adjective")
	if err != nil {
		t.Fatal(err)
	}
	after := time.Now().Add(time.Microsecond)

	if err := validators.UUID(word.WordID); err != nil {
		t.Fatal(word.WordID)
	}

	if word.Word != "Grand" {
		t.Fatal(word.Word)
	}
	if word.Language != "en" {
		t.Fatal(word.Language)
	}
	if word.Part != "adjective" {
		t.Fatal(word.Part)
	}
	if word.CreatedAt.Before(before) || word.CreatedAt.After(after) {
		t.Fatal(word.CreatedAt)
	}
	if !word.CreatedAt.Equal(word.UpdatedAt) {
		t.Fatal(word.UpdatedAt)
	}
	if word.ArchivedAt != nil {
		t.Fatal(word.ArchivedAt)
	}
}

// -----------------------------------------------------------------------------
// DBAL.WordGet
// -----------------------------------------------------------------------------
func TestDBAL_WordGet(t *testing.T) {
	t.Parallel()
	dbal, close := NewTestDBAL()
	defer close()

	word_in, err := dbal.WordCreate("Grand", "en", "adjective")
	if err != nil {
		t.Fatal(err)
	}

	word_out, err := dbal.WordGet(word_in.WordID)
	if err != nil {
		t.Fatal(err)
	}

	if word_in.Word != word_out.Word {
		t.Fatal(word_out.Word)
	}
}

// -----------------------------------------------------------------------------
// DBAL.WordSetWord
// -----------------------------------------------------------------------------
func TestDBAL_WordSetWord(t *testing.T) {
	t.Parallel()
	dbal, close := NewTestDBAL()
	defer close()

	word_in, err := dbal.WordCreate("Grand", "en", "adjective")
	if err != nil {
		t.Fatal(err)
	}

	err = dbal.WordSetWord(word_in.WordID, "Big")
	if err != nil {
		t.Fatal(err)
	}

	word_out, err := dbal.WordGet(word_in.WordID)
	if err != nil {
		t.Fatal(err)
	}

	if word_out.Word != "Big" {
		t.Fatal(word_out.Word)
	}

	if word_in.UpdatedAt.Equal(word_out.UpdatedAt) {
		t.Fatal(word_out.UpdatedAt)
	}
}

func TestDBAL_WordSetWord_Duplicate(t *testing.T) {
	t.Parallel()
	dbal, close := NewTestDBAL()
	defer close()

	word_in, err := dbal.WordCreate("Grand", "en", "adjective")
	if err != nil {
		t.Fatal(err)
	}
	_, err = dbal.WordCreate("Pink", "en", "adjective")
	if err != nil {
		t.Fatal(err)
	}

	err = dbal.WordSetWord(word_in.WordID, "Pink")
	if err != errors.WordDuplicate {
		t.Fatal(err)
	}
}

func TestDBAL_WordSetWord_WordNotFound_validUUID(t *testing.T) {
	t.Parallel()
	dbal, close := NewTestDBAL()
	defer close()

	err := dbal.WordSetWord(crypto.NewUUID(), "Grand")
	if err != errors.WordNotFound {
		t.Fatal(err)
	}
}

func TestDBAL_WordSetWord_WordNotFound_invalidUUID(t *testing.T) {
	t.Parallel()
	dbal, close := NewTestDBAL()
	defer close()

	err := dbal.WordSetWord("invalidUUID", "Grand")
	if err != errors.WordNotFound {
		t.Fatal(err)
	}
}

// -----------------------------------------------------------------------------
// DBAL.WordSetLanguage
// -----------------------------------------------------------------------------
func TestDBAL_WordSetLanguage(t *testing.T) {
	t.Parallel()
	dbal, close := NewTestDBAL()
	defer close()

	word_in, err := dbal.WordCreate("Grand", "en", "adjective")
	if err != nil {
		t.Fatal(err)
	}

	err = dbal.WordSetLanguage(word_in.WordID, "fr")
	if err != nil {
		t.Fatal(err)
	}

	word_out, err := dbal.WordGet(word_in.WordID)
	if err != nil {
		t.Fatal(err)
	}

	if word_out.Language != "fr" {
		t.Fatal(word_out.Language)
	}
	if word_in.UpdatedAt.Equal(word_out.UpdatedAt) {
		t.Fatal(word_out.UpdatedAt)
	}
}

func TestDBAL_WordSetLanguage_Duplicate(t *testing.T) {
	t.Parallel()
	dbal, close := NewTestDBAL()
	defer close()

	word_in, err := dbal.WordCreate("Grand", "en", "adjective")
	if err != nil {
		t.Fatal(err)
	}
	_, err = dbal.WordCreate("Grand", "fr", "adjective")
	if err != nil {
		t.Fatal(err)
	}

	err = dbal.WordSetLanguage(word_in.WordID, "fr")
	if err != errors.WordDuplicate {
		t.Fatal(err)
	}
}

func TestDBAL_WordSetLanguage_WordNotFound_validUUID(t *testing.T) {
	t.Parallel()
	dbal, close := NewTestDBAL()
	defer close()

	err := dbal.WordSetLanguage(crypto.NewUUID(), "fr")
	if err != errors.WordNotFound {
		t.Fatal(err)
	}
}

func TestDBAL_WordSetLanguage_WordNotFound_invalidUUID(t *testing.T) {
	t.Parallel()
	dbal, close := NewTestDBAL()
	defer close()

	err := dbal.WordSetLanguage("invalidUUID", "fr")
	if err != errors.WordNotFound {
		t.Fatal(err)
	}
}

// -----------------------------------------------------------------------------
// DBAL.WordSetPart
// -----------------------------------------------------------------------------
func TestDBAL_WordSetPart(t *testing.T) {
	t.Parallel()
	dbal, close := NewTestDBAL()
	defer close()

	word_in, err := dbal.WordCreate("Grand", "en", "adjective")
	if err != nil {
		t.Fatal(err)
	}

	err = dbal.WordSetPart(word_in.WordID, "article")
	if err != nil {
		t.Fatal(err)
	}

	word_out, err := dbal.WordGet(word_in.WordID)
	if err != nil {
		t.Fatal(err)
	}

	if word_out.Part != "article" {
		t.Fatal(word_out.Part)
	}
	if word_in.UpdatedAt.Equal(word_out.UpdatedAt) {
		t.Fatal(word_out.UpdatedAt)
	}
}

func TestDBAL_WordSetPart_Duplicate(t *testing.T) {
	t.Parallel()
	dbal, close := NewTestDBAL()
	defer close()

	word_in, err := dbal.WordCreate("Grand", "en", "adjective")
	if err != nil {
		t.Fatal(err)
	}
	_, err = dbal.WordCreate("Grand", "en", "article")
	if err != nil {
		t.Fatal(err)
	}

	err = dbal.WordSetPart(word_in.WordID, "article")
	if err != errors.WordDuplicate {
		t.Fatal(err)
	}
}

func TestDBAL_WordSetPart_WordNotFound_validUUID(t *testing.T) {
	t.Parallel()
	dbal, close := NewTestDBAL()
	defer close()

	err := dbal.WordSetPart(crypto.NewUUID(), "article")
	if err != errors.WordNotFound {
		t.Fatal(err)
	}
}

func TestDBAL_WordSetPart_WordNotFound_invalidUUID(t *testing.T) {
	t.Parallel()
	dbal, close := NewTestDBAL()
	defer close()

	err := dbal.WordSetPart("invalidUUID", "article")
	if err != errors.WordNotFound {
		t.Fatal(err)
	}
}

// -----------------------------------------------------------------------------
// DBAL.WordSetArchive
// -----------------------------------------------------------------------------
func TestDBAL_WordSetArchive(t *testing.T) {
	t.Parallel()
	dbal, close := NewTestDBAL()
	defer close()

	word_in, err := dbal.WordCreate("Grand", "en", "adjective")
	if err != nil {
		t.Fatal(err)
	}

	err = dbal.WordSetArchive(word_in.WordID)

	word_out, err := dbal.WordGet(word_in.WordID)
	if err != nil {
		t.Fatal(err)
	}

	if word_out.ArchivedAt == nil {
		t.Fatal(word_out)
	}
}

// -----------------------------------------------------------------------------
// DBAL.WordSetUnarchive
// -----------------------------------------------------------------------------
func TestDBAL_WordSetUnarchive(t *testing.T) {
	t.Parallel()
	dbal, close := NewTestDBAL()
	defer close()

	word_in, err := dbal.WordCreate("Grand", "en", "adjective")
	if err != nil {
		t.Fatal(err)
	}

	err = dbal.WordSetUnArchive(word_in.WordID)

	word_out, err := dbal.WordGet(word_in.WordID)
	if err != nil {
		t.Fatal(err)
	}

	if word_out.ArchivedAt != nil {
		t.Fatal(word_out)
	}
}

// -----------------------------------------------------------------------------
// DBAL.WordList
// -----------------------------------------------------------------------------
func TestDBAL_WordList(t *testing.T) {
	t.Parallel()
	dbal, close := NewTestDBAL()
	defer close()

	w1, err := dbal.WordCreate("Grand", "en", "adjective")
	if err != nil {
		t.Fatal(err)
	}
	w2, err := dbal.WordCreate("Hotel", "en", "noun")
	if err != nil {
		t.Fatal(err)
	}
	w3, err := dbal.WordCreate("The", "en", "article")
	if err != nil {
		t.Fatal(err)
	}
	w4, err := dbal.WordCreate("Budapest", "en", "place")
	if err != nil {
		t.Fatal(err)
	}

	limit := 10
	trueVar := true
	listargs := WordListArgs{
		Limit:       &limit,
		OrderByWord: &trueVar,
	}

	results, err := dbal.WordList(listargs)
	if err != nil {
		t.Fatal(err)
	}

	if results[0].WordID != w4.WordID {
		t.Fatal(results[0])
	}
	if results[1].WordID != w1.WordID {
		t.Fatal(results[1])
	}
	if results[2].WordID != w2.WordID {
		t.Fatal(results[2])
	}
	if results[3].WordID != w3.WordID {
		t.Fatal(results[3])
	}
}
