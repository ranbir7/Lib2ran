package internal

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func SearchLibgen(query string) []Book {
	var books []Book
	client := &http.Client{Timeout: 30 * time.Second}
	searchPath := "/search.php?req=%s&res=100&view=simple&phrase=1&column=def"
	queryEscaped := strings.ReplaceAll(query, " ", "+")

	for _, mirror := range LibgenMirrors {
		url := fmt.Sprintf("%s%s", mirror, fmt.Sprintf(searchPath, queryEscaped))
		resp, err := client.Get(url)
		if err != nil {
			continue
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			continue
		}
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			continue
		}
		doc.Find("table.c tr").Each(func(i int, s *goquery.Selection) {
			if i == 0 {
				return // skip header
			}
			tds := s.Find("td")
			if tds.Length() < 9 {
				return
			}
			book := Book{
				Title:     strings.TrimSpace(tds.Eq(2).Text()),
				Author:    strings.TrimSpace(tds.Eq(1).Text()),
				Year:      strings.TrimSpace(tds.Eq(4).Text()),
				Size:      strings.TrimSpace(tds.Eq(7).Text()),
				Extension: strings.TrimSpace(tds.Eq(8).Text()),
				Mirrors:   make(map[string]string),
			}
			// Get download page links from the last few columns
			tds.Slice(9, tds.Length()).Each(func(j int, td *goquery.Selection) {
				link, exists := td.Find("a").Attr("href")
				if exists && strings.HasPrefix(link, "http") {
					book.Mirrors[fmt.Sprintf("Mirror%d", j+1)] = link
				}
			})
			books = append(books, book)
		})
		if len(books) > 0 {
			break // stop at first mirror with results
		}
	}
	return books
}
