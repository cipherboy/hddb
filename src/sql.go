package src

import (
"database/sql"
"strings"
"log"
)

func InsertDB(tx *sql.Tx, fields []string, file []string) {
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

