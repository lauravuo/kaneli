package main

import "testing"

func TestTrimEscValue(t *testing.T) {
	got := trimEscValue("CDATA[  ]")
	if got != "" {
		t.Errorf("Trimming value failed: %s", got)
	}
}
