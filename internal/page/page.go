package page

import (
	"fmt"
	"io"
	"strings"
  "net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/laureanray/clibgen/internal/book"
	"github.com/laureanray/clibgen/internal/mirror"
)

type Page struct {
	doc *goquery.Document
}

func New(document *goquery.Document) *Page {
	return &Page{doc: document}
}

func NewFromReader(r io.Reader) *Page {
  document, _ := goquery.NewDocumentFromReader(r)
  return &Page{doc: document}
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


func (p *Page) getDownloadLinkFromDocument(m mirror.Mirror) (string, bool){
  switch m.(type) {
    case *mirror.OldMirror:
    return p.doc.Find("#download > ul > li > a").First().Attr("href")
  }

  return "", false
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


func (p *Page) GetDirectDownloadLink(mirror mirror.Mirror, link string) string {
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

  page := NewFromReader(resp.Body)
  // TODO: I think this can be improved
  directDownloadLink, exists := page.getDownloadLinkFromDocument(mirror)

  if exists {
    return directDownloadLink
  }
  
  return ""
}
