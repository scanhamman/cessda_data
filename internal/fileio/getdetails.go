package fileio

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func ScrapeDetailsFile(id string) (string, error) {

	// Request the HTML page.
	var url = ""
	if id == "" {
		url = "https://datacatalogue.cessda.eu/oai-pmh/v0/oai?verb=GetRecord&metadataPrefix=oai_ddi25&identifier"
	} else {
		url = "https://datacatalogue.cessda.eu/oai-pmh/v0/oai?verb=GetRecord&metadataPrefix=oai_ddi25&identifier=" + id
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

	// Check the XML document, initially as a text string
	b, err := io.ReadAll(res.Body)
	if err != nil {
		os.Exit(1)
	}
	d := string(b)

	// check if d is nil?

	// replace html tags
	d = ReplaceAngleBrackets(d)
	d = ReplaceBrs(d)
	d = ReplaceNonbreakingSpaces(d)
	d = CleanTags(d)

	// Load the XML document
	var s string
	r := strings.NewReader(s)
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// get the ListIdentifiers element
	doc.Find("ListIdentifiers")

	return d, nil
}

func ReplaceAngleBrackets(s string) string {
	s = strings.ReplaceAll(s, "&lt;", "<")
	return strings.ReplaceAll(s, "&gt;", ">")
}

func ReplaceBrs(s string) string {
	replacer := strings.NewReplacer("<br>", "\n", "<BR>", "\n", "<br/>", "\n", "<br />", "\n", "<br/ >", "\n", "< br / >", "\n")
	return replacer.Replace(s)
}

func ReplaceNonbreakingSpaces(s string) string {
	s = strings.ReplaceAll(s, "&amp;nbsp;", " ")
	return s
}

func CleanTags(s string) string {
	doc := []rune(s)
	var t strings.Builder // string builder for building destinatioon string
	tw := []rune{}        // word or phrase in tag
	bt := []rune{}        // word or phrase between tags

	in_tag := false
	for i := 0; i < len(doc); i++ {
		c := doc[i]

		if c == '<' {
			in_tag = true
			// add current bw (between words)
			element_text := string(bt)
			t.WriteString(element_text)
			tw = []rune{'<'} // start new tag word

		} else if c == '>' {
			// complete the tag words
			// if the tag is not one to ignore...
			tag := string(tw) + ">"
			lose_tag, tag_type := CheckTagToLose(tag)

			if tag_type == "p" {
				t.WriteString("\n")
				tag = "\n"
			}
			if tag_type == "li" {
				t.WriteString("\n\t-> ")
				tag = "\n\t -> "
			}

			if !lose_tag {
				t.WriteString(tag) // Add to t

				// strip <> , trim to first space to get tag name
				// check does not end with a />
				if strings.HasPrefix(tag, "/") {
					// end tag - pop from tag stack
					// (get tag name) if it is the same as top of stack
					println(tag)
					// if not there is a problem
					// if the tag is lower down the stack - pop that
					// and add any non-matched tags as text but mark with asterisks

					// if not matched something weird has happened - post an error

				} else {
					// start tag - push to tag stack
					// as tag name
					println(tag)

				}

			}

			in_tag = false
			bt = []rune{}
		} else {
			if in_tag {
				tw = append(tw, c)
			} else {
				bt = append(bt, c)
			}
		}
	}

	return t.String()
}

func CheckTagToLose(tag string) (bool, string) {
	res := false
	tag_type := ""
	tag = strings.ToLower(tag) // make tag lower case

	if tag == "</i>" || tag == "</b>" ||
		tag == "</u>" || tag == "</p>" {
		res = true
	} else if tag == "</div>" || tag == "</span>" ||
		tag == "</strong>" || tag == "</em>" {
		res = true
	} else if tag == "</ul>" || tag == "</ol>" ||
		tag == "</li>" {
		res = true
	}

	tag = strings.ReplaceAll(tag, ">", " ") // change > to space

	if strings.HasPrefix(tag, "<p ") {
		tag_type = "p"
		res = true
	} else if strings.HasPrefix(tag, "<i ") ||
		strings.HasPrefix(tag, "<b ") ||
		strings.HasPrefix(tag, "<u ") {
		res = true
	} else if strings.HasPrefix(tag, "<div ") ||
		strings.HasPrefix(tag, "<span ") ||
		strings.HasPrefix(tag, "<strong ") ||
		strings.HasPrefix(tag, "<em ") {
		res = true
	} else if strings.HasPrefix(tag, "<li ") {
		tag_type = "li"
		res = true
	} else if strings.HasPrefix(tag, "<ul ") ||
		strings.HasPrefix(tag, "<ol ") {
		res = true
	}

	return res, tag_type
}
