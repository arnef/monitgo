package parser_test

import (
	"testing"

	"github.com/arnef/monitgo/interal/watcher/parser"
)

func TestParserFree(t *testing.T) {
	in := `               total        used        free      shared  buff/cache   available
	Mem:     16451350528  3031846912 10451210240   559390720  2968293376 12510736384
	Swap:              0           0           0
	`

	total, used, err := parser.Free(in)
	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}
	if total != 16451350528 {
		t.Errorf("expected total bytes to be 16451350528 but got %d", total)
	}
	if used != 3031846912 {
		t.Errorf("expected used bytes to be 3031846912 but got %d", used)
	}
}
