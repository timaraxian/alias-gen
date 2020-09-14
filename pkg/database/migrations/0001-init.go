package migrations

// language=SQL
const CreateWordsTable = `
CREATE TABLE words (
word_id    UUID PRIMARY KEY,
word       TEXT NOT NULL,
language   TEXT NOT NULL,
part       TEXT NOT NULL,
created_at  TIMESTAMPTZ NOT NULL,
updated_at  TIMESTAMPTZ NOT NULL,
archived_at TIMESTAMPTZ,

CONSTRAINT words_word_language_part UNIQUE (word, language, part)
);
`
