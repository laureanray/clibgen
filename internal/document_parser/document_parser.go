package documentparser

import (
	"github.com/laureanray/clibgen/internal/book"
)

// type Page struct {
// 	doc *goquery.Document
// }

type DocumentParser interface {
  GetBookDataFromDocument() []book.Book
  GetDownloadLinkFromDocument() (string, bool)
}

