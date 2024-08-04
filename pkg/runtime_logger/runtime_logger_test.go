package runtime_logger

import (
	"testing"

	"github.com/charmbracelet/log"
	"github.com/ethanquix/go_alfred/pkg/globals"
)

func TestInitRuntimeLogger(t *testing.T) {
	globals.SetIsProd(true)
	InitRuntimeLogger()

	log.Infof("hello world")
}
