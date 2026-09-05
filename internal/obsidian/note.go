package obsidian

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/brunohradec/bookkeeper/internal/calibre"
)

const (
	// Tag every book note gets, physical books included.
	bookTag = "book"

	// Tag marking the notes generated from the Calibre library.
	ebookTag = "ebook"

	// The default status of the book is "unknown" as we cannot
	// know if the book was read or not by looking at just the
	// Calibre library.
	defaultStatus = "unknown"
)

func NoteName(b calibre.Book) string {
	title, _, _ := strings.Cut(b.Title, ":")
	return fmt.Sprintf("%s (%s)", strings.TrimSpace(title), b.Author)
}

func coverName(b calibre.Book) string {
	return NoteName(b) + filepath.Ext(b.CoverPath)
}

func frontmatter(b calibre.Book) string {
	var sb strings.Builder

	sb.WriteString("---\n")
	sb.WriteString(fmField("book", b.Title))
	sb.WriteString(fmListField("author", splitAuthors(b.Author)))
	sb.WriteString(fmField("cover", coverLink(b)))
	sb.WriteString(fmYearField(b.Year))
	sb.WriteString(fmField("publisher", b.Publisher))
	sb.WriteString(fmListField("tags", []string{bookTag, ebookTag}))
	sb.WriteString(fmField("status", defaultStatus))
	sb.WriteString("---\n")

	return sb.String()
}

func fmField(key, value string) string {
	if value == "" {
		return key + ":\n"
	}
	return key + ": " + fmQuote(value) + "\n"
}

func fmListField(key string, values []string) string {
	var sb strings.Builder

	sb.WriteString(key)
	sb.WriteString(":\n")

	for _, value := range values {
		sb.WriteString("  - ")
		sb.WriteString(fmQuote(value))
		sb.WriteString("\n")
	}

	return sb.String()
}

func fmYearField(year int) string {
	if year == 0 {
		return "year published:\n"
	}

	return fmt.Sprintf("year published: %d\n", year)
}

func coverLink(b calibre.Book) string {
	if b.CoverPath == "" {
		return ""
	}

	return fmt.Sprintf("[[%s/%s]]", coversDir, coverName(b))
}

func splitAuthors(joined string) []string {
	authors := strings.Split(joined, "&")
	for i, author := range authors {
		authors[i] = strings.TrimSpace(author)
	}

	return authors
}

// Wraps a value in double quotes when YAML would otherwise misread it
func fmQuote(value string) string {
	// The characters listed here are the ones YAML reads as syntax
	// rather than as part of a plain value, so they require quotes.
	if !strings.ContainsAny(value, ":#\"'[]{}|>&*!%@`\n") && !looksNumeric(value) {
		return value
	}

	return `"` + strings.NewReplacer(`\`, `\\`, `"`, `\"`).Replace(value) + `"`
}

// Reports whether a value would be read as a number rather than as text.
func looksNumeric(value string) bool {
	return strings.IndexFunc(value, func(r rune) bool {
		return r != '.' && (r < '0' || r > '9')
	}) < 0
}
