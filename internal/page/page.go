package page

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/laureanray/clibgen/internal/book"
)

type Page struct {
	doc *goquery.Document
}

func New(document *goquery.Document) *Page {
	p := &Page{doc: document}

	return p
}

func (p *Page) GetBookDataFromDocument() []book.Book {
	var books []book.Book
	p.doc.Find(".c > tbody > tr").Each(func(resultsRow int, bookRow *goquery.Selection) {
		var id, author, title, publisher, extension, year, fileSize string
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
				case 7:
					fileSize = columnSelection.Text()
				case 8:
					extension = columnSelection.Text()
				case 9, 10, 11:
					href, hrefExists := columnSelection.Find("a").Attr("href")
					if hrefExists {
						mirrors = append(mirrors, href)
					}
				}
			})
			books = append(books, book.Book{
				ID:        id,
				Author:    author,
				Year:      year,
				Title:     title,
				Publisher: publisher,
				Extension: extension,
				Mirrors:   mirrors,
				FileSize:  fileSize,
			})
		}
	})

	return books
}


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
