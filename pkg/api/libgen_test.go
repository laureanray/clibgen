package api

import (
	"testing"
)

func TestBook(t *testing.T) {
	books, err := SearchBookByTitle("Search")

	if len(books) == 0 || err != nil {
		t.Fatalf(`Empty books: %d`, len(books))
	}
}
