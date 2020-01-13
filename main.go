package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

func main() {
	c := colly.NewCollector(
		// Visit only these
		colly.AllowedDomains("www.sjcc.edu", "sjcc.edu"),
		// Cache responses to prevent multiple download of pages
		// even if the collector is restarted
		colly.CacheDir("./sjcc_cache"),
	)

	// On every <a> element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		url := e.Attr("href")

		// Print link
		fmt.Printf("Link found: %q -> %s\n", e.Text, url)
		// Visit link found on page
		// Only those links are visited which are in Allowed Domains
		c.Visit(e.Request.AbsoluteURL(url))
	})

	// Before making a request
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Set error handling
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL: ", r.Request.URL, "failed with response: ", r, "\nError: ", err)
	})

	// start scraping
	c.Visit("https://www.sjcc.edu")
}
