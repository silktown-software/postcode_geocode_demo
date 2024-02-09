package database

import (
	"database/sql"
	"fmt"
	"strings"
)

const sqliteMaxVariables = 999

// SqlitePostCodeRepository represents an Postcode repository for SQLite
type SqlitePostCodeRepository struct {
	db *sql.DB
}

// NewSqlitePostCodeRepository creates a new PostCodeRepository that will work with SQLite.
func NewSqlitePostCodeRepository(db *sql.DB) (PostCodeRepository, error) {
	return SqlitePostCodeRepository{
		db: db,
	}, nil
}

// Migrate this creates the database tables and adds an index to the postcode column
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

// Get this retrieves a PostcodeEntity when supplied postcode that exists in the database.
// If it does not find the post code it will return the entity as nil.
// This method will return ErrEmptyPostCode if the postcode is empty.
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

// Upsert this takes a singular postcode entity and updates the row if the postcode is already present in the database
func (r SqlitePostCodeRepository) Upsert(entity PostCodeEntity) error {
	stmt, err := r.db.Prepare(`INSERT INTO postcode VALUES (?, ?, ?) ON CONFLICT (postcode) DO UPDATE SET lng=excluded.lng, lat=excluded.lat`)

	defer func(stmt *sql.Stmt) {
		err = stmt.Close()
	}(stmt)

	if err != nil {
		return err
	}

	_, err = stmt.Exec(entity.Postcode, entity.Lng, entity.Lat)

	return err
}

// UpsertMany this takes a slice of postcode entities and batch inserts them into the SQLite database.
// If there is a conflict on the UNIQUE constraint we UPDATE the row with the new lat and lng.
//
// SQLite has a max number of variable of 999 so if the length of the slice is longer than 999, therefore we want to
// batch the slices update/insert the rows
func (r SqlitePostCodeRepository) UpsertMany(postcodes []PostCodeEntity) error {
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
	sqlStmt := `INSERT INTO postcode (postcode, lng, lat) VALUES %s ON CONFLICT (postcode) DO UPDATE SET lng=excluded.lng, lat=excluded.lat`

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
