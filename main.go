package main

import (
	"fmt"
	"hmerritt/go-ics-to-markdown/version"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	md "github.com/JohannesKaufmann/html-to-markdown"
	ics "github.com/arran4/golang-ical"
	mdFmt "github.com/shurcooL/markdownfmt/markdown"
)

type Event struct {
	Summary     string
	Start       time.Time
	End         time.Time
	Description string
	Location    string
}

func main() {
	version.PrintTitle()

	if len(os.Args) < 2 {
		fmt.Println("Please provide the path to the ICS file")
		return
	}

	icsFilePath := os.Args[1]
	markdownTable := convertIcsToMarkdown(icsFilePath)
	markdownFormatted, err := mdFmt.Process("calendar.md", []byte(markdownTable), nil)
	if err != nil {
		fmt.Printf("Error formatting markdown: %v\n", err)
	}

	fmt.Println(string(markdownFormatted))

	err = os.WriteFile("calendar.md", markdownFormatted, 0644)
	if err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
	}
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

			events = append(events, Event{
				Summary:     summary,
				Start:       start,
				End:         end,
				Location:    location,
				Description: removeLineBreaks(description),
			})
		}
	}

	sort.Slice(events, func(i, j int) bool {
		return events[i].Start.Before(events[j].Start)
	})

	markdown := "| Date | Time | Location | Event | Description |\n"
	markdown += "|------|------|----------|-------|-------------|\n"

	for _, event := range events {
		date := event.Start.Format("2006-01-02")
		startTime := event.Start.Format("15:04")
		endTime := event.End.Format("15:04")
		markdown += fmt.Sprintf("| %s | %s-%s | %s | %s | %s |\n",
			date, startTime, endTime, event.Location, event.Summary, event.Description)
	}

	return markdown
}

func removeLineBreaks(text string) string {
	re := regexp.MustCompile(`\x{000D}\x{000A}|[\x{000A}\x{000B}\x{000C}\x{000D}\x{0085}\x{2028}\x{2029}]`)
	return re.ReplaceAllString(text, `<br>`)
}
