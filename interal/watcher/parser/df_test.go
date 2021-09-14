package parser_test

import (
	"testing"

	"github.com/arnef/monitgo/interal/watcher/parser"
)

func TestParserDf(t *testing.T) {
	in := `Filesystem     1K-blocks     Used
	dev              8023116        0
	run              8032884     1924
	/dev/nvme0n1p5 102626052 30811364
	tmpfs            8032884    94124
	tmpfs            8032888    73756
	/dev/nvme0n1p6 388355236 90397596
	/dev/nvme0n1p1    262144    97532
	tmpfs            1606576      100`

	total, used, err := parser.Df(in)
	if err != nil {
		t.Errorf("expected err to be nil but got %v", err)
	}
	if total != 491243432 {
		t.Errorf("expected total to be 491243432 but got %d", total)
	}
	if used != 121306492 {
		t.Errorf("expected used to be 121306492 but got %d", used)
	}
}
