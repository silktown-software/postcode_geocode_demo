package database

import (
	"database/sql"
	"fmt"
	"strings"
)

const sqliteMaxVariables = 999

type SqlitePostCodeRepository struct {
	db *sql.DB
}

func NewSqlitePostCodeRepository(db *sql.DB) (PostCodeRepository, error) {
	return SqlitePostCodeRepository{
		db: db,
	}, nil
}

func (r SqlitePostCodeRepository) Migrate() error {
	const create string = `
		CREATE TABLE IF NOT EXISTS postcode (
			postcode TEXT,
			lng REAL,
			lat REAL
		);
		
		DROP INDEX IF EXISTS idx_postcode;
		
		CREATE UNIQUE INDEX idx_postcode ON postcode (postcode);
	`

	if _, err := r.db.Exec(create); err != nil {
		return err
	}

	return nil
}

func (r SqlitePostCodeRepository) Get(postcode string) (*PostCodeEntity, error) {
	if postcode == "" {
		return nil, ErrEmptyPostcode
	}

	stmt, err := r.db.Prepare("SELECT postcode, lng, lat FROM postcode WHERE replace(postcode, ' ', '') = replace(?, ' ', '') COLLATE NOCASE")

	defer func(stmt *sql.Stmt) {
		err = stmt.Close()
	}(stmt)

	if err != nil {
		return nil, err
	}

	var lat float64
	var lng float64

	rows, err := stmt.Query(postcode)

	defer func(rows *sql.Rows) {
		err = rows.Close()
	}(rows)

	if !rows.Next() {
		return nil, nil
	}

	err = rows.Scan(&postcode, &lng, &lat)

	if err != nil {
		return nil, err
	}

	return &PostCodeEntity{Postcode: postcode, Lat: lat, Lng: lng}, nil
}

func (r SqlitePostCodeRepository) Insert(entity PostCodeEntity) error {
	stmt, err := r.db.Prepare("INSERT INTO postcode VALUES (?, ?, ?)")

	defer func(stmt *sql.Stmt) {
		err = stmt.Close()
	}(stmt)

	if err != nil {
		return err
	}

	_, err = stmt.Exec(entity.Postcode, entity.Lng, entity.Lat)

	return err
}

func (r SqlitePostCodeRepository) InsertMany(postcodes []PostCodeEntity) error {
	var valueStrings []string
	var valueArgs []interface{}

	for idx, p := range postcodes {
		valueStrings = append(valueStrings, "(?, ?, ?)")
		valueArgs = append(valueArgs, p.Postcode)
		valueArgs = append(valueArgs, p.Lng)
		valueArgs = append(valueArgs, p.Lat)

		if idx == 0 {
			continue
		}

		rem := idx % sqliteMaxVariables

		if rem == 0 || idx == len(postcodes)-1 {
			err := r.batchInsert(valueStrings, valueArgs)

			if err != nil {
				return err
			}

			valueStrings = nil
			valueArgs = nil

			continue
		}
	}

	return nil
}

func (r SqlitePostCodeRepository) batchInsert(valueStrings []string, valueArgs []interface{}) error {
	sqlStmt := `INSERT INTO postcode (postcode, lng, lat) VALUES %s`

	sqlStmt = fmt.Sprintf(sqlStmt, strings.Join(valueStrings, ","))

	stmt, err := r.db.Prepare(sqlStmt)

	defer func(stmt *sql.Stmt) {
		err = stmt.Close()
	}(stmt)

	defer func(stmt *sql.Stmt) {
		err = stmt.Close()
	}(stmt)

	if err != nil {
		return fmt.Errorf("could not prepare bulk insert statement: %w", err)
	}

	_, err = stmt.Exec(valueArgs...)

	if err != nil {
		return fmt.Errorf("could not execute bulk insert statement: %w", err)
	}

	return nil
}
