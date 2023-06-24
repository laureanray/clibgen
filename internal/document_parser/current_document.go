package documentparser

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/laureanray/clibgen/internal/book"
)

type CurrentDocumentParser struct {
	doc *goquery.Document
}

func NewCurrentDocumentParser(document *goquery.Document) *CurrentDocumentParser {
	return &CurrentDocumentParser{doc: document}
}

func NewCurrentDocumentParserFromReader(r io.Reader) *CurrentDocumentParser {
	document, _ := goquery.NewDocumentFromReader(r)
	return &CurrentDocumentParser{doc: document}
}

func (cdp *CurrentDocumentParser) GetBookDataFromDocument() []book.Book {
	var books []book.Book
	cdp.doc.Find("#tablelibgen > tbody > tr").Each(func(resultsRow int, bookRow *goquery.Selection) {
		var id, author, title, publisher, extension, year, fileSize string
		var mirrors []string
		if resultsRow != 0 {
			bookRow.Find("td").Each(func(column int, columnSelection *goquery.Selection) {
				switch column {
				case 0:
					title = columnSelection.Find("a").First().Text()
				case 1:
					author = columnSelection.Text()
				case 2:
					publisher = columnSelection.Text()
				case 3:
					year = columnSelection.Text()
				case 6:
					fileSize = columnSelection.Text()
				case 7:
					extension = columnSelection.Text()
				case 8:
					columnSelection.Find("a").Each(func(linkCol int, link *goquery.Selection) {
						href, hrefExists := link.Attr("href")
						if hrefExists {
							mirrors = append(mirrors, href)
						}
					})
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

func (cdp *CurrentDocumentParser) getDownloadLinkFromDocument() (string, bool) {
  return cdp.doc.Find("#main a").First().Attr("href")
}

func (cdp *CurrentDocumentParser) getBookTitleFromSelection(selection *goquery.Selection) string {
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

func (cdp *CurrentDocumentParser) GetDirectDownloadLink(selectedBook book.Book) string {
	fmt.Println("Obtaining direct download link")

  // TODO Implement retry?
  link := selectedBook.Mirrors[0]

	resp, err := http.Get(link)
	defer func(Body io.ReadCloser) {
		err := Body.Close()

		if err != nil {
			fmt.Println("Error closing body:", err)
		}
	}(resp.Body)

	if err != nil {
		fmt.Println("Error getting response:", err)
	}

	directDownloadLink, exists :=
		NewCurrentDocumentParserFromReader(resp.Body).getDownloadLinkFromDocument()

	if exists {
		return directDownloadLink
	}

	return ""
}
