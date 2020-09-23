package database

import "testing"

func TestDBAL_GetDistinctLanguage(t *testing.T) {
	t.Parallel()
	dbal, close := NewTestDBAL()
	defer close()

	_, err := dbal.WordCreate("Grand", "en", "adjective")
	if err != nil {
		t.Fatal(err)
	}
	_, err = dbal.WordCreate("Ookii", "jp", "adjective")
	if err != nil {
		t.Fatal(err)
	}
	_, err = dbal.WordCreate("Le", "fr", "article")
	if err != nil {
		t.Fatal(err)
	}
	_, err = dbal.PatternCreate("adjective,place", "en")
	if err != nil {
		t.Fatal(err)
	}
	_, err = dbal.PatternCreate("adjective,place", "fr")
	if err != nil {
		t.Fatal(err)
	}

	languages, err := dbal.GetDistinctLanguages()
	if err != nil {
		t.Fatal(err)
	}

	if len(languages) != 2 {
		t.Fatal(languages)
	}
	if languages[0] != "en" {
		t.Fatal(languages)
	}
	if languages[1] != "fr" {
		t.Fatal(languages)
	}

}
