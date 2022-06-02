package api

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// TODO: Introduce proper types with pagination
func SearchBookByTitle(query string) ([]Book, error) {
	log.Println("Query: ", query, url.QueryEscape(query))
	queryString := fmt.Sprintf("https://libgen.is/search.php?req=%s&res=25&view=simple&phrase=1&column=def", url.QueryEscape(query))
	var e error
	resp, err := http.Get(queryString)

	log.Println("HTTP GET: ", queryString)

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			e = err
		}
	}(resp.Body)

	if err != nil {
		e = err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		e = err
	}

	var books []Book

	doc.Find(".c > tbody > tr").Each(func(i int, s *goquery.Selection) {
		var id, author, title string

		if i != 0 {
			s.Find("td").Each(func(column int, s *goquery.Selection) {
				switch column {
				//case 0:
				case 0:
					id = s.Text()
				case 1:
					author = s.Text()
				case 2:
					s.Find("a").Each(func(v int, s *goquery.Selection) {
						_, exists := s.Attr("title")

						if exists {
							title = s.Text()
						}
					})
					s.Find("a > font").Each(func(v int, s *goquery.Selection) {
						a := s.Text()
						title = strings.ReplaceAll(title, a, "")
						log.Println(title, a)
					})

					fmt.Println("Q: ", title)
				case 3:

				}
			})
			books = append(books, Book{
				ID:     id,
				Author: author,
				Title:  title,
			})

		}
	})

	return books, e
}
