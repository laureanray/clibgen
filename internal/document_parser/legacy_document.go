package documentparser

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/laureanray/clibgen/internal/book"
)

type LegacyDocumentParser struct {
  doc *goquery.Document
}

func NewLegacyDocumentParser(document *goquery.Document) *LegacyDocumentParser {
	return &LegacyDocumentParser{doc: document}
}

func NewLegacyDocumentParserFromReader(r io.Reader) *LegacyDocumentParser {
  document, _ := goquery.NewDocumentFromReader(r)
  return &LegacyDocumentParser{doc: document}
}

func (ldp *LegacyDocumentParser) GetBookDataFromDocument() []book.Book {
	var books []book.Book
	ldp.doc.Find(".c > tbody > tr").Each(func(resultsRow int, bookRow *goquery.Selection) {
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


func (ldp *LegacyDocumentParser) getDownloadLinkFromDocument() (string, bool){
  return ldp.doc.Find("#download > ul > li > a").First().Attr("href")
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


func GetDirectDownloadLinkFromLegacy(link string) string {
	fmt.Println("Obtaining direct download link")
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

  page := NewLegacyDocumentParserFromReader(resp.Body)
  // TODO: I think this can be improved
  directDownloadLink, exists := page.getDownloadLinkFromDocument()

  if exists {
    return directDownloadLink
  }
  
  return ""
}


