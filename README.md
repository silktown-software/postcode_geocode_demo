# Postcode Geocode Demo

A simple go demo web-application for to look up a UK postcode from a database and mark the location on the map.

# Dependencies

- `go1.21.6`
- `make`
- `staticcheck`

## Retrieving the data

The importer assumes that the postcode data is in CSV files. 
The CSV files will have 3 columns, and they should **not** have a header.

i.e.

| Postcode  | Lat  | Lng |
|-----------|------|-----|
| ZZ99 9ZZ  | 0    | 0   |

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

You can verify if the importer has worked correctly by comparing the total count of lines in the CSV file with the total row count in the `postcode` table: 

```bash
$ find . -name '*.csv' -type f -exec cat {} + | wc -l
1736857
```

If you connect to the sqlite database:

```bash
$ sqlite3 postcodes.db 
SQLite version 3.40.1 2022-12-28 14:03:47
Enter ".help" for usage hints.
sqlite> select count(*) from postcode;
1736857
```

## Running the demo web-application

```bash
$ go run cmd/demo/main.go
```

By default, it will run on PORT 5000. You can specify an alternative port via an environment variable.

It will expect the `postcodes.db` file to be in the root of the repository if running it via `go run cmd/demo/main.go`.

## TODO:

* Split the importer logic into separate package.
* 
* Change Lat/Lng column order, so it is Lng/Lat. Lng/Lat is the preferred order.
* Write batch insert test.





