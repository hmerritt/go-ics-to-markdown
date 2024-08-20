package parse

import (
	"fmt"
	"path/filepath"
	"strings"
)

var (
	ICSExtensions      = [1]string{".ics"}
	DefaultICSFileName = "calendar"
	DefaultICSFile     = fmt.Sprintf("%s%s", DefaultICSFileName, ICSExtensions[0])
)

func FileExtension(filename string) string {
	return filepath.Ext(filename)
}

func IsICSFile(filename string) bool {
	for _, icsExt := range ICSExtensions {
		if strings.HasSuffix(filename, icsExt) {
			return true
		}
	}
	return false
}

func AddICSExtension(filename string) string {
	if IsICSFile(filename) {
		return filename
	}

	return filename + ICSExtensions[0]
}

func ElasticExtension(filename string) string {
	if IsICSFile(filename) || FileExists(filename) {
		return filename
	}

	// Check for an existing file using all ics extensions
	for icsExtIndex := range ICSExtensions {
		elastic := filename + ICSExtensions[icsExtIndex]
		if FileExists(elastic) {
			return elastic
		}
	}

	// Return original filename if nothing found
	return filename
}
