package gob

import "log/slog"

type Config struct {
	OutDir     string
	Selector   string
	VersionVar string
	Dynamic    bool
	Env        map[string]string
}

func defaultConfig() Config {
	return Config{
		OutDir:     "bin",
		Selector:   "./...",
		VersionVar: "main.version",
		Dynamic:    false,
		Env:        map[string]string{},
	}
}

type BuilderConfig struct {
	DefaultTarget string
	Logger        *slog.Logger
	ExitOnError   *bool
	Args          []string
}

func defaultBuilderConfig() BuilderConfig {
	return BuilderConfig{
		DefaultTarget: "",
		Logger:        slog.Default(),
		// if your LSP errors here it might not be on the latest go
		// version, the line below works on go >1.26
		ExitOnError: new(true),
		Args:        nil, // nil means use os.Args
	}
}
