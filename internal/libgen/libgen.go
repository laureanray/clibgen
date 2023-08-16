package libgen

type Domain string

const (
	RS = "rs"
	IS = "is"
	ST = "st"
)

type Filter string

const (
	TITLE     = "title"
	AUTHOR    = "author"
	SERIES    = "series"
	PUBLISHER = "publisher"
	YEAR      = "year"
	ISBN      = "isbn"
	MD5       = "md5"
	TAGS      = "tags"
	EXTENSION = "extension"
)
