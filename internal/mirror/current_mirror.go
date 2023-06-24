package mirror

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
	"github.com/laureanray/clibgen/internal/book"
	"github.com/laureanray/clibgen/internal/console"
	"github.com/laureanray/clibgen/internal/document_parser"
	"github.com/laureanray/clibgen/internal/downloader"
	"github.com/laureanray/clibgen/internal/libgen"
)

type CurrentMirror struct {
  domain libgen.Domain
  filter libgen.Filter
  config Configuration
}

func NewCurrentMirror(domain libgen.Domain) *CurrentMirror {
  return &CurrentMirror{
    domain: domain,
    // TODO: Make this configurable
    filter: libgen.TITLE,
    config: Configuration{
      numberOfResults: 5,
    },
  }
}

func (m *CurrentMirror) SearchByTitle(query string) ([]book.Book, error) {
	fmt.Println("Searching for: ", console.Higlight(query))
	var document *goquery.Document

  document, err := m.searchSite(query)

	if err != nil {
		fmt.Println(console.Error("Error searching for book: %s", query))
    // TODO: Implement retrying
		// fmt.Println(infoColor("Retrying with other site"))
		// document, e = searchLibgen(query, siteToUse) // If this also fails then we have a problem
	}
	fmt.Println(console.Success("Search complete, parsing the document..."))
  
  page := documentparser.NewCurrentDocumentParser(document)
  bookResults := page.GetBookDataFromDocument()

	// if len(bookResults) >= limit {
	// 	bookResults = bookResults[:limit]
	// }

	return bookResults, err
}


// Search the libgen site returns the document 
// of the search results page
func (m *CurrentMirror) searchSite(query string) (*goquery.Document, error) {

  baseUrl := fmt.Sprintf("https://libgen.%s/index.php", m.domain)

	queryString := fmt.Sprintf(
    "%s?req=%s&res=25&view=simple&phrase=1&column=%s",
    baseUrl,
    url.QueryEscape(query),
    m.filter,
  )

  fmt.Println(console.Info(queryString))

	resp, e := http.Get(queryString)

	if e != nil {
		return nil, e
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			e = err
		}
	}(resp.Body)

  document, e := goquery.NewDocumentFromReader(resp.Body)

	if e != nil {
		fmt.Println(e)
		return nil, e
	}

	return document, e
}

func (m *CurrentMirror) DownloadSelection(selectedBook book.Book) {
  fmt.Println(console.Info("Downloading book..."))
  directLink := documentparser.GetDirectDownloadLinkFromCurrent(selectedBook.Mirrors[0])
  downloader.NewDownloader(selectedBook, directLink).Download()
}
