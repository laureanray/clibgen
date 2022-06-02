package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func main() {
	r := gin.Default()
	r.GET("/api/book/:title", getBook)
	r.Run() // listen and serve on 0.0.0.0:8080
}

type Book struct {
	ID        string `json:"id"`
	ISBN      string `json:"isbn"`
	Title     string `json:"title"`
	Author    string `json:"author"`
	Year      string `json:"year"`
	Publisher string `json:"publisher"`
}

func getBook(c *gin.Context) {
	title := c.Param("title")
	queryString := fmt.Sprintf("https://libgen.li/index.php?req=%s&res=25&filesuns=all", title)
	resp, err := http.Get(queryString)
	fmt.Println("query: ", queryString)
	defer resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var books []Book
	// Find the review items
	doc.Find(".table-striped  > tbody > tr").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title

		//html, e := s.Html()
		//
		//if e != nil {
		//	log.Fatal(e)
		//}
		//fmt.Println("s: ", html)

		var id, author, title string

		s.Find("td").Each(func(column int, s *goquery.Selection) {
			//html, e := s.Html()
			//
			//if e != nil {
			//	log.Fatal(e)
			//}
			//fmt.Println("td: ", html)

			//fmt.Println(s.Text())
			switch column {
			case 0:
				title = s.Find("a").Text()
				fmt.Println("a: ", title)

			case 1:
				author = s.Text()
			case 2:
				title = s.Find(".badge-secondary").Text()

			}
		})

		fmt.Println("Appending> ", id, author, title)

		books = append(books, Book{
			ID:     id,
			Author: author,
			Title:  title,
		})

	})

	fmt.Println(books)

	//fmt.Println("query:", queryString)
	//fmt.Println("data:", responseString)
	//fmt.Println(resp.Body)
	//fmt.Println(err)
	c.JSON(200, gin.H{
		"message": "test",
	})
}
