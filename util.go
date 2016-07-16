package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/user"
	"path"
)

import _ "github.com/mattn/go-sqlite3"

var DB *sql.DB

func ConnectDB(db string) {
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

	DB.Exec("CREATE TABLE IF NOT EXISTS files (size INTEGER, md5 STRING, sha1 STRING, sha256 STRING, tiger STRING, whirlpool STRING, filename STRING, scandate DATE);")
	DB.Exec("CREATE TABLE IF NOT EXISTS ignored (filename STRING, ignored BOOLEAN);")
}
