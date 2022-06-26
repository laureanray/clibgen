package api

type Book struct {
	ID        string
	Title     string
	Author    string
	Year      string
	Publisher string
	Extension string
	Mirrors   []string
	FileSize  string
}
