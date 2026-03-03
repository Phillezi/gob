# gob

Build go projects with recipes written in go.

## Example

Create a go file in your project, you can name it what you want `build.go` for example

```go
//go:build ignore

package main

import (
	"log/slog"

	"github.com/phillezi/gob"
)

func init() {
	logger := slog.New(gob.NewPrettyHandler(nil))
	slog.SetDefault(logger)
}

func main() {
  // takes in options
	gob.New(gob.WithDefaultTarget("all")).Add(
		"all",
    // takes in options to customize
    // can also be chaned with .For(os, arch)
    // or .Matrix([]string{"linux"}, []string{"amd64", "arm64"})
		gob.Static(),
	).Add(
		"clean",
		gob.Clean(),
	).Run()
}
```

To build your project you simply run:

```bash
go run build.go # the path to the file you created
```

Your project will be built statically for your platform and all the binaries will end up in `./bin/` by default. There are more options, like building for multiple architectures etc.

See [examples](./examples/) for more examples.

## Why

This is just a small tool that automates setting common build flags that I usually set for my go projects, such as `CGO_ENABLED=0` `-ldflags="-w -s"` and also version tagging using git, making the version string available in the binary `-ldflags="-X main.version=vX.Y.Z"`.

So instead of having to write:

```bash
GIT_TAG=$(git describe --tags --abbrev=0 2>/dev/null)
GIT_COMMIT=$(git rev-parse HEAD 2>/dev/null)
GIT_TAG_COMMIT=$(git rev-list -n 1 "$GIT_TAG" 2>/dev/null)

if [ -n "$(git status --porcelain 2>/dev/null)" ]; then
    GIT_DIRTY=true
else
    GIT_DIRTY=false
fi

GIT_SHORT=${GIT_COMMIT:0:7}

if [ "$GIT_DIRTY" = false ]; then
    if [ "$GIT_TAG_COMMIT" = "$GIT_COMMIT" ]; then
        GIT_VERSION="$GIT_TAG"
    else
        GIT_VERSION="${GIT_TAG}-${GIT_SHORT}"
    fi
else
    GIT_VERSION="${GIT_TAG}-dirty-${GIT_SHORT}"
fi

# and then finally building
CGO_ENABLED=0 go build -ldflags="-w -s -X main.version=$GIT_VERSION" -o ./bin/ ./...
```

This can be used to do it in a cross-platform way.

You can of-course use a `Makefile` for this, but it adds another dependency. And this is just go code, which also makes it possible to re-use the build recipe when building Docker images with slim builder images that doesnt have Make installed.

