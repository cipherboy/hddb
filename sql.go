package main

import (
	"database/sql"
	"log"
	"strings"
)

func insertDB(tx *sql.Tx, fields []string, file []string) {
	values := make([]interface{}, len(file))
	for i, v := range file {
		values[i] = v
	}

	var query string = "INSERT INTO files (" + strings.Join(fields, ",") + ") VALUES (" + strings.Repeat("?,", len(fields)-1) + "?)"

	_, err := tx.Exec(query, values...)

	if err != nil {
		log.Fatal("Error inserting into DB: ", err, " ::::: ", query)
	}
}

func ignoreDB(tx *sql.Tx, file string) {
	var query string = "INSERT INTO ignored (filename, ignored) VALUES (?, ?)"

	_, err := tx.Exec(query, file, true)

	if err != nil {
		log.Fatal("Error inserting into DB: ", err, " ::::: ", query)
	}
}
