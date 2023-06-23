package mirror

import (
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
func (m *OldMirror) searchSite() {

}

