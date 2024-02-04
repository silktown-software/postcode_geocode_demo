# Postcode Geocode Demo

A simple go demo web-application for to lookup a UK postcode from a database and mark the location on the map.

# Dependencies

- `go1.21.6`
- `make`
- `staticcheck`

## Retrieving the data

The importer assumes that the postcode data is in CSV files. 
The CSV files will have 3 columns, and they should **not** have a header.

i.e.

|Postcode| Lat |Lng
|-----|-----|----|
|ZZ99 9ZZ| 0   | 0

This data can be obtained by using the postcode coordinate converter for the Ordnance Survey Codepoint Open dataset:

https://github.com/silktown-software/postcode-coordinate-converter

## Importer Usage

First run the importer to build the postcode database:

```bash
$ make build-importer
$ ./dist/importer/postcode-importer -h
Usage of ./dist/importer/postcode-importer:
  -directory string
        the directory that contains the CSV files to import
```
e.g.

If you run the importer in the root of the source directory:

```bash
$ ./dist/importer/postcode-importer -directory=<directory with CSV files>
```

## Running the demo web-application

```bash
$ go run cmd/demo/main.go
```

By default, it will run on PORT 5000. You can specify an alternative port via an environment variable.

It will expect the `postcodes.db` file to be in the root of the repository if running it via `go run cmd/demo/main.go`

## TODO:

* Split importer into separate package
* Handle 404/500 errors properly on the frontend
* Change Lat/Lng column order so it is Lng/Lat





