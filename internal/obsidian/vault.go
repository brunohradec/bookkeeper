// Package obsidian writes Calibre books as notes in an Obsidian vault.
package obsidian

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/brunohradec/bookkeeper/internal/calibre"
)

const (
	// The vault directory name notes are written to.
	booksDir = "Books"

	// The directory name inside booksDir holding cover images.
	coversDir = "covers"

	// The name of the Obsidian base written into booksDir.
	baseName = "Library.base"
)

// The content of the Obsidian base, a card view of the library.
const base = `views:
  - type: cards
    name: Bookshelf
    filters:
      and:
        - file.inFolder("Books")
        - file.ext == "md"
    order:
      - file.name
      - author
      - status
    image: note.cover
    imageAspectRatio: 1.55
`

// An Obsidian vault book notes are written to.
type Vault struct {
	path string
}

// Opens the Obsidian vault at path and creates the directories notes
// and cover images are written to. It fails if the vault directory
// does not exist.
func OpenVault(path string) (Vault, error) {
	info, err := os.Stat(path)
	if err != nil {
		return Vault{}, fmt.Errorf("vault %q: %w", path, err)
	}

	if !info.IsDir() {
		return Vault{}, fmt.Errorf("vault %q is not a directory", path)
	}

	vault := Vault{path: path}
	if err := os.MkdirAll(vault.coversPath(), 0o755); err != nil {
		return Vault{}, err
	}

	return vault, nil
}

// Creates the note and cover image for a book, reporting whether
// anything was written. Books whose note already exists are skipped.
func (v Vault) Write(b calibre.Book) (written bool, err error) {
	notePath := filepath.Join(v.booksPath(), NoteName(b)+".md")

	found, err := exists(notePath)
	if err != nil {
		return false, err
	}

	if found {
		return false, nil
	}

	if b.CoverPath != "" {
		dst := filepath.Join(v.coversPath(), coverName(b))
		if err := copyFile(b.CoverPath, dst); err != nil {
			return false, fmt.Errorf("copying cover: %w", err)
		}
	}

	if err := os.WriteFile(notePath, []byte(frontmatter(b)), 0o644); err != nil {
		return false, fmt.Errorf("writing note: %w", err)
	}

	return true, nil
}

// Creates the Obsidian base for the library, reporting whether it was
// written. An existing base is left untouched.
func (v Vault) WriteBase() (written bool, err error) {
	basePath := filepath.Join(v.booksPath(), baseName)

	found, err := exists(basePath)
	if err != nil {
		return false, err
	}

	if found {
		return false, nil
	}

	if err := os.WriteFile(basePath, []byte(base), 0o644); err != nil {
		return false, fmt.Errorf("writing base: %w", err)
	}

	return true, nil
}

// Reports whether a file exists.
func exists(path string) (bool, error) {
	switch _, err := os.Stat(path); {
	case err == nil:
		return true, nil
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	default:
		return false, err
	}
}

func (v Vault) booksPath() string {
	return filepath.Join(v.path, booksDir)
}

func (v Vault) coversPath() string {
	return filepath.Join(v.booksPath(), coversDir)
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	return os.WriteFile(dst, data, 0o644)
}
