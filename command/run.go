package command

import (
	"fmt"
	"hmerritt/go-ics-to-markdown/parse"
	"os"
	"strings"
	"time"

	mdFmt "github.com/shurcooL/markdownfmt/markdown"
)

type RunCommand struct {
	*BaseCommand
}

func (c *RunCommand) Synopsis() string {
	return "Convert ICS file into Markdown table"
}

func (c *RunCommand) Help() string {
	helpText := `
Usage: ics-to-markdown run [options] FILE
  
  Convert ICS file into Markdown table.
`

	return strings.TrimSpace(helpText)
}

func (c *RunCommand) Flags() *FlagMap {
	return GetFlagMap(FlagNamesGlobal)
}

func (c *RunCommand) strictExit() {
	if c.Flags().Get("strict").Value == true {
		c.UI.Error("\nAn error occured while using the '--strict' flag.")
		os.Exit(1)
	}
}

func (c *RunCommand) Run(args []string) int {
	// Record the total duration of this command
	timeStart := time.Now()

	// Initiate error slice and counter
	// Collects errors
	// TODO: make this a type + methods
	// errorSlice := make([]error, 0, 1)
	errorCount := 0

	args = c.Flags().Parse(c.UI, args)

	var icsPath string

	if len(args) == 0 {
		// Use default ICS file
		icsPath = parse.AddICSExtension(parse.ElasticExtension(parse.DefaultICSFileName))
		c.UI.Warn("No file entered.")
		c.strictExit()
		c.UI.Warn("Trying default '" + icsPath + "' instead.\n")
	} else {
		icsPath = parse.ElasticExtension(args[0])
	}

	icsData, err, isURL := parse.FetchICS(icsPath)
	if err != nil {
		if isURL {
			c.UI.Error("Unable to fetch URL data.")
			c.UI.Error(fmt.Sprint(err))
			c.UI.Warn("\nMake sure the link is accessible and try again.")

		} else {
			c.UI.Error("Unable to open file.")
			c.UI.Error(fmt.Sprint(err))
			c.UI.Warn("\nCheck the file is exists and try again.")
			os.Exit(2)
		}
		return 2
	}

	icsEvents, hasEventValue, err := parse.IcsToEvent(icsData)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error parsing ICS file: %v\n", err))
		return 1
	}

	// Print ICS file stats
	c.UI.Output("ICS File")
	c.UI.Output("└── Events       " + fmt.Sprint(len(icsEvents)))
	c.UI.Output("")

	markdownFinal := ""
	markdownTable := parse.ICSEventToMarkdown(icsEvents, hasEventValue)
	markdownFormatted, err := mdFmt.Process("calendar.md", []byte(markdownTable), nil)

	if err == nil {
		markdownFinal = string(markdownFormatted)
	} else {
		markdownFinal = markdownTable
		c.UI.Error(fmt.Sprintf("Error formatting markdown: %v\n", err))
		errorCount++
		c.strictExit()
	}

	err = os.WriteFile("calendar.md", []byte(markdownFinal), 0644)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error writing to file: %v\n", err))
		errorCount++
		c.strictExit()
	}

	if errorCount > 0 {
		c.UI.Warn("Use '--strict' flag to stop immediately if any errors occur\n")

		c.UI.Output(fmt.Sprintf("%s in %s", c.UI.Colorize("ICS file converted (with "+fmt.Sprint(errorCount)+" errors)", c.UI.WarnColor), time.Since(timeStart)))
		return 1
	}

	c.UI.Output(fmt.Sprintf("%s in %s", c.UI.Colorize("ICS file converted", c.UI.SuccessColor), time.Since(timeStart)))

	return 0
}
