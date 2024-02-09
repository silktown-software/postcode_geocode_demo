package database_test

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"mapboxdemo/database"

	_ "github.com/mattn/go-sqlite3"
)

var repo database.PostCodeRepository
var db *sql.DB

func TestMain(m *testing.M) {
	var err error

	// open and in memory database https://www.sqlite.org/inmemorydb.html
	db, err = sql.Open("sqlite3", ":memory:")

	if err != nil {
		log.Fatalf("could not open in memory database: %s", err)
	}

	repo, err = database.NewSqlitePostCodeRepository(db)

	if err != nil {
		log.Fatal(err)
	}

	code := m.Run()

	os.Exit(code)
}

func TestRepository_Migrate(t *testing.T) {
	err := repo.Migrate()

	assert.Nil(t, err)

	rows, err := db.Query(`SELECT name FROM sqlite_master WHERE type='table' AND name='postcode'`)

	assert.Nil(t, err)

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			t.Fatalf("could not close rows")
		}
	}(rows)

	next := rows.Next()

	if !next {
		assert.FailNow(t, "should be at least one table")
	}

	var tableName string

	err = rows.Scan(&tableName)

	assert.Nil(t, err)

	assert.Equal(t, tableName, "postcode")
}

func cleanUp(t *testing.T) {
	_, err := db.Exec(`DELETE FROM postcode`)

	if err != nil {
		assert.FailNow(t, "could not truncate postcode table")
	}
}

func TestSqlitePostCodeRepository_Get(t *testing.T) {
	err := repo.Migrate()

	if err != nil {
		assert.FailNow(t, "could not migrate db", err)
	}

	stmt, err := db.Prepare("INSERT INTO postcode VALUES (?, ?, ?)")

	if err != nil {
		assert.FailNow(t, "could not prepare insert statement", err)
	}

	_, err = stmt.Exec("AB10 1AB", 0, 0)

	if err != nil {
		assert.FailNow(t, "could not execute insert: %s", err)
	}

	cases := []struct {
		name     string
		postcode string
		entity   *database.PostCodeEntity
		err      error
	}{
		{
			name:     "empty postcode should return nil and error",
			postcode: "",
			entity:   nil,
			err:      database.ErrEmptyPostcode,
		},
		{
			name:     "existing postcode should return correct result",
			postcode: "AB10 1AB",
			entity: &database.PostCodeEntity{
				Postcode: "AB10 1AB",
				Lat:      0.00,
				Lng:      0.00,
			},
			err: nil,
		},
		{
			name:     "exiting postcode lower case should return correct result",
			postcode: "ab10 1ab",
			entity: &database.PostCodeEntity{
				Postcode: "AB10 1AB",
				Lat:      0.00,
				Lng:      0.00,
			},
			err: nil,
		},
		{
			name:     "existing postcode without spaces should return correct result",
			postcode: "AB101AB",
			entity: &database.PostCodeEntity{
				Postcode: "AB10 1AB",
				Lat:      0.00,
				Lng:      0.00,
			},
			err: nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			entity, err := repo.Get(c.postcode)

			assert.ErrorIs(t, err, c.err)
			assert.Equal(t, entity, c.entity)
		})
	}

	cleanUp(t)
}

func TestSqlitePostCodeRepository_Insert(t *testing.T) {
	err := repo.Migrate()

	if err != nil {
		assert.FailNow(t, "could not migrate db", err)
	}

	stmt, err := db.Prepare("INSERT INTO postcode VALUES (?, ?, ?)")

	if err != nil {
		assert.FailNow(t, "could not prepare insert statement", err)
	}

	_, err = stmt.Exec("AB10 1AB", 0, 0)

	if err != nil {
		assert.FailNow(t, "could not execute insert: %s", err)
	}

	cases := []struct {
		name   string
		entity database.PostCodeEntity
		error  error
	}{
		{
			name: "insert new record into the database",
			entity: database.PostCodeEntity{
				Postcode: "ZZ99 9ZZ",
				Lat:      0,
				Lng:      0,
			},
			error: nil,
		},
		{
			name: "insert existing postcode record should no return error",
			entity: database.PostCodeEntity{
				Postcode: "AB10 1AB",
				Lat:      0.00,
				Lng:      0.00,
			},
			error: nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := repo.Upsert(c.entity)

			assert.ErrorIs(t, c.error, err)
		})
	}

	cleanUp(t)
}
