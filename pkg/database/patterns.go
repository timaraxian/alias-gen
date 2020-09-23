package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/timaraxian/alias-gen/pkg/errors"
	"github.com/timaraxian/alias-gen/pkg/helpers/crypto"
	"github.com/timaraxian/alias-gen/pkg/helpers/validators"
)

type Pattern struct {
	PatternID  string     `json:"patternID"`
	Pattern    string     `json:"pattern"`
	Language   string     `json:language`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
	ArchivedAt *time.Time `json:"archivedAt"`
}

func (dbal *DBAL) PatternCreate(pattern_in, language string) (pattern Pattern, err error) {
	pattern.PatternID = crypto.NewUUID()

	// todo: validate pattern, language
	pattern.Pattern = pattern_in
	pattern.Language = language

	pattern.CreatedAt = time.Now()
	pattern.UpdatedAt = pattern.CreatedAt

	stmt := `INSERT INTO patterns (
		pattern_id,
		pattern,
		language,
		created_at,
		updated_at,
		archived_at
	) VALUES ($1, $2, $3, $4, $5, NULL);`

	_, err = dbal.Exec(stmt,
		pattern.PatternID,
		pattern.Pattern,
		pattern.Language,
		pattern.CreatedAt,
		pattern.UpdatedAt,
	)

	if err == nil {
		return pattern, nil
	}

	if dbIsDuplicateErr(err, "patterns_pattern_language") {
		return pattern, errors.PatternDuplicate
	}

	return pattern, errors.UnexpectedError(err, "Failed creating pattern")
}

func (dbal *DBAL) PatternGet(patternID string) (pattern Pattern, err error) {
	if err := validators.UUID(patternID); err != nil {
		return pattern, errors.PatternNotFound
	}

	stmt := `SELECT
                pattern_id,
                pattern,
                language,
                created_at,
                updated_at,
                archived_at FROM patterns WHERE pattern_id=$1;`

	err = dbal.QueryRow(stmt, patternID).Scan(
		&pattern.PatternID,
		&pattern.Pattern,
		&pattern.Language,
		&pattern.CreatedAt,
		&pattern.UpdatedAt,
		&pattern.ArchivedAt,
	)

	if err == nil {
		return pattern, nil
	}

	if err == sql.ErrNoRows {
		return pattern, errors.PatternNotFound
	}

	return pattern, errors.UnexpectedError(err, "Failed getting pattern")
}

func (dbal DBAL) PatternSetPattern(patternID, pattern string) (err error) {
	if err := validators.UUID(patternID); err != nil {
		return errors.PatternNotFound
	}

	// todo: validate pattern

	stmt := `UPDATE patterns SET pattern=$1, updated_at=$2 WHERE pattern_id=$3 AND archived_at IS NULL;`

	_, n, err := dbal.ExecOne(stmt, pattern, time.Now(), patternID)
	if dbIsDuplicateErr(err, "patterns_pattern_language") {
		return errors.PatternDuplicate
	}
	if err != nil {
		return errors.UnexpectedError(err, "Failed to set pattern")
	} else if n == 0 {
		return errors.PatternNotFound
	}

	return nil
}

func (dbal DBAL) PatternSetLanguage(patternID, language string) (err error) {
	if err := validators.UUID(patternID); err != nil {
		return errors.PatternNotFound
	}

	// todo: validate language
	stmt := `UPDATE patterns SET language=$1, updated_at=$2 WHERE pattern_id=$3 AND archived_at IS NULL;`

	_, n, err := dbal.ExecOne(stmt, language, time.Now(), patternID)
	if dbIsDuplicateErr(err, "patterns_pattern_language") {
		return errors.PatternDuplicate
	}
	if err != nil {
		return errors.UnexpectedError(err, "Failed to set language")
	} else if n == 0 {
		return errors.PatternNotFound
	}

	return nil
}

func (dbal DBAL) PatternSetArchive(patternID string) (err error) {
	if err := validators.UUID(patternID); err != nil {
		return errors.PatternNotFound
	}

	stmt := `UPDATE patterns SET archived_at=COALESCE(archived_at, NOW()) WHERE pattern_id=$1;`

	_, n, err := dbal.ExecOne(stmt, patternID)
	if err != nil {
		return errors.UnexpectedError(err, "Failed to archive pattern")
	} else if n == 0 {
		return errors.PatternNotFound
	}

	return nil
}

func (dbal DBAL) PatternSetUnArchive(patternID string) (err error) {
	if err := validators.UUID(patternID); err != nil {
		return errors.PatternNotFound
	}

	stmt := `UPDATE patterns SET archived_at=NULL WHERE pattern_id=$1;`

	_, n, err := dbal.ExecOne(stmt, patternID)
	if err != nil {
		return errors.UnexpectedError(err, "Failed to archive pattern")
	} else if n == 0 {
		return errors.PatternNotFound
	}

	return nil
}

type PatternListArgs struct {
	Limit            *int
	Offset           *int
	OrderByPattern   *bool
	DescPattern      *bool
	OrderByLanguage  *bool
	DescLanguage     *bool
	OrderByUpdatedAt *bool
	DescUpdatedAt    *bool
	OrderByCreatedAt *bool
	DescCreatedAt    *bool
	ShowArchived     *bool
}

func (dbal DBAL) PatternList(listArgs PatternListArgs) (patterns []Pattern, err error) {
	// ------ build statement
	stmt := `SELECT 
		pattern_id,
		pattern,
		language,
		created_at,
		updated_at,
		archived_at FROM patterns %s %s %s %s;`

	// %s(1) show archived or not
	showArchived := ""
	if listArgs.ShowArchived != nil && !*listArgs.ShowArchived {
		showArchived = "WHERE archived_at IS NULL"
	}

	// %s(2) orderbys
	orderBy := "ORDER BY "
	if listArgs.OrderByPattern != nil && *listArgs.OrderByPattern {
		orderBy += "pattern"

		if listArgs.DescPattern != nil && *listArgs.DescPattern {
			orderBy += " desc, "
		} else {
			orderBy += ", "
		}
	}
	if listArgs.OrderByLanguage != nil && *listArgs.OrderByLanguage {
		orderBy += "language"

		if listArgs.DescLanguage != nil && *listArgs.DescLanguage {
			orderBy += " desc, "
		} else {
			orderBy += ", "
		}
	}
	if listArgs.OrderByUpdatedAt != nil && *listArgs.OrderByUpdatedAt {
		orderBy += "updated_at"

		if listArgs.DescUpdatedAt != nil && *listArgs.DescUpdatedAt {
			orderBy += " desc, "
		} else {
			orderBy += ", "
		}
	}
	if listArgs.OrderByCreatedAt != nil && *listArgs.OrderByCreatedAt {
		orderBy += "created_at"

		if listArgs.DescCreatedAt != nil && *listArgs.DescCreatedAt {
			orderBy += " desc, "
		} else {
			orderBy += ", "
		}
	}
	if orderBy == "ORDER BY " {
		orderBy = ""
	} else {
		orderBy = orderBy[:len(orderBy)-2]
	}

	// %s(3) limit
	limit := ""
	if listArgs.Limit != nil && *listArgs.Limit > 0 {
		limit = fmt.Sprintf("LIMIT %d", *listArgs.Limit)
	} else {
		limit = "LIMIT 50"
	}

	// %s(4) offset
	offset := ""
	if listArgs.Offset != nil && *listArgs.Offset > 0 {
		offset = fmt.Sprintf("OFFSET %d", *listArgs.Offset)
	}

	stmt = fmt.Sprintf(stmt, showArchived, orderBy, limit, offset)

	// ------- statement built

	rows, err := dbal.Query(stmt)
	if err != nil {
		return patterns, errors.UnexpectedError(err, "Failed listing patterns")
	}
	defer rows.Close()

	for rows.Next() {
		pattern := Pattern{}
		if err := rows.Scan(
			&pattern.PatternID,
			&pattern.Pattern,
			&pattern.Language,
			&pattern.CreatedAt,
			&pattern.UpdatedAt,
			&pattern.ArchivedAt,
		); err != nil {
			return patterns, errors.UnexpectedError(err, "Failed scanning patterns")
		}

		patterns = append(patterns, pattern)
	}

	if err := rows.Err(); err != nil {
		return patterns, errors.UnexpectedError(err, "Failed iterating pattern rows")
	}

	return patterns, err
}

func (dbal DBAL) PatternRandom(language string) (pattern Pattern, err error) {
	// todo: validate language

	stmt := `SELECT
                pattern_id,
                pattern,
                language,
                created_at,
                updated_at,
                archived_at FROM patterns WHERE language=$1 AND archived_at IS NULL ORDER BY RANDOM() LIMIT 1;`

	err = dbal.QueryRow(stmt, language).Scan(
		&pattern.PatternID,
		&pattern.Pattern,
		&pattern.Language,
		&pattern.CreatedAt,
		&pattern.UpdatedAt,
		&pattern.ArchivedAt,
	)

	if err == nil {
		return pattern, nil
	}

	if err == sql.ErrNoRows {
		return pattern, errors.PatternNotFound
	}

	return pattern, errors.UnexpectedError(err, "Failed getting pattern")
}
