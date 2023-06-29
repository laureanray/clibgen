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

	if len(bookResults) >= m.config.numberOfResults {
		bookResults = bookResults[:m.config.numberOfResults]
	}

	return bookResults, err
}

func (m *CurrentMirror) SearchByAuthor(query string) ([]book.Book, error) {
	fmt.Println("Searching for: ", console.Higlight(query))
	var document *goquery.Document

  m.filter = libgen.AUTHOR
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

	if len(bookResults) >= m.config.numberOfResults {
		bookResults = bookResults[:m.config.numberOfResults]
	}

	return bookResults, err
}


// Search the libgen site returns the document 
// of the search results page
func (m *CurrentMirror) searchSite(query string) (*goquery.Document, error) {

  baseUrl := fmt.Sprintf("https://libgen.%s/index.php", m.domain)

	queryString := fmt.Sprintf(
    "%s?req=\"%s\"",
    baseUrl,
    url.QueryEscape(query),
  )

  filter := string(string(m.filter)[0])

  reqString:= queryString + "&columns%5B%5D=" + filter + "&objects%5B%5D=f&objects%5B%5D=e&objects%5B%5D=s&objects%5B%5D=a&objects%5B%5D=p&objects%5B%5D=w&topics%5B%5D=l&topics%5B%5D=c&topics%5B%5D=f&res=25&gmode=on&filesuns=all"

  fmt.Println(reqString)

	resp, e := http.Get(reqString)

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
