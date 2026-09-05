# bookkeeper

Generates an Obsidian note for every book in a Calibre library.

Each note is named `Book (Author).md`, holds the book's metadata as YAML
frontmatter, and gets its cover image copied into the vault. Book data is read
through the Calibre command line tool, never from the Calibre database.

## Requirements

- Calibre, with `calibredb` on your `PATH`

## How to use

Close the Calibre desktop app first, as it locks the library.

```sh
bookkeeper --vault ~/Documents/Notes
```

or from this repo using `go run`:
```sh
go run cmd/bookkeeper/main.go --vault ~/Documents/Notes
```

The `--vault` flag is required and must point at an existing Obsidian vault.

To leave out books carrying a particular Calibre tag, pass `--exclude-tag`:

```sh
bookkeeper --vault ~/Documents/Notes --exclude-tag notmine
```

The tag is matched ignoring case.

## What it writes

Notes go into `Books/` inside the vault and cover images into `Books/covers/`.
Both are created if missing, along with an Obsidian base named `Library.base`
that shows the library in a card view.

Existing notes and an existing base are never overwritten, so the program is
safe to rerun after adding books to Calibre.

## Build

```sh
go build -o bookkeeper ./cmd/bookkeeper
```

To install it on your `PATH` instead, run `go install ./cmd/bookkeeper`, which
puts the binary in `~/go/bin`.
