package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/kennygrant/sanitize"
	"github.com/laureanray/clibgen/internal/book"
	"github.com/laureanray/clibgen/internal/console"
	"github.com/schollz/progressbar/v3"
)

type Downloader struct {
  selectedBook book.Book
}

func NewDownloader(selectedBook book.Book) *Downloader {
  return &Downloader{
    selectedBook: selectedBook,
  }
}

func (d *Downloader) Download() error {
	fmt.Println(console.Info("Initializing download "))

  // TODO: implement retry
	req, _ := http.NewRequest("GET", d.selectedBook.Mirrors[0], nil)
	resp, error := http.DefaultClient.Do(req)

	if error != nil {
		fmt.Println(console.Error("Error downloading file: %s", error.Error()))
	}

	defer resp.Body.Close()
	filename := sanitize.Path(strings.Trim(d.selectedBook.Title, " ") + "." + d.selectedBook.Extension)

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
		fmt.Println(console.Success("File successfully downloaded: %s", f.Name()))
	}

  return err
}
