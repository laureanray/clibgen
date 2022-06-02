package main

import (
	"fmt"
	"genread-api/pkg/api"
	"log"
)

func main() {
	books, err := api.SearchBookByTitle("Elon Musk")

	for _, book := range books {
		fmt.Println("Title: ", book.Title)
	}

	if err != nil {
		log.Fatal(err)
	}

	//fmt.Println(books)
}

//
//type Book struct {
//	ID        string `json:"id"`
//	ISBN      string `json:"isbn"`
//	Title     string `json:"title"`
//	Author    string `json:"author"`
//	Year      string `json:"year"`
//	Publisher string `json:"publisher"`
//}

//func getBook(c *gin.Context) {
//	title := c.Param("title")
//	queryString := fmt.Sprintf("https://libgen.li/index.php?req=%s&res=25&filesuns=all", title)
//	resp, err := http.Get(queryString)
//
//	defer resp.Body.Close()
//
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	doc, err := goquery.NewDocumentFromReader(resp.Body)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	var books []Book
//
//	doc.Find(".table-striped  > tbody > tr").Each(func(i int, s *goquery.Selection) {
//		var id, author, title string
//
//		s.Find("td").Each(func(column int, s *goquery.Selection) {
//			switch column {
//			case 0:
//				title = s.Find("a").Text()
//				fmt.Println("a: ", title)
//
//			case 1:
//				author = s.Text()
//			case 2:
//				title = s.Find(".badge-secondary").Text()
//			}
//		})
//		books = append(books, Book{
//			ID:     id,
//			Author: author,
//			Title:  title,
//		})
//
//	})
//	c.JSON(200, gin.H{
//		"message": "test",
//	})
//}
