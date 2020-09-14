package database

import (
	"database/sql"
	"github.com/pkg/errors"
)

type Migrator struct {
	db         *sql.DB
	migrations []string
}

func NewMigrator(db *sql.DB, migrations []string) Migrator {
	return Migrator{db: db, migrations: migrations}
}

func (m Migrator) Fresh() error {
	if err := m.Drop(); err != nil {
		return err
	}

	return m.Migrate()
}

func (m Migrator) Drop() error {
	rows, err := m.db.Query(`SELECT 'DROP TABLE IF EXISTS "' || tablename || '" CASCADE;' FROM pg_tables WHERE schemaname='public';`)
	if err != nil {
		return err
	}

	for rows.Next() {
		var query string
		err := rows.Scan(&query)
		if err != nil {
			return err
		}
		if _, err := m.db.Exec(query); err != nil {
			return err
		}
	}
	return nil
}

func (m Migrator) Migrate() (err error) {
	tx, err := m.db.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec("SET TRANSACTION ISOLATION LEVEL SERIALIZABLE;"); err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`CREATE TABLE IF NOT EXISTS migrations(version INTEGER NOT NULL);`)
	if err != nil {
		panic(err)
	}

	var version int
	row := tx.QueryRow(`SELECT version FROM migrations LIMIT 1;`)
	if err := row.Scan(&version); err != nil {
		if errors.Cause(err) != sql.ErrNoRows {
			tx.Rollback()
			return err
		}

		if _, err := tx.Exec(`INSERT INTO migrations(version) VALUES(0);`); err != nil {
			tx.Rollback()
			return err
		}

		version = 0
	}

	for i := version; i < len(m.migrations); i++ {
		if _, err := tx.Exec(m.migrations[i]); err != nil {
			tx.Rollback()
			return err
		}

		if _, err = tx.Exec(`UPDATE migrations SET version=$1;`, i+1); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
