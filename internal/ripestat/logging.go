package ripestat

import (
	"io"
	"log"
	"os"
	"strings"
	"sync/atomic"
)

const debugEnvVar = "RIPESTAT_DEBUG_HTTP"

var (
	debugLogger      = log.New(io.Discard, "[ripestat] ", log.LstdFlags)
	httpDebugEnabled atomic.Bool
)

func init() {
	if isTruthy(os.Getenv(debugEnvVar)) {
		debugLogger.SetOutput(os.Stderr)
		httpDebugEnabled.Store(true)
	}
}

func debugf(format string, args ...any) {
	if !httpDebugEnabled.Load() {
		return
	}
	debugLogger.Printf(format, args...)
}

func isTruthy(v string) bool {
	v = strings.TrimSpace(v)
	if v == "" {
		return false
	}
	switch strings.ToLower(v) {
	case "1", "true", "yes", "on":
		return true
	default:
		return false
	}
}

func setDebugOutput(w io.Writer) {
	if w == nil {
		debugLogger.SetOutput(io.Discard)
		return
	}
	debugLogger.SetOutput(w)
}

func setHTTPDebug(enabled bool) {
	httpDebugEnabled.Store(enabled)
}
