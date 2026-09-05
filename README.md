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

The `--vault` flag is required and must point at an existing Obsidian vault.

## Build

```sh
go build -o bookkeeper ./cmd/bookkeeper
```

To install it on your `PATH` instead, run `go install ./cmd/bookkeeper`, which
puts the binary in `~/go/bin`.
