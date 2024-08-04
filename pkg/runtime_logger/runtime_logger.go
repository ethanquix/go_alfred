package runtime_logger

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/ethanquix/go_alfred/pkg/globals"
)

type gcpLogEntry struct {
	Severity string      `json:"severity"`
	Message  interface{} `json:"message"`
	//HttpRequest    *httpRequest      `json:"httpRequest,omitempty"`
	Timestamp time.Time         `json:"timestamp"`
	Labels    map[string]string `json:"logging.googleapis.com/labels,omitempty"`
	//Operation      *operation        `json:"logging.googleapis.com/operation,omitempty"`
	//SourceLocation *sourceLocation   `json:"logging.googleapis.com/sourceLocation,omitempty"`
	SpanID       string `json:"logging.googleapis.com/spanId,omitempty"`
	TraceID      string `json:"logging.googleapis.com/trace,omitempty"`
	TraceSampled bool   `json:"logging.googleapis.com/trace_sampled,omitempty"`
}

func envToLogLevel(envVar string) log.Level {
	switch strings.ToLower(envVar) {
	case "debug":
		return log.DebugLevel
	case "info":
		return log.InfoLevel
	case "warn":
		return log.WarnLevel
	case "error":
		return log.ErrorLevel
	case "":
		fmt.Printf("WARNING: Log level set to Debug by default. Set 'RUNTIME_LOG_LEVEL' env var to one of 'debug/info/warn/error' to modify this behavior!\n")
		return log.DebugLevel // default is debug
	}
	panic("Invalid value for 'RUNTIME_LOG_LEVEL'. Must be one of 'debug/info/warn/error' (or empty)")
}

func InitRuntimeLogger(forceProd ...bool) {
	// get log level from ENV (RUNTIME_LOG_LEVEL)
	runtimeLogLevel := os.Getenv("RUNTIME_LOG_LEVEL")
	log.SetLevel(envToLogLevel(runtimeLogLevel))

	if globals.IS_PROD() || len(forceProd) > 0 {
		// if prod set, we use GCP optimized logging (JSON)
		log.SetTimeFormat(time.RFC3339)
		log.MessageKey = "message"
		log.LevelKey = "severity"
		log.TimestampKey = "time"
		log.SetFormatter(log.JSONFormatter)
	}
}

func GetTraceID(header http.Header) string {
	projectID := "nofo-prod"

	// Derive the traceID associated with the current request.
	var trace string
	traceHeader := header.Get("X-Cloud-Trace-Context")
	traceParts := strings.Split(traceHeader, "/")
	if len(traceParts) > 0 && len(traceParts[0]) > 0 {
		trace = fmt.Sprintf("projects/%s/traces/%s", projectID, traceParts[0])
	}
	return trace
}

func GetLoggerForRequest(r *http.Request) *log.Logger {
	trace := GetTraceID(r.Header)
	return log.With("trace", trace)
}
