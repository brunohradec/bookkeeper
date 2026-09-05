// Reads books from a Calibre library through the Calibre CLI.
package calibre

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

const (
	binary = "calibredb"
	fields = "id,title,authors,cover,pubdate,publisher,tags"

	// Calibre uses 101 as a placeholder year for books with an unknown year.
	placeholderYear = 101
)

// Book holds the metadata of a single book in the Calibre library.
type Book struct {
	ID        int
	Title     string
	Author    string
	CoverPath string // empty if the book has no cover
	Year      int    // zero if the publication date is unknown
	Publisher string
	Tags      []string
}

func (b Book) HasTag(tag string) bool {
	for _, t := range b.Tags {
		if strings.EqualFold(t, tag) {
			return true
		}
	}

	return false
}

func (b Book) HasAnyTag(tags []string) bool {
	for _, tag := range tags {
		if b.HasTag(tag) {
			return true
		}
	}

	return false
}

// Mirrors a single book of the calibredb list JSON output.
type listEntry struct {
	ID        int      `json:"id"`
	Title     string   `json:"title"`
	Authors   string   `json:"authors"`
	Cover     string   `json:"cover"`
	PubDate   string   `json:"pubdate"`
	Publisher string   `json:"publisher"`
	Tags      []string `json:"tags"`
}

// Returns every book of the default Calibre library.
func ListBooks() ([]Book, error) {
	out, err := readFromCLI()
	if err != nil {
		return nil, err
	}

	var entries []listEntry
	if err := json.Unmarshal(out, &entries); err != nil {
		return nil, fmt.Errorf("parsing %s output: %w", binary, err)
	}

	books := make([]Book, len(entries))
	for i, e := range entries {
		books[i] = e.toBook()
	}

	return books, nil
}

// Runs calibredb and returns its raw JSON output.
func readFromCLI() ([]byte, error) {
	path, err := exec.LookPath(binary)
	if err != nil {
		return nil, fmt.Errorf("%s not found in PATH, is Calibre installed?", binary)
	}

	var stdout, stderr bytes.Buffer

	cmd := exec.Command(path, "list", "--for-machine", "--fields", fields)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("running %s: %w: %s", binary, err, strings.TrimSpace(stderr.String()))
	}

	return stdout.Bytes(), nil
}

// Converts a calibredb list entry into a Book.
func (e listEntry) toBook() Book {
	return Book{
		ID:        e.ID,
		Title:     e.Title,
		Author:    e.Authors,
		CoverPath: e.Cover,
		Year:      year(e.PubDate),
		Publisher: e.Publisher,
		Tags:      e.Tags,
	}
}

// Extracts the publication year from a calibredb date,
// returning zero when the date is missing or unknown.
func year(pubDate string) int {
	t, err := time.Parse(time.RFC3339, pubDate)
	if err != nil || t.Year() == placeholderYear {
		return 0
	}

	return t.Year()
}
