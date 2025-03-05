//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	// mage:import
	build "github.com/grafana/grafana-plugin-sdk-go/build"
)

// Default configures the default target.
var Default = build.BuildAll

// BuildPlugin builds the plugin
func BuildPlugin() error {
	// Build for current platform
	return buildForPlatform(runtime.GOOS, runtime.GOARCH)
}

// BuildAllPlatforms builds the plugin for all supported platforms
func BuildAllPlatforms() error {
	platforms := []struct {
		os   string
		arch string
	}{
		{"linux", "amd64"},
		{"linux", "arm64"},
		{"darwin", "amd64"},
		{"darwin", "arm64"},
		{"windows", "amd64"},
	}

	for _, p := range platforms {
		if err := buildForPlatform(p.os, p.arch); err != nil {
			return err
		}
	}

	return nil
}

// buildForPlatform builds the plugin for a specific platform
func buildForPlatform(goos, goarch string) error {
	output := fmt.Sprintf("dist/autotask-datasource_%s_%s", goos, goarch)
	if goos == "windows" {
		output += ".exe"
	}

	cmd := exec.Command("go", "build",
		"-o", output,
		"-ldflags", "-w -s",
		"./pkg/main.go",
	)
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("GOOS=%s", goos),
		fmt.Sprintf("GOARCH=%s", goarch),
	)

	return cmd.Run()
}

// CleanBuild cleans build artifacts
func CleanBuild() error {
	return os.RemoveAll("dist")
}

// DevPlugin builds and runs the plugin in development mode
func DevPlugin() error {
	if err := BuildPlugin(); err != nil {
		return err
	}

	cmd := exec.Command("go", "run", "./pkg/main.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
