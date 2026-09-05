// Command bookkeeper generates Obsidian notes from a Calibre library.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/brunohradec/bookkeeper/internal/calibre"
	"github.com/brunohradec/bookkeeper/internal/obsidian"
)

func main() {
	path := flag.String("vault", "", "path to the Obsidian vault (required)")
	flag.Parse()

	if *path == "" {
		fmt.Fprintln(os.Stderr, "bookkeeper: the -vault flag is required")
		os.Exit(1)
	}

	vault, err := obsidian.OpenVault(*path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "bookkeeper:", err)
		os.Exit(1)
	}

	books, err := calibre.ListBooks()
	if err != nil {
		fmt.Fprintln(os.Stderr, "bookkeeper:", err)
		os.Exit(1)
	}

	var written, skipped int
	for _, b := range books {
		ok, err := vault.Write(b)
		if err != nil {
			fmt.Fprintf(os.Stderr, "bookkeeper: %s: %v\n", b.Title, err)
			os.Exit(1)
		}

		if ok {
			written++
			fmt.Println("+", obsidian.NoteName(b))
		} else {
			skipped++
		}
	}

	fmt.Printf("\n%d written, %d skipped\n", written, skipped)
}
