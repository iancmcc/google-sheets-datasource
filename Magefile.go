//+build mage

package main

import (
	"fmt"
	"os"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

const dsName string = "sheets-datasource"

func buildBackend(variant string, enableDebug bool, env map[string]string) error {
	varStr := ""
	if variant != "" {
		varStr = fmt.Sprintf("_%s", variant)
	}
	args := []string{
		"build", "-o", fmt.Sprintf("dist/%s%s", dsName, varStr), "-tags", "netgo",
	}
	if enableDebug {
		args = append(args, "-gcflags=all=-N -l")
	} else {
		args = append(args, []string{"-ldflags", "-w"}...)
	}
	args = append(args, "./pkg")
	// TODO: Change to sh.RunWithV once available.
	if err := sh.RunWith(env, "go", args...); err != nil {
		return err
	}

	return nil
}

// Build is a namespace.
type Build mg.Namespace

// BackendLinux builds the back-end plugin for Linux.
func (Build) BackendLinux() error {
	env := map[string]string{
		"GOARCH": "amd64",
		"GOOS":   "linux",
	}
	return buildBackend("linux_amd64", false, env)
}

// BackendLinuxDebug builds the back-end plugin for Linux in debug mode.
func (Build) BackendLinuxDebug() error {
	env := map[string]string{
		"GOARCH": "amd64",
		"GOOS":   "linux",
	}
	return buildBackend("linux_amd64", true, env)
}

// Frontend builds the front-end for production.
func (Build) Frontend() error {
	mg.Deps(Deps)
	return sh.RunV("./node_modules/.bin/grafana-toolkit", "plugin:build")
}

// FrontendDev builds the front-end for development.
func (Build) FrontendDev() error {
	mg.Deps(Deps)
	return sh.RunV("./node_modules/.bin/grafana-toolkit", "plugin:dev")
}

// BuildAll builds both back-end and front-end components.
func BuildAll() {
	b := Build{}
	mg.Deps(b.BackendLinux, b.Frontend)
}

// Deps installs dependencies.
func Deps() error {
	return sh.RunV("yarn", "install")
}

// Test runs all tests.
func Test() error {
	mg.Deps(Deps)

	if err := sh.RunV("go", "test", "./pkg/..."); err != nil {
		return nil
	}
	return sh.RunV("yarn", "test")
}

// Lint lints the sources.
func Lint() error {
	return sh.RunV("golangci-lint", "run", "./...")
}

// Format formats the sources.
func Format() error {
	if err := sh.RunV("gofmt", "-w", "."); err != nil {
		return err
	}

	return nil
}

// Dev builds the plugin in dev mode.
func Dev() error {
	b := Build{}
	mg.Deps(b.BackendLinuxDebug, b.FrontendDev) // TODO: only the current architecture
	return nil
}

// Watch will build the plugin in dev mode and then update when the frontend files change.
func Watch() error {
	b := Build{}
	mg.Deps(b.BackendLinuxDebug)

	// The --watch will never return
	return sh.RunV("./node_modules/.bin/grafana-toolkit", "plugin:dev", "--watch")
}

// Clean cleans build artifacts, by deleting the dist directory.
func Clean() error {
	return os.RemoveAll("dist")
}

// Default configures the default target.
var Default = BuildAll
