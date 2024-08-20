package command

import (
	"fmt"
	"os"

	"hmerritt/go-ics-to-markdown/version"

	"github.com/mitchellh/cli"
)

func Run() {
	// Initiate new CLI app
	app := cli.NewCLI("ics-to-markdown", version.GetVersion().VersionNumber())
	app.Args = os.Args[1:]

	// Feed active commands to CLI app
	app.Commands = map[string]cli.CommandFactory{
		"list": func() (cli.Command, error) {
			return &ListCommand{
				BaseCommand: GetBaseCommand(),
			}, nil
		},
		"run": func() (cli.Command, error) {
			return &RunCommand{
				BaseCommand: GetBaseCommand(),
			}, nil
		},
	}

	// Run app
	exitStatus, err := app.Run()
	if err != nil {
		os.Stderr.WriteString(fmt.Sprint(err))
	}

	// Exit without an error if no arguments were passed
	if len(app.Args) == 0 {
		os.Exit(0)
	}

	os.Exit(exitStatus)
}
