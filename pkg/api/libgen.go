package api

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
	"github.com/kennygrant/sanitize"
	"github.com/schollz/progressbar/v3"
)

/*

Currently there are three libgen domains:

 - libgen.is -> primary domain (old interface)
 - libgen.li -> secondary (fallback) newer interface

These domains might change in the future
*/

type Site int32

const (
	LibgenOld Site = 0
	LibgenNew Site = 1
)

// Note: Applicable only for libgen.is! (We might need to implemented something like this for libgen.li)
// Get book title from the selection, in most cases the title is hidden through nested anchor tags.
// In order to produce a clean output extra texts are also removed.
func getBookTitleFromSelection(selection *goquery.Selection) string {
	var title string
	selection.Find("a").Each(func(v int, s *goquery.Selection) {
		_, exists := s.Attr("title")
		if exists {
			title = s.Text()
		}
	})
	selection.Find("a > font").Each(func(v int, s *goquery.Selection) {
		a := s.Text()
		title = strings.ReplaceAll(title, a, "")
	})
	return title
}

func getBookDataFromDocumentOld(document *goquery.Document) []Book {
	var books []Book
	document.Find(".c > tbody > tr").Each(func(resultsRow int, bookRow *goquery.Selection) {
		var id, author, title, publisher, extension, year, fileSize string
		var mirrors []string
		if resultsRow != 0 {
			bookRow.Find("td").Each(func(column int, columnSelection *goquery.Selection) {
				switch column {
				case 0:
					id = columnSelection.Text()
				case 1:
					author = columnSelection.Text()
				case 2:
					title = getBookTitleFromSelection(columnSelection)
				case 3:
					publisher = columnSelection.Text()
				case 4:
					year = columnSelection.Text()
				case 7:
					fileSize = columnSelection.Text()
				case 8:
					extension = columnSelection.Text()
				case 9, 10, 11:
					href, hrefExists := columnSelection.Find("a").Attr("href")
					if hrefExists {
						mirrors = append(mirrors, href)
					}
				}
			})
			books = append(books, Book{
				ID:        id,
				Author:    author,
				Year:      year,
				Title:     title,
				Publisher: publisher,
				Extension: extension,
				Mirrors:   mirrors,
				FileSize:  fileSize,
			})
		}
	})
	return books
}

func getBookDataFromDocumentNew(document *goquery.Document) []Book {
	var books []Book
	document.Find("#tablelibgen > tbody > tr").Each(func(resultsRow int, bookRow *goquery.Selection) {
		var id, author, title, publisher, extension, year, fileSize string
		var mirrors []string
		if resultsRow != 0 {
			bookRow.Find("td").Each(func(column int, columnSelection *goquery.Selection) {
				switch column {
				case 0:
					title = columnSelection.Find("a").First().Text()
				case 1:
					author = columnSelection.Text()
				case 2:
					publisher = columnSelection.Text()
				case 3:
					year = columnSelection.Text()
				case 6:
					fileSize = columnSelection.Text()
				case 7:
					extension = columnSelection.Text()
				case 8:
					columnSelection.Find("a").Each(func(linkCol int, link *goquery.Selection) {
						href, hrefExists := link.Attr("href")
						if hrefExists {
							mirrors = append(mirrors, href)
						}
					})
				}
			})
			books = append(books, Book{
				ID:        id,
				Author:    author,
				Year:      year,
				Title:     title,
				Publisher: publisher,
				Extension: extension,
				Mirrors:   mirrors,
				FileSize:  fileSize,
			})
		}
	})
	return books
}

// Parse HTML and get the data from the table by parsing and iterating through them.
func getBookDataFromDocument(document *goquery.Document, libgenSite Site) []Book {
	switch libgenSite {
	case LibgenOld:
		return getBookDataFromDocumentOld(document)
	case LibgenNew:
		return getBookDataFromDocumentNew(document)
	}
	return []Book{}
}

func getLinkFromDocumentOld(document *goquery.Document) (string, bool) {
	return document.Find("#download > ul > li > a").First().Attr("href")
}

func getLinkFromDocumentNew(document *goquery.Document) (string, bool) {
	return document.Find("#main a").First().Attr("href")
}

func getDirectDownloadLink(link string, libgenType Site) string {
	fmt.Println("Obtaining direct download link")

	resp, err := http.Get(link)

	defer func(Body io.ReadCloser) {
		err := Body.Close()

		if err != nil {
			fmt.Println("Error closing body:", err)
		}
	}(resp.Body)

	if err != nil {
		fmt.Println("Error getting response:", err)
	}

	document, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		fmt.Println("Error creating document:", err)
	}

	var directDownloadLink string
	var relativeLink string
	var exists bool

	switch libgenType {
	case LibgenOld:
		directDownloadLink, exists = getLinkFromDocumentOld(document)
	case LibgenNew:
		u, _ := url.Parse(link)
		// TODO: Add proper err handling
		var host = u.Host
		var protocol = u.Scheme
		relativeLink, exists = getLinkFromDocumentNew(document)

		directDownloadLink = fmt.Sprintf("%s://%s/%s", protocol, host, relativeLink)
	}

	if exists {
		fmt.Println(successColor("Direct download link found"))
		return directDownloadLink
	}

	fmt.Println(errorColor("Direct download link not found"))
	return ""
}

func highlight(s string) string {
	magenta := color.New(color.FgHiWhite).Add(color.BgBlack).SprintFunc()
	return magenta(s)
}

func errorColor(s string) string {
	red := color.New(color.FgRed).SprintFunc()
	return red(s)
}

func infoColor(s string) string {
	yellow := color.New(color.FgYellow).SprintFunc()
	return yellow(s)
}

func successColor(s string) string {
	green := color.New(color.FgHiGreen).SprintFunc()
	return green(s)
}

func searchLibgen(query string, libgenSite Site) (document *goquery.Document, e error) {
	var baseUrl string

	switch libgenSite {
	case LibgenOld:
		baseUrl = "https://libgen.is/search.php"
	case LibgenNew:
		baseUrl = "https://libgen.li/index.php"
	}

	queryString := fmt.Sprintf("%s?req=%s&res=25&view=simple&phrase=1&column=def", baseUrl, url.QueryEscape(query))
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

func SearchBookByTitle(query string, limit int, libgenSite Site) (bookResults []Book, siteToUse Site, e error) {
	fmt.Println("Searching for:", highlight(query))
	var document *goquery.Document
	document, e = searchLibgen(query, siteToUse)

	if e != nil {
		fmt.Println(errorColor("Error searching for book: " + query))
		failedSite := libgenSite
		if failedSite == LibgenOld {
			siteToUse = LibgenNew
		} else {
			siteToUse = LibgenOld
		}

		fmt.Println(infoColor("Retrying with other site"))
		document, e = searchLibgen(query, siteToUse) // If this also fails then we have a problem
	}
	fmt.Println(successColor("Search complete, parsing the document..."))

	bookResults = getBookDataFromDocument(document, siteToUse)

	if len(bookResults) >= limit {
		bookResults = bookResults[:limit]
	}

	return bookResults, siteToUse, e
}

// DownloadSelection Downloads the file to current working directory
func DownloadSelection(selectedBook Book, libgenType Site) {
	link := getDirectDownloadLink(selectedBook.Mirrors[0], libgenType)
	fmt.Println(infoColor("Initializing download "))
	req, _ := http.NewRequest("GET", link, nil)
	resp, error := http.DefaultClient.Do(req)

	if error != nil {
		fmt.Println(errorColor("Error downloading file: " + error.Error()))
	}

	defer resp.Body.Close()
	filename := sanitize.Path(strings.Trim(selectedBook.Title, " ") + "." + selectedBook.Extension)

	f, _ := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0666)
	defer f.Close()

	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		"Downloading",
	)

	bytes, err := io.Copy(io.MultiWriter(f, bar), resp.Body)

	if bytes == 0 || err != nil {
		fmt.Println(bytes, err)
	} else {
		fmt.Println(successColor("File successfully downloaded:"), f.Name())
	}
}
