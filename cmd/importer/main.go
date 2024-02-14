package main

import (
	"database/sql"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
	"mapboxdemo/database"
)

const dbname = "postcodes.db"

func main() {
	var path string

	flag.StringVar(&path, "directory", "", "The directory that contains the CSV files to import")
	flag.Parse()

	if path == "" {
		log.Fatal("no directory specified")
	}

	dir, err := os.Stat(path)

	if err != nil {
		log.Fatalf("could not stat dir: %s", err)
	}

	if !dir.IsDir() {
		log.Fatalf("path specified is not a directory: %s", path)
	}

	pattern := filepath.Join(path, "*.csv")

	matches, err := filepath.Glob(pattern)

	if err != nil {
		log.Fatalf("could not get CSV files: %s", err)
	}

	if len(matches) == 0 {
		log.Fatalf("no csv files in the directory")
	}

	db, err := sql.Open("sqlite3", dbname)

	if err != nil {
		log.Fatalf("could not open sqlite database %s: %s", dbname, err)
	}

	repo, err := database.NewSqlitePostCodeRepository(db)

	if err != nil {
		log.Fatalf("could not initialise repository: %s", err)
	}

	if err = repo.Migrate(); err != nil {
		log.Fatalf("could not migrate database: %s", err)
	}

	for _, match := range matches {
		if err = processCSVFile(repo, match); err != nil {
			log.Fatalf("could not process CSV file: %s", err)
		}
	}
}

func processCSVFile(repo database.PostCodeRepository, match string) error {
	f, err := os.Open(match)

	defer func(f *os.File) {
		err = f.Close()
	}(f)

	if err != nil {
		return fmt.Errorf("could not open csv file: %s", err)
	}

	reader := csv.NewReader(f)

	var postcodes []database.PostCodeEntity

	for {
		rec, err := reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("error while reading csv file: %w", err)
		}

		lat, err := strconv.ParseFloat(rec[1], 64)

		if err != nil {
			return fmt.Errorf("could not convert latitude on postcode: %w", err)
		}

		lng, err := strconv.ParseFloat(rec[2], 64)

		if err != nil {
			return fmt.Errorf("could not convert longitude on record: %w", err)
		}

		ent := database.PostCodeEntity{
			Postcode: rec[0],
			Lng:      lng,
			Lat:      lat,
		}

		postcodes = append(postcodes, ent)
	}

	log.Printf("bulk inserting: %s", match)

	err = repo.UpsertMany(postcodes)

	if err != nil {
		return fmt.Errorf("error processing file: %s, error: %w", match, err)
	}

	return nil
}
