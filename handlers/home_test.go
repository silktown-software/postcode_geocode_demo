package handlers_test

import (
	"errors"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"mapboxdemo"
	"mapboxdemo/database"
	"mapboxdemo/handlers"
)

type stubRepo struct {
	entities []database.PostCodeEntity
}

func newStubRepo() stubRepo {
	repo := stubRepo{}

	repo.entities = append(repo.entities, database.PostCodeEntity{
		Postcode: "AB10 1AB", Lng: 0, Lat: 0,
	})

	return repo
}

func (r stubRepo) Migrate() error {
	return nil
}

func (r stubRepo) Get(postcode string) (*database.PostCodeEntity, error) {
	if postcode == "ERROR" {
		return nil, errors.New("a generic error")
	}

	for _, entity := range r.entities {
		if entity.Postcode == postcode {
			return &entity, nil
		}
	}

	return nil, nil
}

func (r stubRepo) Upsert(entity database.PostCodeEntity) error {
	return nil
}

func (r stubRepo) UpsertMany(postcodes []database.PostCodeEntity) error {
	return nil
}

func TestHomeHandler(t *testing.T) {
	tmpl, err := template.ParseFS(mapboxdemo.TemplateFS, "templates/*.gohtml")

	if err != nil {
		log.Fatalf("could not parse templates: %s", err)
	}

	handlerFunc := handlers.HomeHandler(tmpl)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handlerFunc(w, r)

	rs := w.Result()

	assert.Equal(t, rs.StatusCode, http.StatusOK)

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			assert.Fail(t, "could not close body reader")
		}
	}(rs.Body)

	html, err := io.ReadAll(rs.Body)

	if err != nil {
		assert.Fail(t, "could not read response body: %+v", err)
	}

	strings.Contains(string(html[:]), "<!DOCTYPE html>")
}

func TestGeoCodeHandler(t *testing.T) {
	cases := []struct {
		name       string
		postcode   string
		statusCode int
		body       string
	}{
		{
			name:       "existing postcode should return 200 with json",
			statusCode: http.StatusOK,
			postcode:   "AB10 1AB",
			body:       "{\"postcode\":\"AB10 1AB\",\"lng\":0,\"lat\":0}",
		},
		{
			name:       "empty postcode should return a 400 badrequest",
			statusCode: http.StatusBadRequest,
			postcode:   "",
			body:       "",
		},
		{
			name:       "postcode not present should return 404",
			statusCode: 404,
			postcode:   "ZZ99 ZZZ",
			body:       "",
		},
		{
			name:       "should return 500 when there is an error from the repo",
			statusCode: 500,
			postcode:   "ERROR", // this is a hardcoded error string in the stup repo
			body:       "",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/geocode", nil)

			q := r.URL.Query()
			q.Add("postcode", c.postcode)

			r.URL.RawQuery = q.Encode()

			w := httptest.NewRecorder()

			handlerFunc := handlers.GeoCodeHandler(newStubRepo())

			handlerFunc(w, r)

			res := w.Result()

			assert.Equal(t, c.statusCode, res.StatusCode)

			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					assert.NoError(t, err)
				}
			}(res.Body)

			bytes, err := io.ReadAll(res.Body)

			assert.NoError(t, err)
			assert.Equal(t, c.body, string(bytes[:]))
		})
	}
}
