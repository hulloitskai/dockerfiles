package main

import (
	"os"
	"testing"
)

func TestConfig_resolveWithShell(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %e", err)
	}
	cfg := &config{Context: dir}

	in := "$(hostname)"
	out, err := cfg.resolveWithShell(in)
	if err != nil {
		t.Error(err)
	}
	if out == in {
		t.Error("Expected result not to equal input.")
	} else {
		t.Logf("Hostname: %q", out)
	}

	in = "$(cat asdfasdfasdf234234q)"
	out, err = cfg.resolveWithShell(in)
	if err == nil {
		t.Errorf("Expected invalid command to error out, instead got: %q", out)
	} else {
		t.Logf("Error (expected): %v", err)
	}
}
