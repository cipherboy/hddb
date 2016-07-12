package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"path"
	"time"
	"strings"
)

import _ "github.com/mattn/go-sqlite3"

var DB *sql.DB

func connectDB(db string) {
	var err error

	if db == "~/.hddb/main.db" {
		usr, _ := user.Current()
		dir := usr.HomeDir
		db = path.Join(dir, "/.hddb/")
		os.Mkdir(db, 0750)
		db = path.Join(db, "main.db")
	}

	DB, err = sql.Open("sqlite3", db)
	if err != nil {
		log.Fatal("Unable to load database: ", err)
	}

	DB.Exec("CREATE TABLE IF NOT EXISTS files (size INTEGER, md5 STRING, sha1 STRING, sha256 STRING, tiger STRING, whirlpool STRING, filename STRING, scandate DATE);")
	DB.Exec("CREATE TABLE IF NOT EXISTS ignored (filename STRING, ignored BOOLEAN);")
}

func parseFile(file string, date string, bp string) {
	var line int = 0
	var fields []string
	var basepath string = bp
	var lines [][]string

	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if line == 0 {
		} else if line == 1 {
			fields = strings.Split(scanner.Text()[5:], ",")
			fields = append(fields, "scandate")
		} else if line == 2 {
			if bp == "" {
				basepath = scanner.Text()[16:]
			} else {
				fmt.Println("bp: " + bp)
			}
		} else if line > 4 {
			lines = append(lines, append(strings.SplitN(scanner.Text(), ",", len(fields)-1), date))
		}

		line += 1
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Fields: ", fields)
	fmt.Println("Basepath: ", basepath)

	tx, err := DB.Begin()

	if err != nil {
		log.Fatal("Error beginning DB transaction: ", err)
	}

	for _, value := range lines {
		if len(value) != len(fields) {
			for k, v := range(value) {
				fmt.Println(k, ":", v)
			}
			log.Fatal("Error! Mismatch line length: ", len(value), " for ", len(fields), " fields")
		}

		for i, _ := range value {
			if fields[i] == "filename" {
				value[i] = path.Join(basepath, value[i])
			}
		}
		insertDB(tx, fields, value)
	}

	tx.Commit()
}

func insertDB(tx *sql.Tx, fields []string, file []string) {
	values := make([]interface{}, len(file))
	for i, v := range file {
	    values[i] = v
	}

	var query string = "INSERT INTO files (" + strings.Join(fields, ",") + ") VALUES (" +  strings.Repeat("?,", len(fields) - 1) + "?)"

	_, err := tx.Exec(query, values...)

	if err != nil {
		log.Fatal("Error inserting into DB: ", err, " ::::: ", query)
	}
}

func main() {
	var db string
	var file string
	var date string
	var basepath string
	var check bool
	var ignore string

	flag.StringVar(&db, "database", "~/.hddb/main.db", "path to hddb database")
	flag.StringVar(&file, "import", "", "hashdeep file to import into the database")
	flag.StringVar(&basepath, "basepath", "", "override file's basepath")
	flag.StringVar(&date, "date", "", "date of hashdeep scan")
	flag.BoolVar(&check, "check", false, "perform check for changes on known files against latest changeset")
	flag.StringVar(&ignore, "ignore", "", "exclude file from checks")

	flag.Parse()

	if file == "" && check == false && ignore == "" {
		flag.Usage()
		return
	}

	connectDB(db)

	if file != "" {
		if _, err := os.Stat(file); err == nil {
			fmt.Println("Parsing file... `" + basepath + "`")

			if (date == "") {
				t := time.Now()
				date = fmt.Sprintf("%d-%d-%d", t.Month(), t.Day(), t.Year())
			}

			fmt.Println("date:", date)

			parseFile(file, date, basepath)
		}
	}

	DB.Close()
}
