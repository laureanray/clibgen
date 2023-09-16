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
	"github.com/manifoldco/promptui"
	"github.com/schollz/progressbar/v3"
)

type Downloader struct {
	selectedBook  book.Book
	directLink    string
	outputFileDir string
	linkType      string
}

func NewDownloader(selectedBook book.Book, directLink string, outputFileDir string) *Downloader {
	return &Downloader{
		selectedBook:  selectedBook,
		directLink:    directLink,
		outputFileDir: outputFileDir,
	}
}

func (d *Downloader) Download() error {
	fmt.Println(console.Info("Initializing download "))

	req, _ := http.NewRequest("GET", d.directLink, nil)
	resp, error := http.DefaultClient.Do(req)

	if error != nil {
		fmt.Println(console.Error("Error downloading file: %s", error.Error()))
	}

	defer resp.Body.Close()
	filename := sanitize.Path(strings.Trim(d.selectedBook.Title, " ") + "." + d.selectedBook.Extension)
	filename = filepath.Clean(d.outputFileDir + "/" + filename)

	fmt.Println("Downloading to: ", filename)

	var f *os.File
	// Check if file exists
	if _, err := os.Stat(filename); err == nil {
		prompt := promptui.Select{
			Label: "File Exists. Do you want to replace file?",
			Items: []string{"Yes", "No"},
		}
		_, result, err := prompt.Run()

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return err
		}

		if result == "Yes" {
			// file doesnt exists
			f, _ = os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0666)
			defer f.Close()
		} else {
			os.Exit(1)	
		}
	} else {
		// file doesnt exists
		f, _ = os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0666)
		defer f.Close()
	}

	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		"Downloading",
	)

	bytes, err := io.Copy(io.MultiWriter(f, bar), resp.Body)

	// Check if byte size is unusually low
	if bytes <= 200 || err != nil {
		fmt.Println(console.Error("File downloaded with unusually low bytes size: %d bytes", bytes))
	} else {
		fmt.Println(console.Success("File successfully downloaded: %s", f.Name()))
	}

	return err
}
