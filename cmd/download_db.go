package main

import (
	"fmt"
	"os"

	"github.com/scanhamman/cessda_data/internal/fileio"
)

func main() {

	var thisToken string = ""
	var file_end bool = false
	var n int = 0

	// store identifiers for records
	for !file_end && n < 50 {
		nextToken, fileend, err := fileio.ScrapeHeaderList(thisToken)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		thisToken = nextToken
		file_end = fileend
		n++
	}

}
