package gob

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

// CMD represents a generic command/task
type CMD struct {
	cfg    Config
	action func(ctx context.Context, cfg Config) error // the function to run
	desc   string                                      // string for help/description
}

// NewCMD creates a generic command with functional options and a runnable action
func NewCMD(action func(ctx context.Context, cfg Config) error, desc string, opts ...Option) *CMD {
	cfg := defaultConfig()
	for _, opt := range opts {
		opt.apply(&cfg)
	}

	return &CMD{
		cfg:    cfg,
		action: action,
		desc:   desc,
	}
}

// Run executes the command
func (c *CMD) Run(ctx context.Context) error {
	if c.action == nil {
		return fmt.Errorf("no action defined for CMD")
	}
	return c.action(ctx, c.cfg)
}

// String returns a description for help/pretty printing
func (c *CMD) String() string {
	if c.desc != "" {
		return c.desc
	}
	return "Generic command"
}

// Clean returns a CMD that removes the output directory
func Clean(opts ...Option) *CMD {
	return NewCMD(func(ctx context.Context, cfg Config) error {
		if cfg.OutDir == "" {
			return fmt.Errorf("no output directory set")
		}

		absPath, err := filepath.Abs(cfg.OutDir)
		if err != nil {
			return fmt.Errorf("failed to resolve path %s: %w", cfg.OutDir, err)
		}

		if err := os.RemoveAll(absPath); err != nil {
			return fmt.Errorf("failed to remove directory %s: %w", absPath, err)
		}

		return nil
	}, fmt.Sprintf("Clean task => %s", defaultConfig().OutDir), opts...)
}
