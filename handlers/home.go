package handlers

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"mapboxdemo/database"
)

// HomeHandler this returns a http.HandlerFunc that will render the homepage
func HomeHandler(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var b bytes.Buffer

		if err := tmpl.ExecuteTemplate(&b, "index.gohtml", nil); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		if _, err := w.Write(b.Bytes()); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

// GeoCodeHandler this returns a handler which will return the postcode location data in JSON format
func GeoCodeHandler(postCodeRepo database.PostCodeRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		if !query.Has("postcode") {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		postcode := query.Get("postcode")

		if postcode == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		entity, err := postCodeRepo.Get(postcode)

		if err != nil {
			log.Printf("internal server error: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if entity == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		b, err := json.Marshal(entity)

		if err != nil {
			log.Printf("could not marshal json: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err = w.Write(b)

		if err != nil {
			log.Printf("could not write body: %s", err)
		}
	}
}
