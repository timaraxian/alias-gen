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

CONSTRAINT words_language_part_word UNIQUE (language, part, word)
);
`

// language=SQL
const CreatePatternsTable = `
CREATE TABLE patterns (
pattern_id UUID PRIMARY KEY,
pattern TEXT NOT NULL,
language TEXT NOT NULL,
created_at TIMESTAMPTZ NOT NULL,
updated_at TIMESTAMPTZ NOT NULL,
archived_at TIMESTAMPTZ,

CONSTRAINT patterns_pattern_language UNIQUE (pattern, language)
);
`
