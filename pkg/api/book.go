package api

type Book struct {
	ID        string   `json:"ID"`
	ISBN      string   `json:"ISBN"`
	Title     string   `json:"title"`
	Author    string   `json:"author"`
	Year      string   `json:"year"`
	Publisher string   `json:"publisher"`
	Extension string   `json:"extension"`
	Mirrors   []string `json:"mirrors"`
}
