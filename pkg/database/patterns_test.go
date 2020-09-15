package database

import (
	"testing"
	"time"

	"github.com/timaraxian/hotel-gen/pkg/errors"
	"github.com/timaraxian/hotel-gen/pkg/helpers/crypto"
	"github.com/timaraxian/hotel-gen/pkg/helpers/validators"
)

// -----------------------------------------------------------------------------
// DBAL.PatternCreate
// -----------------------------------------------------------------------------
func TestDBAL_PatternCreate(t *testing.T) {
	t.Parallel()
	dbal, close := NewTestDBAL()
	defer close()

	before := time.Now().Round(time.Microsecond)
	pattern, err := dbal.PatternCreate("article,adjective,noun", "en")
	if err != nil {
		t.Fatal(err)
	}
	after := time.Now().Add(time.Microsecond)

	if err := validators.UUID(pattern.PatternID); err != nil {
		t.Fatal(pattern.PatternID)
	}

	if pattern.Pattern != "article,adjective,noun" {
		t.Fatal(pattern.Pattern)
	}
	if pattern.Language != "en" {
		t.Fatal(pattern.Language)
	}
	if pattern.CreatedAt.Before(before) || pattern.CreatedAt.After(after) {
		t.Fatal(pattern.CreatedAt)
	}
	if !pattern.CreatedAt.Equal(pattern.UpdatedAt) {
		t.Fatal(pattern.UpdatedAt)
	}
	if pattern.ArchivedAt != nil {
		t.Fatal(pattern.ArchivedAt)
	}
}

// -----------------------------------------------------------------------------
// DBAL.PatternGet
// -----------------------------------------------------------------------------
func TestDBAL_PatternGet(t *testing.T) {
	t.Parallel()
	dbal, close := NewTestDBAL()
	defer close()

	pattern_in, err := dbal.PatternCreate("article,adjective,noun", "en")
	if err != nil {
		t.Fatal(err)
	}

	pattern_out, err := dbal.PatternGet(pattern_in.PatternID)
	if err != nil {
		t.Fatal(err)
	}

	if pattern_in.Pattern != pattern_out.Pattern {
		t.Fatal(pattern_out.Pattern)
	}
}

// -----------------------------------------------------------------------------
// DBAL.PatternSetPattern
// -----------------------------------------------------------------------------
func TestDBAL_PatternSetPattern(t *testing.T) {
	t.Parallel()
	dbal, close := NewTestDBAL()
	defer close()

	pattern_in, err := dbal.PatternCreate("article,adjective,noun", "en")
	if err != nil {
		t.Fatal(err)
	}

	err = dbal.PatternSetPattern(pattern_in.PatternID, "adjective,adjective")
	if err != nil {
		t.Fatal(err)
	}

	pattern_out, err := dbal.PatternGet(pattern_in.PatternID)
	if err != nil {
		t.Fatal(err)
	}

	if pattern_out.Pattern != "adjective,adjective" {
		t.Fatal(pattern_out.Pattern)
	}

	if pattern_in.UpdatedAt.Equal(pattern_out.UpdatedAt) {
		t.Fatal(pattern_out.UpdatedAt)
	}
}

func TestDBAL_PatternSetPattern_Duplicate(t *testing.T) {
	t.Parallel()
	dbal, close := NewTestDBAL()
	defer close()

	pattern_in, err := dbal.PatternCreate("article,adjective,noun", "en")
	if err != nil {
		t.Fatal(err)
	}
	_, err = dbal.PatternCreate("adjective,adjective", "en")
	if err != nil {
		t.Fatal(err)
	}

	err = dbal.PatternSetPattern(pattern_in.PatternID, "adjective,adjective")
	if err != errors.PatternDuplicate {
		t.Fatal(err)
	}
}

func TestDBAL_PatternSetPattern_PatternNotFound_validUUID(t *testing.T) {
	t.Parallel()
	dbal, close := NewTestDBAL()
	defer close()

	err := dbal.PatternSetPattern(crypto.NewUUID(), "article,adjective,noun")
	if err != errors.PatternNotFound {
		t.Fatal(err)
	}
}

func TestDBAL_PatternSetPattern_PatternNotFound_invalidUUID(t *testing.T) {
	t.Parallel()
	dbal, close := NewTestDBAL()
	defer close()

	err := dbal.PatternSetPattern("invalidUUID", "article,adjective,noun")
	if err != errors.PatternNotFound {
		t.Fatal(err)
	}
}

// -----------------------------------------------------------------------------
// DBAL.PatternSetLanguage
// -----------------------------------------------------------------------------
func TestDBAL_PatternSetLanguage(t *testing.T) {
	t.Parallel()
	dbal, close := NewTestDBAL()
	defer close()

	pattern_in, err := dbal.PatternCreate("article,adjective,noun", "en")
	if err != nil {
		t.Fatal(err)
	}

	err = dbal.PatternSetLanguage(pattern_in.PatternID, "fr")
	if err != nil {
		t.Fatal(err)
	}

	pattern_out, err := dbal.PatternGet(pattern_in.PatternID)
	if err != nil {
		t.Fatal(err)
	}

	if pattern_out.Language != "fr" {
		t.Fatal(pattern_out.Language)
	}
	if pattern_in.UpdatedAt.Equal(pattern_out.UpdatedAt) {
		t.Fatal(pattern_out.UpdatedAt)
	}
}

func TestDBAL_PatternSetLanguage_Duplicate(t *testing.T) {
	t.Parallel()
	dbal, close := NewTestDBAL()
	defer close()

	pattern_in, err := dbal.PatternCreate("article,adjective,noun", "en")
	if err != nil {
		t.Fatal(err)
	}
	_, err = dbal.PatternCreate("article,adjective,noun", "fr")
	if err != nil {
		t.Fatal(err)
	}

	err = dbal.PatternSetLanguage(pattern_in.PatternID, "fr")
	if err != errors.PatternDuplicate {
		t.Fatal(err)
	}
}

func TestDBAL_PatternSetLanguage_PatternNotFound_validUUID(t *testing.T) {
	t.Parallel()
	dbal, close := NewTestDBAL()
	defer close()

	err := dbal.PatternSetLanguage(crypto.NewUUID(), "fr")
	if err != errors.PatternNotFound {
		t.Fatal(err)
	}
}

func TestDBAL_PatternSetLanguage_PatternNotFound_invalidUUID(t *testing.T) {
	t.Parallel()
	dbal, close := NewTestDBAL()
	defer close()

	err := dbal.PatternSetLanguage("invalidUUID", "fr")
	if err != errors.PatternNotFound {
		t.Fatal(err)
	}
}

// -----------------------------------------------------------------------------
// DBAL.PatternSetArchive
// -----------------------------------------------------------------------------
func TestDBAL_PatternSetArchive(t *testing.T) {
	t.Parallel()
	dbal, close := NewTestDBAL()
	defer close()

	pattern_in, err := dbal.PatternCreate("article,adjective,noun", "en")
	if err != nil {
		t.Fatal(err)
	}

	err = dbal.PatternSetArchive(pattern_in.PatternID)

	pattern_out, err := dbal.PatternGet(pattern_in.PatternID)
	if err != nil {
		t.Fatal(err)
	}

	if pattern_out.ArchivedAt == nil {
		t.Fatal(pattern_out)
	}
}

// -----------------------------------------------------------------------------
// DBAL.PatternSetUnarchive
// -----------------------------------------------------------------------------
func TestDBAL_PatternSetUnarchive(t *testing.T) {
	t.Parallel()
	dbal, close := NewTestDBAL()
	defer close()

	pattern_in, err := dbal.PatternCreate("article,adjective,noun", "en")
	if err != nil {
		t.Fatal(err)
	}

	err = dbal.PatternSetUnArchive(pattern_in.PatternID)

	pattern_out, err := dbal.PatternGet(pattern_in.PatternID)
	if err != nil {
		t.Fatal(err)
	}

	if pattern_out.ArchivedAt != nil {
		t.Fatal(pattern_out)
	}
}

// -----------------------------------------------------------------------------
// DBAL.PatternList
// -----------------------------------------------------------------------------
func TestDBAL_PatternList(t *testing.T) {
	t.Parallel()
	dbal, close := NewTestDBAL()
	defer close()

	p1, err := dbal.PatternCreate("article,adjective,place,noun", "en")
	if err != nil {
		t.Fatal(err)
	}
	p2, err := dbal.PatternCreate("adjective,noun", "en")
	if err != nil {
		t.Fatal(err)
	}
	p3, err := dbal.PatternCreate("place,article,noun", "en")
	if err != nil {
		t.Fatal(err)
	}
	p4, err := dbal.PatternCreate("adjective,noun", "fr")
	if err != nil {
		t.Fatal(err)
	}

	limit := 10
	trueVar := true
	listargs := PatternListArgs{
		Limit:           &limit,
		OrderByPattern:  &trueVar,
		OrderByLanguage: &trueVar,
	}

	results, err := dbal.PatternList(listargs)
	if err != nil {
		t.Fatal(err)
	}

	if results[0].PatternID != p2.PatternID {
		t.Fatal(results[0])
	}
	if results[1].PatternID != p4.PatternID {
		t.Fatal(results[1])
	}
	if results[2].PatternID != p1.PatternID {
		t.Fatal(results[2])
	}
	if results[3].PatternID != p3.PatternID {
		t.Fatal(results[3])
	}
}
