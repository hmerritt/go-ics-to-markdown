# Go Calendar ICS -> Markdown

CLI program to convert an ICS calendar file into a markdown table.

[Download binaries here 💾](https://github.com/hmerritt/go-ics-to-markdown/releases)

## Example usage

Show help:

```bash
$ ics-to-markdown
```

Convert ICS file to markdown:

```bash
$ ics-to-markdown run <path-to-ics>
```

## Developer setup

Setup by running the following bootstrap commands:

```bash
$ go install github.com/magefile/mage
$ mage -v bootstrap
```

Local debug build:

```bash
$ mage -v build:debug
```

Cross-platform release builds (`release` zips up build):

```bash
$ mage -v build:release && mage -v release
```
