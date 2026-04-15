package version

import (
	"runtime"
)

var (
	// Current is set at build time via -ldflags
	Current = "dev"
	// Go is the Go runtime version
	Go = runtime.Version()[2:] // "go1.22.2" -> "1.22.2"
)
