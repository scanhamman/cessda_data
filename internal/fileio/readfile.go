package fileio

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/scanhamman/cessda_data/internal/database"

	"github.com/PuerkitoBio/goquery"
)

func ScrapeHeaderList(token string) (string, bool, error) {

	// Request the HTML page.
	var url = ""
	if token == "" {
		url = "https://datacatalogue.cessda.eu/oai-pmh/v0/oai?verb=ListIdentifiers&metadataPrefix=oai_ddi25"
	} else {
		url = "https://datacatalogue.cessda.eu/oai-pmh/v0/oai?verb=ListIdentifiers&metadataPrefix=oai_ddi25&resumptionToken=" + token
	}

	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		fmt.Println(res.Status)
		os.Exit(1)
	}

	time.Sleep(2 * time.Second)

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// get the ListIdentifiers element
	list := doc.Find("ListIdentifiers")
	var idents []database.Identifier

	// get resumption token and identifiers list
	thistoken := list.Find("resumptionToken")
	thistokenvalue := thistoken.Text()
	list.Find("header").Each(func(i int, header *goquery.Selection) {
		id := header.Find("identifier").Text()
		status, _ := header.Attr("status")
		idents = append(idents, database.Identifier{Id: id, Status: status})
		fmt.Println(id, status)
	})

	var file_ended bool = false
	if thistokenvalue == "" {
		file_ended = true
	}
	fmt.Println(thistokenvalue, file_ended, err)
	database.StoreIdentifiers(idents)
	return thistokenvalue, file_ended, err
}
