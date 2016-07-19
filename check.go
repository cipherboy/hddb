package main

import (
	"fmt"
    "log"
)

func CheckAllFiles() {
	tx, err := DB.Begin()

	if err != nil {
		log.Fatal("Unable to begin transaction: ", err)
	}

	var files []FileList = getAllFileNamesDB(tx)

    fmt.Println("Checking:", len(files), "files")

    var onepercent int = len(files) / 100

	for i, file := range files {

        if i != 0 && (i % onepercent) == 0 {
            fmt.Println(".")
        }

		if !file.Ignored {
			var hashes []FileHashes = getAllHashesDB(tx, file.FileName)
            var issue = false

			for i, nhash := range hashes {
				if i+1 < len(hashes) {
					var ohash FileHashes = hashes[i+1]
					if nhash.Size != ohash.Size {
						fmt.Println(nhash.FileName, ":::: changed sizes between", ohash.ScanDate, "and", nhash.ScanDate)
                        issue = true
						break
					} else if nhash.MD5 != "" && ohash.MD5 != "" && nhash.MD5 != ohash.MD5 {
						fmt.Println(nhash.FileName, ":::: md5 changed between", ohash.ScanDate, "and", nhash.ScanDate)
                        issue = true
						break
					} else if nhash.SHA1 != "" && ohash.SHA1 != "" && nhash.SHA1 != ohash.SHA1 {
						fmt.Println(nhash.FileName, ":::: sha1 changed between", ohash.ScanDate, "and", nhash.ScanDate)
                        issue = true
						break
					} else if nhash.SHA256 != "" && ohash.SHA256 != "" && nhash.SHA256 != ohash.SHA256 {
						fmt.Println(nhash.FileName, ":::: sha256 changed between", ohash.ScanDate, "and", nhash.ScanDate)
                        issue = true
						break
					} else if nhash.Tiger != "" && ohash.Tiger != "" && nhash.Tiger != ohash.Tiger {
						fmt.Println(nhash.FileName, ":::: tiger changed between", ohash.ScanDate, "and", nhash.ScanDate)
                        issue = true
						break
					} else if nhash.Whirlpool != "" && ohash.Whirlpool != "" && nhash.Whirlpool != ohash.Whirlpool {
						fmt.Println(nhash.FileName, ":::: whirlpool changed between", ohash.ScanDate, "and", nhash.ScanDate)
                        issue = true
						break
					}
				}
			}

            if !issue {
                issue = true
            }
		}
	}

	tx.Commit()
}
