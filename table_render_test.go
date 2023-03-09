package gutils

import (
	"testing"
)

func TestTerminalRender(t *testing.T) {
	TerminalRender(SetIsTerminal())
}
