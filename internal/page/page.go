package page

import "github.com/PuerkitoBio/goquery"

type Page struct {
	doc *goquery.Document
}

func New(document *goquery.Document) *Page {
	p := &Page{doc: document}

	return p
}
