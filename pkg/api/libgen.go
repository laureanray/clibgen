package api

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
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
	log.Println("Downloading book ", link)

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

// TODO: Introduce proper types with pagination
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

func BookPrinter(book []Book) {
	for i, v := range book {
		fmt.Printf("-[%d]----------------\n", i)

		fmt.Println("Title: ", v.Title)
		fmt.Println("Year: ", v.Year, " Ext: ", v.Extension)
		fmt.Println("Author: ", v.Author)
		//fmt.Println("----------------------")
	}
}

// Downloads the file to current working directorhy
func DownloadSelection(selectedBook Book) {

	path, err := os.Getwd()
	link := getDirectDownloadLink(selectedBook.Mirrors[0])

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Downloading file to current directory: ", path)

	out, err := os.Create(selectedBook.Title + "." + selectedBook.Extension)
	defer out.Close()

	resp, err := http.Get(link)
	defer resp.Body.Close()

	n, err := io.Copy(out, resp.Body)

	fmt.Println("written: ", n)
}
