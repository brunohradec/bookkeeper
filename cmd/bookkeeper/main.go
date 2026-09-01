// Command bookkeeper generates Obsidian notes from a Calibre library.
package main

import (
	"fmt"
	"os"

	"github.com/brunohradec/bookkeeper/internal/book"
	"github.com/brunohradec/bookkeeper/internal/calibre"
)

func main() {
	var source book.Source = calibre.Reader{}

	books, err := source.ListBooks()
	if err != nil {
		fmt.Fprintln(os.Stderr, "bookkeeper:", err)
		os.Exit(1)
	}

	for _, b := range books {
		fmt.Printf("%s (%s)\n", b.Title, b.Author)
		fmt.Printf("  year:      %d\n", b.Year)
		fmt.Printf("  publisher: %s\n", b.Publisher)
		fmt.Printf("  cover:     %s\n", b.CoverPath)
	}

	fmt.Printf("\n%d books\n", len(books))
}
