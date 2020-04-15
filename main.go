package main

import (
	"log"
	"os"
	"strings"

	//"github.com/josephsintum/go-webscraper/write.go"
	"github.com/gocolly/colly"
)

// Link stores information about links
type Link struct {
	Title    string
	URL      string
	Origin   string
	FormType string
}

func write(filename string, data string) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := file.WriteString(data); err != nil {
		file.Close()
		// ignore error; Write error takes precedence
		log.Fatal(err)
	}
	if err := file.Close(); err != nil {
		log.Fatal(err)
	}
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
				Origin:   e.Request.URL.String(),
				FormType: "PDF",
			}

		// ? query for Google Forms
		case strings.Contains(url, "docs.google.com/forms") || strings.Contains(url, "goo.gl/forms"):
			links["https://www.sjcc.edu"+url] = Link{
				Title:    e.Text,
				URL:      url,
				Origin:   e.Request.URL.String(),
				FormType: "Google Docs",
			}

		// ? query for Office Forms
		case strings.Contains(url, "forms.office.com/"):
			links["https://www.sjcc.edu"+url] = Link{
				Title:    e.Text,
				URL:      url,
				Origin:   e.Request.URL.String(),
				FormType: "Office Forms",
			}

		// ? query for formsite
		case strings.Contains(url, "formsite.com"):
			links["https://www.sjcc.edu"+url] = Link{
				Title:    e.Text,
				URL:      url,
				Origin:   e.Request.URL.String(),
				FormType: "FormSite",
			}

		// ? query for smartsheets
		case strings.Contains(url, "smartsheet.com/b/form"):
			links["https://www.sjcc.edu"+url] = Link{
				Title:    e.Text,
				URL:      url,
				Origin:   e.Request.URL.String(),
				FormType: "SmartSheets",
			}

		// ? query for docusign
		case strings.Contains(url, "docusign.net/Member/PowerFormSigning.aspx"):
			links["https://www.sjcc.edu"+url] = Link{
				Title:    e.Text,
				URL:      url,
				Origin:   e.Request.URL.String(),
				FormType: "DocuSign",
			}

		}

		// Print link
		// fmt.Printf("Link found: %q -> %s\n", e.Text, url)
		// Visit link found on page
		// Only those links are visited which are in Allowed Domains
		err := c.Visit(e.Request.AbsoluteURL(url))
		if err != nil {
		}
	})

	// ! Set error handling
	// c.OnError(func(r *colly.Response, err error) {
	// 	fmt.Println("Request URL: ", r.Request.URL, "failed with response: ", r, "\nError: ", err)
	// })

	// start scraping
	err := c.Visit("http://www.sjcc.edu")
	if err != nil {
	}

	for _, record := range links {
		write("results.csv", record.Title+"\t"+record.URL+"\t"+record.FormType+"\t"+record.Origin+"\n")
	}

	write("results.csv", "domino and sunlight")
}
