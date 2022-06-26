package api

import "testing"

// TODO: Use mocks, introduce better assertions
func TestSearchBookByTitle(t *testing.T) {
	books, err := SearchBookByTitle("Test", 5)

	if err != nil {
		t.Errorf("err")
	}

	if len(books) == 0 {
		t.Errorf("got 0 books")
	}
}

// TODO: Use mocks
func TestDownloadSelection(t *testing.T) {
	selectedBook := Book{
		ID:        "ID",
		Extension: ".epub",
		Title:     "Elon Musk",
		Mirrors: []string{
			"http://library.lol/main/580000A1CAA698C2EFD8F5439E9A1F26",
			"http://library.lol/main/2F2DBA2A621B693BB95601C16ED680F8",
			"http://library.lol/main/2F2DBA2A621B693BB95601C16ED680F8",
		},
		Author:    "Vance",
		Year:      "2012",
		Publisher: "Someone",
	}

	DownloadSelection(selectedBook)
}
