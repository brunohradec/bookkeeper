// Defines the book domain model and the sources books are read from.
package book

// Metadata of the single book, independent of the source.
type Book struct {
	ID        int
	Title     string
	Author    string
	CoverPath string // empty if the book has no cover
	Year      int    // zero if the publication date is unknown
	Publisher string
}

// Source is a library books can be read from, such as Calibre, Goodreads, Hardcover, etc.
type Source interface {
	ListBooks() ([]Book, error)
}
