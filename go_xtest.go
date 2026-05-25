package gob

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

type GoTest struct {
	cfg TestConfig
}

type TestConfig struct {
	Selector  string
	Verbose   bool
	Race      bool
	Cover     bool
	CoverMode string
	CoverPkg  string
	CoverOut  string
	Timeout   time.Duration
	Count     int
	Tags      []string
	Run       string
	FailFast  bool
	Shuffle   string
	JSON      bool
	Parallel  int
	Env       map[string]string
}

type TestOption interface {
	apply(*TestConfig)
}

type testOptionFunc func(*TestConfig)

func (f testOptionFunc) apply(cfg *TestConfig) {
	f(cfg)
}

func defaultTestConfig() TestConfig {
	return TestConfig{
		Selector:  "./...",
		Timeout:   1 * time.Minute,
		Count:     1,
		CoverMode: "atomic",
		// CoverOut:  "coverage.out",
		Env: map[string]string{},
	}
}

func Test(opts ...TestOption) *GoTest {
	cfg := defaultTestConfig()

	for _, opt := range opts {
		opt.apply(&cfg)
	}

	return &GoTest{cfg: cfg}
}

func (g *GoTest) Run(ctx context.Context) error {
	args := []string{"test"}

	if g.cfg.Verbose {
		args = append(args, "-v")
	}

	if g.cfg.Race {
		args = append(args, "-race")
	}

	if g.cfg.FailFast {
		args = append(args, "-failfast")
	}

	if g.cfg.JSON {
		args = append(args, "-json")
	}

	if g.cfg.Cover {
		args = append(
			args,
			"-cover",
			"-covermode="+g.cfg.CoverMode,
		)

		if g.cfg.CoverOut != "" {
			args = append(args, "-coverprofile="+g.cfg.CoverOut)
		}

		if g.cfg.CoverPkg != "" {
			args = append(args, "-coverpkg="+g.cfg.CoverPkg)
		}
	}

	if g.cfg.Timeout > 0 {
		args = append(args, "-timeout", g.cfg.Timeout.String())
	}

	if g.cfg.Count > 0 {
		args = append(args, "-count", fmt.Sprintf("%d", g.cfg.Count))
	}

	if g.cfg.Run != "" {
		args = append(args, "-run", g.cfg.Run)
	}

	if len(g.cfg.Tags) > 0 {
		args = append(args, "-tags", strings.Join(g.cfg.Tags, ","))
	}

	if g.cfg.Shuffle != "" {
		args = append(args, "-shuffle", g.cfg.Shuffle)
	}

	if g.cfg.Parallel > 0 {
		args = append(args, "-parallel", fmt.Sprintf("%d", g.cfg.Parallel))
	}

	args = append(args, g.cfg.Selector)

	cmd := exec.CommandContext(ctx, "go", args...)

	env := os.Environ()

	for k, v := range g.cfg.Env {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}

	cmd.Env = env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

func (g *GoTest) HTMLCoverage(ctx context.Context, output string) error {
	if !g.cfg.Cover {
		return fmt.Errorf("coverage must be enabled")
	}

	if output == "" {
		output = "coverage.html"
	}

	cmd := exec.CommandContext(
		ctx,
		"go",
		"tool",
		"cover",
		"-html="+g.cfg.CoverOut,
		"-o",
		output,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (g *GoTest) String() string {
	args := []string{"go test"}

	if g.cfg.Verbose {
		args = append(args, "-v")
	}

	if g.cfg.Race {
		args = append(args, "-race")
	}

	if g.cfg.Cover {
		args = append(args, "-cover")
	}

	if g.cfg.Run != "" {
		args = append(args, fmt.Sprintf("-run %q", g.cfg.Run))
	}

	args = append(args, g.cfg.Selector)

	return strings.Join(args, " ")
}

func WithTestSelector(selector string) TestOption {
	return testOptionFunc(func(cfg *TestConfig) {
		cfg.Selector = selector
	})
}

func WithVerbose() TestOption {
	return testOptionFunc(func(cfg *TestConfig) {
		cfg.Verbose = true
	})
}

func WithRace() TestOption {
	return testOptionFunc(func(cfg *TestConfig) {
		cfg.Race = true
	})
}

func WithCoverage() TestOption {
	return testOptionFunc(func(cfg *TestConfig) {
		cfg.Cover = true
	})
}

func WithCoverageOutput(path string) TestOption {
	return testOptionFunc(func(cfg *TestConfig) {
		cfg.Cover = true
		cfg.CoverOut = path
	})
}

func WithCoveragePackage(pkg string) TestOption {
	return testOptionFunc(func(cfg *TestConfig) {
		cfg.CoverPkg = pkg
	})
}

func WithTimeout(timeout time.Duration) TestOption {
	return testOptionFunc(func(cfg *TestConfig) {
		cfg.Timeout = timeout
	})
}

func WithCount(count int) TestOption {
	return testOptionFunc(func(cfg *TestConfig) {
		cfg.Count = count
	})
}

func WithRun(regex string) TestOption {
	return testOptionFunc(func(cfg *TestConfig) {
		cfg.Run = regex
	})
}

func WithTags(tags ...string) TestOption {
	return testOptionFunc(func(cfg *TestConfig) {
		cfg.Tags = tags
	})
}

func WithFailFast() TestOption {
	return testOptionFunc(func(cfg *TestConfig) {
		cfg.FailFast = true
	})
}

func WithShuffle(mode string) TestOption {
	return testOptionFunc(func(cfg *TestConfig) {
		cfg.Shuffle = mode
	})
}

func WithJSON() TestOption {
	return testOptionFunc(func(cfg *TestConfig) {
		cfg.JSON = true
	})
}

func WithParallel(n int) TestOption {
	return testOptionFunc(func(cfg *TestConfig) {
		cfg.Parallel = n
	})
}

func WithTestEnv(key, value string) TestOption {
	return testOptionFunc(func(cfg *TestConfig) {
		cfg.Env[key] = value
	})
}
