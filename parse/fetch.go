package parse

import (
	"errors"
	"os"
	"strings"

	"hmerritt/go-ics-to-markdown/ui"

	"github.com/imroc/req"
)

func IsUrl(path string) bool {
	return (strings.HasPrefix(path, "www.") || strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://"))
}

func UseUrl(path string) bool {
	return (!FileExists(path) && IsUrl(path))
}

func FileExists(path string) bool {
	stat, err := os.Stat(path)
	return err == nil && !stat.IsDir()
}

func FetchFile(filepath string) ([]byte, error) {
	data, err := os.ReadFile(filepath)
	return data, err
}

func FetchUrl(url string) ([]byte, error) {
	res, err := req.Get(url)
	if err != nil {
		return nil, err
	}

	resStatus := res.Response().Status

	if !strings.HasPrefix(resStatus, "1") && !strings.HasPrefix(resStatus, "2") {
		return nil, errors.New("request returned a bad http status code: " + resStatus + ".")
	}

	return res.Bytes(), nil
}

// Fetch and parse ICS file locally or from a URL
func FetchICS(path string) ([]byte, error, bool) {
	var data []byte
	var err error
	var isURL = false

	// Decide if URL or file
	if UseUrl(path) {
		isURL = true
		ui.Spinner.Start("", " Fetching URL data...")

		data, err = FetchUrl(path)

		ui.Spinner.Stop()

		if err != nil {
			return nil, err, isURL
		}
	} else {
		data, err = FetchFile(path)

		if err != nil {
			return nil, err, isURL
		}
	}

	return data, nil, isURL
}
