package mirror

import (
	"github.com/laureanray/clibgen/internal/book"
	"github.com/laureanray/clibgen/internal/domain"
)

type Mirror interface {
  SearchByTitle(title string) []book.Book
  SearchByAuthor(author string) []book.Book
  SearchByISBN(isbn string) []book.Book
  GetDownloadLink(book book.Book) string
}

// TODO: Make this persistent
type Configuration struct {
  numberOfResults int
}

type NewMirror struct {
  domain domain.Domain
}
