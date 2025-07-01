package internal

type Book struct {
	Title     string
	Author    string
	Year      string
	Size      string
	Extension string
	Mirrors   map[string]string // map[mirrorName]downloadPageURL
} 