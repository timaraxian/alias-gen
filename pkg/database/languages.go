package database

import "github.com/timaraxian/alias-gen/pkg/errors"

func (dbal *DBAL) GetDistinctLanguages() (languages []string, err error) {
	stmt := `SELECT DISTINCT
		p.language
		FROM words w
		INNER JOIN patterns p
		ON w.language = p.language
		WHERE p.archived_at IS NULL AND w.archived_at IS NULL
		ORDER BY p.language;`

	rows, err := dbal.Query(stmt)
	if err != nil {
		return languages, errors.UnexpectedError(err, "Failed getting distinct languages")
	}
	defer rows.Close()

	for rows.Next() {
		var language string
		if err := rows.Scan(
			&language,
		); err != nil {
			return languages, errors.UnexpectedError(err, "Failed scanning languages")
		}
		languages = append(languages, language)
	}

	if err := rows.Err(); err != nil {
		return languages, errors.UnexpectedError(err, "Failed iterating languages")
	}

	return languages, err
}
