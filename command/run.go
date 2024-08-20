package command

import (
	"fmt"
	"hmerritt/go-ics-to-markdown/parse"
	"hmerritt/go-ics-to-markdown/ui"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	md "github.com/JohannesKaufmann/html-to-markdown"
	ics "github.com/arran4/golang-ical"
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

	var icsFilePath string

	if len(args) == 0 {
		// Use default ICS file
		icsFilePath = parse.AddICSExtension(parse.ElasticExtension(parse.DefaultICSFileName))
		c.UI.Warn("No file entered.")
		c.strictExit()
		c.UI.Warn("Trying default '" + icsFilePath + "' instead.\n")
	} else {
		icsFilePath = parse.ElasticExtension(args[0])
	}

	ui.Spinner.Start("", " Running...")

	markdownTable := convertIcsToMarkdown(icsFilePath)
	markdownFormatted, err := mdFmt.Process("calendar.md", []byte(markdownTable), nil)
	if err != nil {
		ui.Spinner.Stop()
		c.UI.Error(fmt.Sprintf("Error formatting markdown: %v\n", err))
		errorCount++
		c.strictExit()
	}

	ui.Spinner.Stop()
	fmt.Println(string(markdownFormatted))

	err = os.WriteFile("calendar.md", markdownFormatted, 0644)
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

type Event struct {
	Summary     string
	Start       time.Time
	End         time.Time
	Description string
	Location    string
}

func convertIcsToMarkdown(filePath string) string {
	icsData, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return ""
	}

	calendar, err := ics.ParseCalendar(strings.NewReader(string(icsData)))
	if err != nil {
		fmt.Printf("Error parsing ICS data: %v\n", err)
		return ""
	}

	htmlToMd := md.NewConverter("", true, nil)

	hasLocation := false

	var events []Event
	for _, component := range calendar.Components {
		if event, ok := component.(*ics.VEvent); ok {
			start, _ := event.GetStartAt()
			end, _ := event.GetEndAt()
			summary := ""
			description := ""
			location := ""

			if summaryProp := event.GetProperty(ics.ComponentPropertySummary); summaryProp != nil {
				summary = summaryProp.Value
			}

			if descProp := event.GetProperty(ics.ComponentPropertyDescription); descProp != nil {
				description = descProp.Value
			}
			markdown, err := htmlToMd.ConvertString(description)
			if err == nil {
				description = markdown
			}

			if locationProp := event.GetProperty(ics.ComponentPropertyLocation); locationProp != nil {
				location = locationProp.Value
			}
			if location != "" {
				hasLocation = true
			}

			events = append(events, Event{
				Summary:     summary,
				Start:       start,
				End:         end,
				Location:    location,
				Description: convertLineBreaks(description),
			})
		}
	}

	sort.Slice(events, func(i, j int) bool {
		return events[i].Start.Before(events[j].Start)
	})

	markdown := ""

	if hasLocation {
		markdown = "| Date | Time | Location | Event | Description |\n"
		markdown += "|------|------|----------|-------|-------------|\n"
	} else {
		markdown = "| Date | Time | Event | Description |\n"
		markdown += "|------|------|-------|-------------|\n"
	}

	for _, event := range events {
		date := event.Start.Format("2006-01-02")
		startTime := event.Start.Format("15:04")
		endTime := event.End.Format("15:04")
		if hasLocation {
			markdown += fmt.Sprintf("| %s | %s-%s | %s | %s | %s |\n",
				date, startTime, endTime, event.Location, event.Summary, event.Description)
		} else {
			markdown += fmt.Sprintf("| %s | %s-%s |  %s | %s |\n",
				date, startTime, endTime, event.Summary, event.Description)
		}
	}

	return markdown
}

func convertLineBreaks(text string) string {
	re := regexp.MustCompile(`\x{000D}\x{000A}|[\x{000A}\x{000B}\x{000C}\x{000D}\x{0085}\x{2028}\x{2029}]`)
	return re.ReplaceAllString(text, `<br>`)
}
