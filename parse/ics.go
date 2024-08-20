package parse

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	md "github.com/JohannesKaufmann/html-to-markdown"
	ics "github.com/arran4/golang-ical"
	"github.com/samber/lo"
)

type ICSEvent struct {
	Summary     string
	Start       time.Time
	End         time.Time
	Description string
	Location    string
}

type ICSEventFilter struct {
	Start time.Time
	End   time.Time
}

func IcsToEvents(icsData []byte) ([]ICSEvent, map[string]bool, error) {
	calendar, err := ics.ParseCalendar(strings.NewReader(string(icsData)))
	if err != nil {
		return nil, nil, err
	}

	htmlToMd := md.NewConverter("", true, nil)

	hasEventValue := map[string]bool{
		"start":       true,
		"end":         true,
		"summary":     false,
		"description": false,
		"location":    false,
	}

	var events []ICSEvent
	for _, component := range calendar.Components {
		if event, ok := component.(*ics.VEvent); ok {
			start, _ := event.GetStartAt()
			end, _ := event.GetEndAt()
			summary := ""
			description := ""
			location := ""

			if summaryProp := event.GetProperty(ics.ComponentPropertySummary); summaryProp != nil && summaryProp.Value != "" {
				summary = summaryProp.Value
				hasEventValue["summary"] = true
			}

			if descProp := event.GetProperty(ics.ComponentPropertyDescription); descProp != nil && descProp.Value != "" {
				description = descProp.Value
				hasEventValue["description"] = true
			}
			markdown, err := htmlToMd.ConvertString(description)
			if err == nil {
				description = markdown
			}

			if locationProp := event.GetProperty(ics.ComponentPropertyLocation); locationProp != nil && locationProp.Value != "" {
				location = locationProp.Value
				hasEventValue["location"] = true
			}

			events = append(events, ICSEvent{
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

	return events, hasEventValue, nil
}

func ICSEventsFilter(events []ICSEvent, filter ICSEventFilter) []ICSEvent {
	switch true {
	case !filter.Start.IsZero() && !filter.End.IsZero():
		// Filter from start AND up-to end time
		return lo.Filter(events, func(e ICSEvent, index int) bool {
			return (e.Start.After(filter.Start) || e.Start.Equal(filter.Start)) &&
				(e.End.Before(filter.End) || e.End.Equal(filter.End))
		})
	case !filter.Start.IsZero():
		// Filter from start time
		return lo.Filter(events, func(e ICSEvent, index int) bool {
			return e.Start.After(filter.Start) || e.Start.Equal(filter.Start)
		})
	case !filter.End.IsZero():
		// Filter up-to end time
		return lo.Filter(events, func(e ICSEvent, index int) bool {
			return e.End.Before(filter.End) || e.End.Equal(filter.End)
		})
	}

	return events
}

func ICSEventsToMarkdown(events []ICSEvent, hasEventValue map[string]bool) string {
	var headerFields, separatorFields []string
	fieldOrder := []string{"Date", "Time", "Location", "Event", "Description"}

	for _, field := range fieldOrder {
		switch field {
		case "Date", "Time":
			if hasEventValue["start"] || hasEventValue["end"] {
				headerFields = append(headerFields, field)
				separatorFields = append(separatorFields, strings.Repeat("-", len(field)))
			}
		case "Location":
			if hasEventValue["location"] {
				headerFields = append(headerFields, field)
				separatorFields = append(separatorFields, strings.Repeat("-", len(field)))
			}
		case "Event":
			if hasEventValue["summary"] {
				headerFields = append(headerFields, field)
				separatorFields = append(separatorFields, strings.Repeat("-", len(field)))
			}
		case "Description":
			if hasEventValue["description"] {
				headerFields = append(headerFields, field)
				separatorFields = append(separatorFields, strings.Repeat("-", len(field)))
			}
		}
	}

	markdown := fmt.Sprintf("| %s |\n", strings.Join(headerFields, " | "))
	markdown += fmt.Sprintf("| %s |\n", strings.Join(separatorFields, " | "))

	for _, event := range events {
		var rowFields []string
		if hasEventValue["start"] || hasEventValue["end"] {
			date := event.Start.Format("2006-01-02")
			rowFields = append(rowFields, date)

			startTime := event.Start.Format("15:04")
			endTime := event.End.Format("15:04")
			rowFields = append(rowFields, fmt.Sprintf("%s-%s", startTime, endTime))
		}
		if hasEventValue["location"] {
			rowFields = append(rowFields, event.Location)
		}
		if hasEventValue["summary"] {
			rowFields = append(rowFields, event.Summary)
		}
		if hasEventValue["description"] {
			rowFields = append(rowFields, event.Description)
		}
		markdown += fmt.Sprintf("| %s |\n", strings.Join(rowFields, " | "))
	}

	return markdown
}

func convertLineBreaks(text string) string {
	re := regexp.MustCompile(`\x{000D}\x{000A}|[\x{000A}\x{000B}\x{000C}\x{000D}\x{0085}\x{2028}\x{2029}]`)
	return re.ReplaceAllString(text, `<br>`)
}
