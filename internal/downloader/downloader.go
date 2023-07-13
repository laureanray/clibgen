package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/kennygrant/sanitize"
	"github.com/laureanray/clibgen/internal/book"
	"github.com/laureanray/clibgen/internal/console"
	"github.com/schollz/progressbar/v3"
)

type Downloader struct {
  selectedBook book.Book
  directLink string
  outputFileDir string
}

func NewDownloader(selectedBook book.Book, directLink string, outputFileDir string) *Downloader {
  return &Downloader{
    selectedBook: selectedBook,
    directLink: directLink,
    outputFileDir: outputFileDir,
  }
}

func (d *Downloader) Download() error {
	fmt.Println(console.Info("Initializing download "))

  // TODO: implement retry
	req, _ := http.NewRequest("GET", d.directLink, nil)
	resp, error := http.DefaultClient.Do(req)

	if error != nil {
		fmt.Println(console.Error("Error downloading file: %s", error.Error()))
	}

	defer resp.Body.Close()
	filename := sanitize.Path(strings.Trim(d.selectedBook.Title, " ") + "." + d.selectedBook.Extension)
  filename = filepath.Clean(d.outputFileDir + "/" + filename)

  fmt.Println("Downloading to: ", filename)

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
