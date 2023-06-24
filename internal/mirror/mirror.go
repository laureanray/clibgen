package mirror

import (
	"github.com/laureanray/clibgen/internal/book"
	"github.com/laureanray/clibgen/internal/libgen"
)

type Mirror interface {
  SearchByTitle(query string) ([]book.Book, error)
  // SearchByAuthor(author string) []book.Book
  // SearchByISBN(isbn string) []book.Book
  // 1GetDownloadLink(book book.Book) string
  DownloadSelection(book book.Book)
}

// TODO: Make this persistent
type Configuration struct {
  numberOfResults int
}

type NewMirror struct {
  domain libgen.Domain
}
