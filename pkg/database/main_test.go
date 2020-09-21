package database

import (
	"log"
	"os"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/timaraxian/alias-gen/pkg/helpers/testdb"
)

var tdb *testdb.Manager

func TestMain(m *testing.M) {
	var config struct{ DB Config }
	_, err := toml.DecodeFile("../../config.test.toml", &config)
	if err != nil {
		log.Fatal(err, "database test config not found", "../..config.test.toml")
	}

	tdb, err = testdb.NewManager(
		"db_testing", 3,
		config.DB.DBHost,
		config.DB.DBUser,
		config.DB.DBPassword,
		config.DB.DBPort,
		config.DB.DBSSLMode,
	)
	if err != nil {
		panic(err)
	}

	var status int

	defer func() {
		recover()
		tdb.TearDown()
		os.Exit(status)
	}()

	status = m.Run()
}

func NewTestDBAL() (dbal *DBAL, close func()) {
	conn, close, err := tdb.NewConn()
	if err != nil {
		panic(err)
	}

	dbal = &DBAL{DB: conn}
	err = NewMigrator(conn, MigrationFiles).Fresh()
	if err != nil {
		panic(err)
	}
	return dbal, close
}
