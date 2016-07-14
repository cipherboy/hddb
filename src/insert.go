package src


import (
"strings"
"os"
"bufio"
"path"
"log"
"fmt"
)

func ImportFile(file string, date string, bp string) {
	var basepath string
	var fields []string
	var lines [][]string

	basepath, fields, lines = ParseFile(file, date, bp)

	tx, err := DB.Begin()

	fmt.Println("Fields: ", fields)
	fmt.Println("Basepath: ", basepath)

	if err != nil {
		log.Fatal("Error beginning DB transaction: ", err)
	}

	for _, value := range lines {
		if len(value) != len(fields) {
			for k, v := range value {
				fmt.Println(k, ":", v)
			}
			log.Fatal("Error! Mismatch line length: ", len(value), " for ", len(fields), " fields")
		}

		for i, _ := range value {
			if fields[i] == "filename" {
				value[i] = strings.TrimSpace(path.Join(basepath, value[i]))
			}
		}
		insertDB(tx, fields, value)
	}

	tx.Commit()
}

func ParseFile(file string, date string, bp string) (string, []string, [][]string) {
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

	return basepath, fields, lines
}

