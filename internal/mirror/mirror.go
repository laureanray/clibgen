package mirror

import (
	"fmt"
	"io"
	"net/http"

	"github.com/laureanray/clibgen/internal/book"
	"github.com/laureanray/clibgen/internal/libgen"
)

type Mirror interface {
  SearchByTitle(query string) ([]book.Book, error)
  // SearchByAuthor(author string) []book.Book
  // SearchByISBN(isbn string) []book.Book
  // 1GetDownloadLink(book book.Book) string
}

// TODO: Make this persistent
type Configuration struct {
  numberOfResults int
}

type NewMirror struct {
  domain libgen.Domain
}
