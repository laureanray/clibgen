package api

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/schollz/progressbar/v3"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// Get book title from the selection, in most cases the title is hidden through nested anchor tags.
// In order to produce a clean output extra texts are also removed.
func getBookTitleFromSelection(selection *goquery.Selection) string {
	var title string
	selection.Find("a").Each(func(v int, s *goquery.Selection) {
		_, exists := s.Attr("title")
		if exists {
			title = s.Text()
		}
	})
	selection.Find("a > font").Each(func(v int, s *goquery.Selection) {
		a := s.Text()
		title = strings.ReplaceAll(title, a, "")
	})
	return title
}

// Parse HTML and get the data from the table by parsing and iterating through them.
func getBookDataFromDocument(document *goquery.Document) []Book {
	var books []Book
	document.Find(".c > tbody > tr").Each(func(resultsRow int, bookRow *goquery.Selection) {
		var id, author, title, publisher, extension, year string
		var mirrors []string
		if resultsRow != 0 {
			bookRow.Find("td").Each(func(column int, columnSelection *goquery.Selection) {
				switch column {
				case 0:
					id = columnSelection.Text()
				case 1:
					author = columnSelection.Text()
				case 2:
					title = getBookTitleFromSelection(columnSelection)
				case 3:
					publisher = columnSelection.Text()
				case 4:
					year = columnSelection.Text()
				case 8:
					extension = columnSelection.Text()
				case 9, 10, 11:
					href, hrefExists := columnSelection.Find("a").Attr("href")
					if hrefExists {
						mirrors = append(mirrors, href)
					}
				}
			})
			books = append(books, Book{
				ID:        id,
				Author:    author,
				Year:      year,
				Title:     title,
				Publisher: publisher,
				Extension: extension,
				Mirrors:   mirrors,
			})

		}
	})

	return books
}

func getLinkFromDocument(document *goquery.Document) (string, bool) {
	return document.Find("#download > h2 > a").First().Attr("href")
}

func getDirectDownloadLink(link string) string {
	log.Println("Downloading book from: ", link)

	resp, err := http.Get(link)

	defer func(Body io.ReadCloser) {
		err := Body.Close()

		if err != nil {
			log.Fatal(err)
		}
	}(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	document, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	directDownloadLink, exists := getLinkFromDocument(document)

	if exists {
		return directDownloadLink
	}

	// do something
	return ""
}

func SearchBookByTitle(query string, limit int) ([]Book, error) {
	var e error
	queryString := fmt.Sprintf("https://libgen.is/search.php?req=%s&res=25&view=simple&phrase=1&column=def", url.QueryEscape(query))
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

	document, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		e = err
	}

	books := getBookDataFromDocument(document)

	if len(books) >= limit {
		books = books[:limit]
	}

	return books, e
}

// DownloadSelection Downloads the file to current working directory
func DownloadSelection(selectedBook Book) {
	link := getDirectDownloadLink(selectedBook.Mirrors[0])
	req, _ := http.NewRequest("GET", link, nil)
	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()
	f, _ := os.OpenFile(strings.Trim(selectedBook.Title, " ")+"."+selectedBook.Extension, os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()

	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		"downloading",
	)
	io.Copy(io.MultiWriter(f, bar), resp.Body)

	fmt.Println("File downloaded. ", f.Name())
}
