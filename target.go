package gob

import (
	"context"
	"fmt"
	"io"
	"sort"
)

type Target interface {
	Run(ctx context.Context) error
	String() string
}

type platform struct {
	GOOS   string
	GOARCH string
}

var (
	PopularOSes   = []string{GOOSLinux, GOOSDarwin, GOOSWindows}
	PopularArches = []string{GOARCHAMD64, GOARCHARM64}
)

// Common GOOS constants
const (
	GOOSLinux   = "linux"
	GOOSDarwin  = "darwin"
	GOOSWindows = "windows"
	GOOSFreeBSD = "freebsd"
	GOOSNetBSD  = "netbsd"
	GOOSOpenBSD = "openbsd"
	GOOSSolaris = "solaris"
	GOOSAIX     = "aix"
)

// Common GOARCH constants
const (
	GOARCHAMD64    = "amd64"
	GOARCHARM64    = "arm64"
	GOARCH386      = "386"
	GOARCHARM      = "arm"
	GOARCHPPC64    = "ppc64"
	GOARCHPPC64LE  = "ppc64le"
	GOARCHMIPS     = "mips"
	GOARCHMIPSLE   = "mipsle"
	GOARCHMIPS64   = "mips64"
	GOARCHMIPS64LE = "mips64le"
	GOARCHS390X    = "s390x"
	GOARCHRISCV64  = "riscv64"
)

// PrintTargets prints available build targets in a clean format.
func PrintTargets(w io.Writer, targets map[string]Target) {
	if len(targets) == 0 {
		fmt.Fprintln(w, "No targets available.")
		return
	}

	names := make([]string, 0, len(targets))
	for name := range targets {
		names = append(names, name)
	}
	sort.Strings(names)

	maxLen := 0
	for _, name := range names {
		if len(name) > maxLen {
			maxLen = len(name)
		}
	}

	fmt.Fprintln(w, "Available targets:")

	for _, name := range names {
		desc := targets[name].String()
		if desc == "" {
			desc = "(no description)"
		}

		fmt.Fprintf(
			w,
			"  %-*s  %s\n",
			maxLen,
			name,
			desc,
		)
	}
}
