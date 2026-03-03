package gob

import (
	"context"
	"fmt"
	"os"
	"os/exec"
)

type GoBuild struct {
	cfg    Config
	matrix []platform
}

func Static(opts ...Option) *GoBuild {
	return newGoBuild(false, opts...)
}

func Dynamic(opts ...Option) *GoBuild {
	return newGoBuild(true, opts...)
}

func newGoBuild(dynamic bool, opts ...Option) *GoBuild {
	cfg := defaultConfig()
	cfg.Dynamic = dynamic

	for _, opt := range opts {
		opt.apply(&cfg)
	}

	return &GoBuild{cfg: cfg}
}

func (g *GoBuild) For(goos, goarch string) *GoBuild {
	g.matrix = append(g.matrix, platform{
		GOOS:   goos,
		GOARCH: goarch,
	})
	return g
}

func (g *GoBuild) Matrix(oses []string, archs []string) *GoBuild {
	for _, osName := range oses {
		for _, archName := range archs {
			g.matrix = append(g.matrix, platform{
				GOOS:   osName,
				GOARCH: archName,
			})
		}
	}
	return g
}

func (g *GoBuild) Run(ctx context.Context) error {
	if err := os.MkdirAll(g.cfg.OutDir, os.ModePerm); err != nil {
		return err
	}

	version := GitVersion()
	ldflags := fmt.Sprintf(
		"-w -s -X %s=%s",
		g.cfg.VersionVar,
		version,
	)

	// If no matrix defined => single build
	if len(g.matrix) == 0 {
		return g.buildOnce(ctx, "", "", ldflags)
	}

	// Matrix builds
	for _, p := range g.matrix {
		if err := g.buildOnce(ctx, p.GOOS, p.GOARCH, ldflags); err != nil {
			return err
		}
	}

	return nil
}

func (g *GoBuild) buildOnce(
	ctx context.Context,
	goos string,
	goarch string,
	ldflags string,
) error {
	output := g.cfg.OutDir

	if goos != "" && goarch != "" {
		output = fmt.Sprintf(
			"%s/%s-%s/",
			g.cfg.OutDir,
			goos,
			goarch,
		)

		/*if goos == "windows" {
			output += ".exe"
		}*/
	}

	cmd := exec.CommandContext(
		ctx,
		"go", "build",
		"-ldflags", ldflags,
		"-o", output,
		g.cfg.Selector,
	)

	env := os.Environ()

	// Static vs Dynamic
	if g.cfg.Dynamic {
		env = append(env, "CGO_ENABLED=1")
	} else {
		env = append(env, "CGO_ENABLED=0")
	}

	// Matrix overrides
	if goos != "" {
		env = append(env, "GOOS="+goos)
	}
	if goarch != "" {
		env = append(env, "GOARCH="+goarch)
	}

	// Custom env
	for k, v := range g.cfg.Env {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}

	cmd.Env = env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

func (g *GoBuild) String() string {
	mode := "static"
	if g.cfg.Dynamic {
		mode = "dynamic"
	}

	selector := g.cfg.Selector
	if selector == "" {
		selector = "./..."
	}

	baseOut := g.cfg.OutDir
	if baseOut == "" {
		baseOut = "bin"
	}

	// Matrix build
	if len(g.matrix) > 0 {
		paths := make([]string, len(g.matrix))
		for i, p := range g.matrix {
			path := fmt.Sprintf("%s/%s-%s", baseOut, p.GOOS, p.GOARCH)
			if p.GOOS == "windows" {
				path += ".exe"
			}
			paths[i] = path
		}
		return fmt.Sprintf(
			"Build %s binaries from %s => %v",
			mode,
			selector,
			paths,
		)
	}

	// Single build
	return fmt.Sprintf(
		"Build %s binary from %s => %s",
		mode,
		selector,
		baseOut,
	)
}
