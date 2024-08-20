package main

import (
	"hmerritt/go-ics-to-markdown/command"
	"hmerritt/go-ics-to-markdown/version"
)

func main() {
	version.PrintTitle()
	command.Run()
}
