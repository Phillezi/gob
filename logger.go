package gob

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"sync"
)

const (
	reset  = "\033[0m"
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	cyan   = "\033[36m"
)

type prettyCore struct {
	out   io.Writer
	mu    sync.Mutex
	level slog.Level
}

type PrettyHandler struct {
	core  *prettyCore
	attrs []slog.Attr
	group string
}

func NewPrettyHandler(opts *slog.HandlerOptions) *PrettyHandler {
	level := slog.LevelInfo
	if opts != nil && opts.Level != nil {
		level = opts.Level.Level()
	}

	core := &prettyCore{
		out:   os.Stderr,
		level: level,
	}

	return &PrettyHandler{
		core: core,
	}
}

func (h *PrettyHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.core.level
}

func (h *PrettyHandler) Handle(_ context.Context, r slog.Record) error {
	h.core.mu.Lock()
	defer h.core.mu.Unlock()

	prefix, color := levelStyle(r.Level)

	var b strings.Builder

	b.WriteString(color)
	b.WriteString(prefix)
	b.WriteString(reset)
	b.WriteString(" ")

	b.WriteString(r.Message)

	for _, a := range h.attrs {
		writeAttr(&b, h.group, a)
	}

	r.Attrs(func(a slog.Attr) bool {
		writeAttr(&b, h.group, a)
		return true
	})

	b.WriteString("\n")

	_, err := h.core.out.Write([]byte(b.String()))
	return err
}

func (h *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &PrettyHandler{
		core:  h.core,
		attrs: append(append([]slog.Attr{}, h.attrs...), attrs...),
		group: h.group,
	}
}

func (h *PrettyHandler) WithGroup(name string) slog.Handler {
	group := name
	if h.group != "" {
		group = h.group + "." + name
	}

	return &PrettyHandler{
		core:  h.core,
		attrs: h.attrs,
		group: group,
	}
}

func writeAttr(b *strings.Builder, group string, a slog.Attr) {
	key := a.Key
	if group != "" {
		key = group + "." + key
	}

	b.WriteString(" ")
	b.WriteString(key)
	b.WriteString("=")
	fmt.Fprint(b, a.Value.Any())
}

func levelStyle(level slog.Level) (prefix, color string) {
	switch {
	case level >= slog.LevelError:
		return "[*]", red
	case level >= slog.LevelWarn:
		return "[!]", yellow
	case level >= slog.LevelInfo:
		return "[+]", green
	default:
		return "[#]", cyan
	}
}
