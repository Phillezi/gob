package gob

import (
	"log/slog"
	"maps"
)

type Option interface {
	apply(*Config)
}

type optionFunc func(*Config)

func (f optionFunc) apply(c *Config) {
	f(c)
}

func WithOutDir(dir string) Option {
	return optionFunc(func(c *Config) {
		c.OutDir = dir
	})
}

func WithSelector(sel string) Option {
	return optionFunc(func(c *Config) {
		c.Selector = sel
	})
}

func WithVersionVar(v string) Option {
	return optionFunc(func(c *Config) {
		c.VersionVar = v
	})
}

func WithDynamic(v bool) Option {
	return optionFunc(func(c *Config) {
		c.Dynamic = v
	})
}

func WithEnv(key, value string) Option {
	return optionFunc(func(c *Config) {
		if c.Env == nil {
			c.Env = make(map[string]string)
		}
		c.Env[key] = value
	})
}

func WithConfig(cfg Config) Option {
	return optionFunc(func(c *Config) {
		mergeConfig(c, cfg)
	})
}

func mergeConfig(dst *Config, src Config) {
	if src.OutDir != "" {
		dst.OutDir = src.OutDir
	}
	if src.Selector != "" {
		dst.Selector = src.Selector
	}
	if src.VersionVar != "" {
		dst.VersionVar = src.VersionVar
	}
	if src.Dynamic {
		dst.Dynamic = true
	}

	if src.Env != nil {
		if dst.Env == nil {
			dst.Env = make(map[string]string)
		}
		maps.Copy(dst.Env, src.Env)
	}
}

type BuilderOption func(cfg *BuilderConfig)

func WithDefaultTarget(name string) BuilderOption {
	return func(c *BuilderConfig) {
		c.DefaultTarget = name
	}
}

func WithLogger(l *slog.Logger) BuilderOption {
	return func(c *BuilderConfig) {
		c.Logger = l
	}
}

func WithExitOnError(v bool) BuilderOption {
	return func(c *BuilderConfig) {
		c.ExitOnError = &v
	}
}

func WithArgs(args []string) BuilderOption {
	return func(c *BuilderConfig) {
		c.Args = args
	}
}

func WithBuilderConfig(cfg BuilderConfig) BuilderOption {
	return func(dst *BuilderConfig) {
		mergeBuilderConfig(dst, cfg)
	}
}

func mergeBuilderConfig(dst *BuilderConfig, src BuilderConfig) {
	if src.DefaultTarget != "" {
		dst.DefaultTarget = src.DefaultTarget
	}

	if src.Logger != nil {
		dst.Logger = src.Logger
	}

	if src.ExitOnError != nil {
		dst.ExitOnError = dst.ExitOnError
	}

	if src.Args != nil {
		dst.Args = src.Args
	}
}
