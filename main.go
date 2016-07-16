package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

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

	ConnectDB(db)
	defer DB.Close()

	if file != "" {
		if _, err := os.Stat(file); err == nil {
			fmt.Println("Parsing file... `" + basepath + "`")

			if date == "" {
				t := time.Now()
				date = fmt.Sprintf("%d-%d-%d", t.Month(), t.Day(), t.Year())
			}

			fmt.Println("date:", date)

			ImportFile(file, date, basepath)
		}
	}

	if ignore != "" {
		fmt.Println("Ignoring files...")
		IgnoreFiles([]string{ignore})
	}
}
