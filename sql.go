package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/user"
	"path"
	"strings"
)

func connectDB(db string) {
	var err error

	if db == "~/.hddb/main.db" {
		usr, err := user.Current()
		if err != nil {
			log.Fatal("Unable to get current user: ", err)
		}

		dir := usr.HomeDir
		db = path.Join(dir, "/.hddb/")
		os.Mkdir(db, 0750)
		db = path.Join(db, "main.db")
	}

	fmt.Println("Database path:", db)
	DB, err = sql.Open("sqlite3", db)
	if err != nil {
		log.Fatal("Unable to load database: ", err)
	}

	DB.Exec("CREATE TABLE IF NOT EXISTS files (size INTEGER DEFAULT 0, md5 STRING DEFAULT '', sha1 STRING DEFAULT '', sha256 STRING DEFAULT '', tiger STRING DEFAULT '', whirlpool STRING DEFAULT '', filename STRING DEFAULT '', scandate DATE DEFAULT '01-02-1970');")
	DB.Exec("CREATE TABLE IF NOT EXISTS ignored (filename STRING, ignored BOOLEAN DEFAULT FALSE);")
}

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

func getAllFileNamesDB(tx *sql.Tx) []FileList {
	var result []FileList
	var query string = "SELECT DISTINCT(files.filename),ignored.ignored FROM files LEFT JOIN ignored;"

	rows, err := tx.Query(query)

	if err != nil {
		log.Fatal("Error in getting all files from DB: ", err, " ::::: ", query)
	}

	for rows.Next() {
		var item FileList
		var ignored interface{}
		err = rows.Scan(&item.FileName, &ignored)

		if ignored == nil {
			item.Ignored = false
		} else {
			item.Ignored = ignored.(bool)
		}

		if err != nil {
			log.Fatal("Error scanning row from DB: ", err, " ::::: ", query)
		}

		result = append(result, item)
	}

	return result
}

func getAllHashesDB(tx *sql.Tx, filename string) []FileHashes {
	var result []FileHashes
	var query string = "SELECT size,md5,sha1,sha256,tiger,whirlpool,filename,scandate FROM files where filename = ? ORDER BY scandate DESC"

	rows, err := tx.Query(query, filename)

	if err != nil {
		log.Fatal("Error in getting all files from DB: ", err, " ::::: ", query)
	}

	for rows.Next() {
		var item iFileHashes
		err = rows.Scan(&item.Size, &item.MD5, &item.SHA1, &item.SHA256, &item.Tiger, &item.Whirlpool, &item.FileName, &item.ScanDate)

		if err != nil {
			log.Fatal("Error scanning row from DB: ", err, " ::::: ", query)
		}

		var ritem FileHashes

		ritem.Size = item.Size
		ritem.MD5 = string(item.MD5)
		ritem.SHA1 = string(item.SHA1)
		ritem.SHA256 = string(item.SHA256)
		ritem.Tiger = string(item.Tiger)
		ritem.Whirlpool = string(item.Whirlpool)
		ritem.FileName = string(item.FileName)
		ritem.ScanDate = item.ScanDate

		result = append(result, ritem)
	}

	return result
}
