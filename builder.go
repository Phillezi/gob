package gob

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"time"
)

var (
	ErrTargetNotFound = errors.New("target not found")
	ErrNoTargets      = errors.New("no targets defined")
)

type Builder struct {
	targets       map[string]Target
	defaultTarget string
	logger        *slog.Logger
	args          []string
	exitOnError   *bool
}

func New(opts ...BuilderOption) *Builder {
	cfg := defaultBuilderConfig()

	for _, opt := range opts {
		opt(&cfg)
	}

	if cfg.Logger == nil {
		cfg.Logger = slog.New(slog.DiscardHandler)
	}

	b := &Builder{
		targets:       make(map[string]Target),
		logger:        cfg.Logger,
		defaultTarget: cfg.DefaultTarget,
		args:          cfg.Args,
		exitOnError:   cfg.ExitOnError,
	}

	return b
}

func (b *Builder) Default(name string) *Builder {
	b.defaultTarget = name
	return b
}

func (b *Builder) Add(name string, t Target) *Builder {
	b.targets[name] = t
	return b
}

func (b *Builder) Run() error {
	targetLen := len(b.targets)
	if targetLen == 0 {
		b.logger.Error("No targets defined in builder")
		if b.exitOnError != nil && *b.exitOnError {
			os.Exit(1)
		}
		return ErrNoTargets
	}
	if targetLen == 1 {
		if _, found := b.targets[b.defaultTarget]; !found {
			for k := range b.targets {
				b.defaultTarget = k
			}
		}
	}
	args := b.args
	if args == nil {
		args = os.Args[1:]
	}
	if len(args) == 0 {
		args = append(args, b.defaultTarget)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	exitCode := 0
	var errs error

targets:
	for _, arg := range args {
		t, ok := b.targets[arg]
		if !ok {
			b.logger.Error("Unknown", "target", arg)
			errs = errors.Join(errs, ErrTargetNotFound)
			exitCode = 1
			break targets
		}

		b.logger.Info("Running", "target", arg)
		start := time.Now()
		if err := t.Run(ctx); err != nil {
			b.logger.Error("Failed", "target", arg, "error", err)
			errs = errors.Join(errs, err)
			exitCode = 1
		}
		end := time.Now()
		b.logger.Info("Built", "target", arg, "duration", end.Sub(start))
	}

	if errors.Is(errs, ErrTargetNotFound) {
		PrintTargets(os.Stderr, b.targets)
	}

	if b.exitOnError != nil && *b.exitOnError {
		os.Exit(exitCode)
	}

	return errs
}
