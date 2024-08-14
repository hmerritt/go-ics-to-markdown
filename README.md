# Go Calendar ICS -> Markdown

CLI program to convert an ICS calendar file into a markdown table.

[Download binaries here 💾](https://github.com/hmerritt/go-ics-to-markdown/releases)

## Example usage

```bash
$ ics-markdown <path-to-ics>
```

## Developer setup

Setup by running the following bootstrap commands:

```bash
$ go install github.com/magefile/mage
$ mage -v bootstrap
```

Local debug build:

```bash
$ mage -b build:debug
```

Cross-platform release builds (`release` zips up build):

```bash
$ mage -b build:release && mage -b release
```
