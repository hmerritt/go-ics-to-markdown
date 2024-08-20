package command

import (
	"fmt"
	"hmerritt/go-ics-to-markdown/parse"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/samber/lo"
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
	return GetFlagMap(lo.Union(FlagNamesGlobal, []string{"start", "end"}))
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

	mdPath := "calendar.md"
	if parse.FileExists(icsPath) {
		mdPath = fmt.Sprintf("%s.md", strings.TrimSuffix(filepath.Base(icsPath), ".ics"))
	}

	flagStart := fmt.Sprint(c.Flags().Get("start").Value)
	flagEnd := fmt.Sprint(c.Flags().Get("end").Value)

	filterStart, err := time.Parse("2006-01-02", flagStart)
	if err != nil {
		filterStart = time.Time{}
		if flagStart != "" {
			c.UI.Error("Unable to parse start date.")
		}
	}
	filterEnd, err := time.Parse("2006-01-02", flagEnd)
	if err != nil {
		filterEnd = time.Time{}
		if flagEnd != "" {
			c.UI.Error("Unable to parse end date.")
		}
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

	icsEventsTotal, hasEventValue, err := parse.IcsToEvents(icsData)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error parsing ICS file: %v\n", err))
		return 1
	}

	icsEvents := parse.ICSEventsFilter(icsEventsTotal, parse.ICSEventFilter{
		Start: filterStart,
		End:   filterEnd,
	})

	// Print ICS file stats
	c.UI.Output("ICS File")
	c.UI.Output("├── Events in total       " + fmt.Sprint(len(icsEventsTotal)))
	c.UI.Output("└── Events after filters  " + fmt.Sprint(len(icsEvents)))
	c.UI.Output("")

	markdownFinal := ""
	markdownTable := parse.ICSEventsToMarkdown(icsEvents, hasEventValue)
	markdownFormatted, err := mdFmt.Process(mdPath, []byte(markdownTable), nil)

	if err == nil {
		markdownFinal = string(markdownFormatted)
	} else {
		markdownFinal = markdownTable
		c.UI.Error(fmt.Sprintf("Error formatting markdown: %v\n", err))
		errorCount++
		c.strictExit()
	}

	err = os.WriteFile(mdPath, []byte(markdownFinal), 0644)
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
