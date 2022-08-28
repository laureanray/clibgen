package api

type SciMag struct {
	ID        string
	Title     string
	Author    string
	Year      string
	Journal   string
	Extension string
	Mirrors   []string
	FileSize  string
}
