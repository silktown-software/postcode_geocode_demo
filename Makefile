build-demo: dep test static-check
	mkdir -p dist/demo
	go build -o dist/demo/demo cmd/demo/main.go
	cp -R static dist/demo

build-importer: dep test static-check
	mkdir -p dist/importer
	go build -o dist/importer/postcode-importer cmd/importer/main.go

run-demo: dep
	go run cmd/demo/main.go

clean:
	rm -rf dist/
	rm coverage.out

test:
	go test ./...

test_coverage:
	go test ./... -coverprofile=coverage.out

dep:
	go mod tidy

static-check:
	staticcheck ./...
