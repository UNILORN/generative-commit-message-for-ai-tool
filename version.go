package main

import (
	"fmt"
	"runtime/debug"
)

// version is set via ldflags during build (e.g., by GoReleaser)
// Example: go build -ldflags "-X main.version=v1.0.0"
var version = ""

// getVersion returns the version information for this binary.
// It checks multiple sources in the following priority order:
//  1. ldflags-injected version (used by GoReleaser)
//  2. Module version from debug.ReadBuildInfo (used by go install)
//  3. "dev" as fallback for development builds
func getVersion() string {
	// 1. Check if version was set via ldflags (GoReleaser)
	if version != "" {
		return version
	}

	// 2. Try to get version from build info (go install)
	if info, ok := debug.ReadBuildInfo(); ok {
		if info.Main.Version != "" && info.Main.Version != "(devel)" {
			return info.Main.Version
		}
	}

	// 3. Fallback to dev for development builds
	return "dev"
}

// printVersion prints the version information to stdout
func printVersion() {
	fmt.Printf("generate-auto-commit-message version %s\n", getVersion())
}
