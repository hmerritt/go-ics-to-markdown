package command

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

type ListCommand struct {
	*BaseCommand
}

func (c *ListCommand) Synopsis() string {
	return "List structure files in the current directory"
}

func (c *ListCommand) Help() string {
	helpText := `
Usage: ics-to-markdown list [options] PATH
  
  Lists all ICS files in the current directory.
`

	return strings.TrimSpace(helpText)
}

func (c *ListCommand) Run(args []string) int {
	path := "./"

	if len(args) > 0 {
		path = args[0]
	}

	files, err := os.ReadDir(path)
	if err != nil {
		c.UI.Error("Unable to read directory files")
		c.UI.Error(fmt.Sprint(err))
		c.UI.Warn("\nThis is most likely due to a lack of permissions,")
		c.UI.Warn("check you have (at least) read access to this directory.")
		return 1
	}

	icsFiles := make([]string, 0)

	for _, f := range files {
		name := f.Name()
		if !f.IsDir() {
			if strings.HasSuffix(name, ".ics") {
				icsFiles = append(icsFiles, name)
			}
		}
	}

	sort.Strings(icsFiles)

	if len(icsFiles) == 0 {
		// Get friendly path name.
		// Converts './' to 'current'
		pathFriendly := "'" + path + "'"
		if path == "./" {
			pathFriendly = "current"
		}

		c.UI.Output("No ICS files found in the " + pathFriendly + " directory.")
		return 2
	}

	fmt.Printf("Found %d ICS files:\n", len(icsFiles))

	for _, yf := range icsFiles {
		fmt.Printf("-- %s\n", yf)
	}

	c.UI.Output("\nConvert ics->markdown file using:")
	c.UI.Output("$ ics-to-markdown run <FILE>")

	return 0
}
