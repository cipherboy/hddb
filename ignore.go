package main

import (
	"fmt"
	"log"
)

func IgnoreFiles(files []string) {
	tx, err := DB.Begin()

	if err != nil {
		log.Fatal("Unable to begin transaction: ", err)
	}

	for _, file := range files {
		fmt.Println("Ignoring file: " + file)
		ignoreDB(tx, file)
	}

	tx.Commit()
}
