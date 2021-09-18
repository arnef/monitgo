package parser_test

import (
	"testing"

	"github.com/arnef/monitgo/interal/watcher/parser"
)

func TestParserLoadAverage(t *testing.T) {
	// linux
	val := " 19:23:11 up 1 day,  2:31,  1 user,  load average: 2,35, 2,36, 1,86"

	load, err := parser.LoadAverage(val)
	if err != nil {
		t.Errorf("expected err to be nil but got %v", err)
	}
	if len(load) != 3 {
		t.Errorf("expected 3 values for load but got %d", len(load))
	}
	if load[0] != 235 || load[1] != 236 || load[2] != 186 {
		t.Errorf("Expected load values [235 236 186] bot got %v", load)
	}
	// diskstation
	val = " 19:11:16 up 2 days,  6:34,  1 user,  load average: 0.25, 0.17, 0.13 [IO: 0.24, 0.13, 0.10 CPU: 0.00, 0.02, 0.00]"
	load, err = parser.LoadAverage(val)
	if err != nil {
		t.Errorf("expected err to be nil but got %v", err)
	}
	if len(load) != 3 {
		t.Errorf("expected 3 values for load but got %d", len(load))
	}
	if load[0] != 25 || load[1] != 17 || load[2] != 13 {
		t.Errorf("Expected load values [25 17 13] bot got %v", load)
	}
}
