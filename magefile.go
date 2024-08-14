//go:build mage

package main

import (
	"os"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type Build mg.Namespace
type Release mg.Namespace

var Aliases = map[string]interface{}{
	"build": Build.Release,
}

// ----------------------------------------------------------------------------
// Test
// ----------------------------------------------------------------------------

// Runs Go tests
func Test() error {
	log := NewLogger()
	defer log.End()
	return RunSync([][]string{
		{"gotestsum", "--format", "pkgname", "--", "--cover", "./..."},
	})
}

func Bench() error {
	log := NewLogger()
	defer log.End()
	return RunSync([][]string{
		{"gotestsum", "--format", "pkgname", "--", "--cover", "-bench", ".", "-benchmem", "./..."},
	})
}

// ----------------------------------------------------------------------------
// Build
// ----------------------------------------------------------------------------

func (Build) Debug() error {
	log := NewLogger()
	defer log.End()
	log.Info("compiling debug binary")
	return RunSync([][]string{
		{"go", "build", "-ldflags", "-s -w", "."},
	})
}

func (Build) Release() error {
	log := NewLogger()
	defer log.End()
	log.Info("compiling release binary")
	return RunSync([][]string{
		{"gox",
			"-osarch",
			"linux/amd64 windows/amd64",
			"-gocmd",
			"go",
			"-ldflags",
			LdFlagString(),
			"-tags",
			"ics-to-markdown",
			"-output",
			"bin/ics-to-markdown",
			"."},
	})
}

// ----------------------------------------------------------------------------
// Bootstrap
// ----------------------------------------------------------------------------

// Bootstraps required packages (installs required linux/macOS packages if needed)
func Bootstrap() error {
	log := NewLogger()
	defer log.End()

	// Install mage bootstrap (the recommended, as seen in https://magefile.org)
	if !ExecExists("mage") && ExecExists("git") {
		log.Info("installing mage")
		tmpDir := "__tmp_mage"

		if err := sh.Run("git", "clone", "https://github.com/magefile/mage", tmpDir); err != nil {
			return log.Error("error: installing mage: ", err)
		}

		if err := os.Chdir(tmpDir); err != nil {
			return log.Error("error: installing mage: ", err)
		}

		if err := sh.Run("go", "run", "bootstrap.go"); err != nil {
			return log.Error("error: installing mage: ", err)
		}

		if err := os.Chdir("../"); err != nil {
			return log.Error("error: installing mage: ", err)
		}

		os.RemoveAll(tmpDir)
	}

	// Install Go dependencies
	log.Info("installing go dependencies")
	return RunSync([][]string{
		{"go", "mod", "vendor"},
		{"go", "mod", "tidy"},
		{"go", "generate", "-tags", "tools", "tools/tools.go"},
	})
}

// ----------------------------------------------------------------------------
// Housekeeping
// ----------------------------------------------------------------------------

// Update all Go dependencies
func UpdateDeps() error {
	log := NewLogger()
	defer log.End()
	return RunSync([][]string{
		{"go", "get", "-u", "all"},
		{"go", "mod", "vendor"},
		{"go", "mod", "tidy"},
	})
}
