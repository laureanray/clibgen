package downloader

import "github.com/laureanray/clibgen/internal/book"

// DownloadSelection Downloads the file to current working directory
func DownloadSelection(selectedBook book.Book) {
  // TODO: implement fallback, and retry here
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
