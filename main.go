package main

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
)

// Link stores information about links
type Link struct {
	Title  string
	URL    string
	Origin string
}

func main() {
	c := colly.NewCollector(
		// Visit only these
		colly.AllowedDomains("www.sjcc.edu", "sjcc.edu"),
		// Cache responses to prevent multiple download of pages
		// even if the collector is restarted
		colly.CacheDir("./sjcc_cache"),
	)

	links := make(map[string]Link)
	var currentPage string

	links["0"] = Link{Title: "Title", URL: "URL", Origin: "Origin"}

	// On every <a> element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		url := e.Attr("href")

		switch true {
		// check for .pdf
		case strings.Contains(url, ".pdf"):
			fmt.Println(".pdf")

		// check for Google Forms
		case (strings.Contains(url, "docs.google.com/forms") || strings.Contains(url, "goo.gl/forms")):
			fmt.Println("google docs")

		// check for Office Forms
		case strings.Contains(url, "forms.office.com/"):
			fmt.Println("google docs")

		// check for formsite
		case strings.Contains(url, "formsite.com"):
			fmt.Println("formsite")

		// check if link contains form in title
		case strings.Contains(e.Text, "form"):
			fmt.Println("form")

		default: 
			fmt.Println("nothing")
		}


		links[url] = Link{
			Title:  e.Text,
			URL:    url,
			Origin: currentPage,
		}

		// Print link
		fmt.Printf("Link found: %q -> %s\n", e.Text, url)
		// Visit link found on page
		// Only those links are visited which are in Allowed Domains
		c.Visit(e.Request.AbsoluteURL(url))
	})

	// Before making a request
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
		currentPage = r.URL.String()
	})

	// Set error handling
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL: ", r.Request.URL, "failed with response: ", r, "\nError: ", err)
	})

	// start scraping
	c.Visit("https://www.sjcc.edu")

	// fmt.Printf("%v \n", links)
}
