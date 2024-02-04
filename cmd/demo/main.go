package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"mapboxdemo"
	"mapboxdemo/database"
	"mapboxdemo/handlers"

	_ "github.com/mattn/go-sqlite3"
)

const dbname = "postcodes.db"

func main() {
	var port int

	flag.IntVar(&port, "p", 5000, "Specify a port which the web server runs on. By default this is port 5000")
	flag.Parse()

	tmpl, err := template.ParseFS(mapboxdemo.TemplateFS, "templates/*.gohtml")

	if err != nil {
		log.Fatalf("could not parse templates: %s", err)
	}

	db, err := sql.Open("sqlite3", dbname)

	if err != nil {
		log.Fatalf("could not open sqlite database %s: %s", dbname, err)
	}

	postCodeRepository, err := database.NewSqlitePostCodeRepository(db)

	if err != nil {
		log.Fatal(err)
	}

	fs := http.FileServer(http.Dir("./static"))

	// setup static routs for CSS, JS and Assets
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// add basic routes
	http.HandleFunc("/", handlers.HomeHandler(tmpl))
	http.HandleFunc("/geocode", handlers.GeoCodeHandler(postCodeRepository))

	log.Printf("starting server on port: %d", port)

	addr := fmt.Sprintf(":%d", port)

	err = http.ListenAndServe(addr, nil)

	if errors.Is(err, http.ErrServerClosed) {
		log.Printf("server closed")
	} else if err != nil {
		log.Fatalf("error starting server: %s", err)
	}
}
