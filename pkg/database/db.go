package database

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/timaraxian/alias-gen/pkg/database/migrations"
)

type Config struct {
	DBHost     string
	DBName     string
	DBUser     string
	DBPassword string
	DBPort     string
	DBSSLMode  string
}

type DBAL struct {
	*sql.DB
}

var MigrationFiles = []string{
	migrations.CreateWordsTable,
	migrations.CreatePatternsTable,
}

func Bootstrap(config Config) (db *DBAL, err error) {
	db = &DBAL{}
	db.DB, err = sql.Open("postgres", fmt.Sprintf("host=%s dbname=%s user=%s password=%s port=%s sslmode=%s",
		config.DBHost, config.DBName, config.DBUser, config.DBPassword, config.DBPort, config.DBSSLMode,
	))
	if err != nil {
		return db, err
	}

	return db, db.Migrate()
}

func (dbal *DBAL) Migrate() error {
	return NewMigrator(dbal.DB, MigrationFiles).Migrate()
}

func (dbal *DBAL) Fresh() error {
	return NewMigrator(dbal.DB, MigrationFiles).Fresh()
}

func (dbal *DBAL) Close() error {
	return dbal.DB.Close()
}

func dbIsDuplicateErr(err error, constraint string) bool {
	err, ok := errors.Cause(err).(*pq.Error)
	if !ok {
		return false
	}
	pqErr := err.(*pq.Error)

	return pqErr.Code == "23505" && pqErr.Constraint == constraint
}

func dbIsForeignKeyErr(err error, constraint string) bool {
	err, ok := errors.Cause(err).(*pq.Error)
	if !ok {
		return false
	}
	pqErr := err.(*pq.Error)
	return pqErr.Code == "23503" && pqErr.Constraint == constraint
}

func dbIsNullConstraintErr(err error, column string) bool {
	err, ok := errors.Cause(err).(*pq.Error)
	if !ok {
		return false
	}
	pqErr := err.(*pq.Error)

	return pqErr.Code == "23502" && pqErr.Column == column
}

func (dbal *DBAL) ExecOne(query string, args ...interface{}) (result sql.Result, n int, err error) {
	return dbExecOne(dbal, query, args...)
}

type dbExec interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

func dbExecOne(db dbExec, query string, args ...interface{}) (result sql.Result, n int, err error) {
	result, err = db.Exec(query, args...)
	if err != nil {
		return result, 0, err
	}

	if n, err := result.RowsAffected(); err != nil {
		return result, 0, err
	} else if n > 1 {
		panic("update too many rows: " + string(n))
	} else {
		return result, int(n), nil
	}
}

type OrderBy struct {
	Key        string `json:"key"`
	Desc       bool   `json:"desc"`
	NullsFirst bool   `json:"nullsFirst"`
}

func dbBuildOrderBy(sorts []OrderBy, keys map[string]string) (result string, ok bool) {
	if len(sorts) == 0 {
		return "", true
	}

	var parts []string
	for _, s := range sorts {
		part, ok := keys[s.Key]
		if !ok {
			return "", false
		}
		if s.Desc {
			part += " DESC"
		} else {
			part += " ASC"
		}
		if s.NullsFirst {
			part += " NULLS FIRST"
		} else {
			part += " NULLS LAST"
		}
		parts = append(parts, part)
	}

	return "ORDER BY " + strings.Join(parts, ", "), true
}

func dbTX(db *sql.DB, txFunc func(*sql.Tx) error) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()
	err = txFunc(tx)
	return err
}
