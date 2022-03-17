package main

import (
	"fmt"
	_ "net/http"
	"os"

	"github.com/scanhamman/cessda_data/internal/fileio"
)

func main() {

	/*
		var thisToken string = ""
		var file_end bool = false
		var n int = 0

		// store identifiers for records
		for !file_end {
			nextToken, fileend, err := fileio.ScrapeHeaderList(thisToken)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			thisToken = nextToken
			file_end = fileend
			n++
			fmt.Printf("Count: %d\n", n)
		}
	*/

	// Use fixed id for now
	var id = "e5ee9b2c41c34cb48bd639678c49562a"
	foo, err := fileio.ScrapeDetailsFile(id)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if foo == "" {
		os.Exit(1)
	}
	println(foo)
}
