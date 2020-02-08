package main

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
)

// Link stores information about links
type Link struct {
	Title    string
	URL      string
	Origin   string
	FormType string
} 

func main() {
	c := colly.NewCollector(
		// Visit only these
		colly.AllowedDomains("www.sjcc.edu", "sjcc.edu", "catalog.sjcc.edu"),
		// Cache responses to prevent multiple download of pages
		// even if the collector is restarted
		colly.CacheDir("./sjcc_cache"),
	)

	links := make(map[string]Link)
	var currentPage string

	links["0"] = Link{
		Title:    "Title",
		URL:      "URL",
		Origin:   "Origin",
		FormType: "Form Type",
	}

	// On every <a> element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		url := e.Attr("href")

		switch {
		// ? query for .pdf
		case strings.Contains(url, ".pdf"):
			if !strings.Contains(url, "://") {
				url = "https://www.sjcc.edu" + url
			}
			links["https://www.sjcc.edu"+url] = Link{
				Title:    e.Text,
				URL:      url,
				Origin:   currentPage,
				FormType: "PDF",
			}

		// ? query for Google Forms
		case (strings.Contains(url, "docs.google.com/forms") || strings.Contains(url, "goo.gl/forms")):
			links["https://www.sjcc.edu"+url] = Link{
				Title:    e.Text,
				URL:      url,
				Origin:   currentPage,
				FormType: "Google Docs",
			}

		// ? query for Office Forms
		case strings.Contains(url, "forms.office.com/"):
			links["https://www.sjcc.edu"+url] = Link{
				Title:    e.Text,
				URL:      url,
				Origin:   currentPage,
				FormType: "Office Forms",
			}

		// ? query for formsite
		case strings.Contains(url, "formsite.com"):
			links["https://www.sjcc.edu"+url] = Link{
				Title:    e.Text,
				URL:      url,
				Origin:   currentPage,
				FormType: "FormSite",
			}

		// TODO add a case for smartsheet forms
		case strings.Contains(url, "smartsheet.com/b/form"):
			links["https://www.sjcc.edu"+url] = Link{
				Title:    e.Text,
				URL:      url,
				Origin:   currentPage,
				FormType: "SmartSheets",
			}

		// check if link contains form in title or link
		// ! refine query for "form" as word not substring
		// case strings.Contains(e.Text, "form") || strings.Contains(url, "form"):
		// 	if !strings.Contains(url, "://") {
		// 		url = "https://www.sjcc.edu" + url
		// 	}
		// 	links["https://www.sjcc.edu"+url] = Link{
		// 		Title:    e.Text,
		// 		URL:      url,
		// 		Origin:   currentPage,
		// 		FormType: "Form",
		// 	}
		}

		// Print link
		// fmt.Printf("Link found: %q -> %s\n", e.Text, url)
		// Visit link found on page
		// Only those links are visited which are in Allowed Domains
		c.Visit(e.Request.AbsoluteURL(url))
	})

	// Before making a request
	c.OnRequest(func(r *colly.Request) {
		// fmt.Println("Visiting", r.URL.String())
		currentPage = r.URL.String()
	})

	// ! Set error handling
	// c.OnError(func(r *colly.Response, err error) {
	// 	fmt.Println("Request URL: ", r.Request.URL, "failed with response: ", r, "\nError: ", err)
	// })

	// start scraping
	c.Visit("http://www.sjcc.edu")

	// TODO print as CSV format
	// TODO remove origin link
	for _, record := range links {
		fmt.Println(record)
	}
	
	// fmt.Printf("%v \n", links)
}
