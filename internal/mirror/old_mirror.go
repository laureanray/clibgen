package mirror

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
	"github.com/laureanray/clibgen/internal/book"
	"github.com/laureanray/clibgen/internal/domain"
)

type OldMirror struct {
  domain domain.Domain
  config Configuration
}

func NewOldMirror(domain domain.Domain) *OldMirror {
  return &OldMirror{
    domain: domain,
    config: Configuration{
      // TODO: Make this configurable
      numberOfResults: 5,
    },
  }
}

func (m *OldMirror) SearchByTitle(query string) []book.Book {
  return []book.Book{}
}

// Search the libgen site returns the document 
// of the search results page
func (m *OldMirror) searchSite() (*goquery.Document, error) {
  baseUrl := fmt.Sprintf("https://libgen.%s/search.php", m.domain)
	queryString := fmt.Sprintf("%s?req=%s&res=25&view=simple&phrase=1&column=%s", baseUrl, url.QueryEscape(query))
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

	document, e = goquery.NewDocumentFromReader(resp.Body)

	if e != nil {
		fmt.Println(e)
		return nil, e
	}

	return document, e
}

