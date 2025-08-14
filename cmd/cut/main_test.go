package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRun_OK(t *testing.T) {
	in := strings.NewReader("a\tb\tc\n1\t2\t3\n")
	var out, errBuf bytes.Buffer
	code := run([]string{"-f", "1,3"}, in, &out, &errBuf)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d, stderr=%q", code, errBuf.String())
	}
	if got, want := out.String(), "a\tc\n1\t3\n"; got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestRun_MissingF(t *testing.T) {
	in := strings.NewReader("")
	var out, errBuf bytes.Buffer
	code := run([]string{}, in, &out, &errBuf)
	if code != 2 {
		t.Fatalf("expected exit 2, got %d", code)
	}
	if errBuf.Len() == 0 {
		t.Fatalf("expected error message in stderr")
	}
}

func TestRun_InvalidF(t *testing.T) {
	in := strings.NewReader("")
	var out, errBuf bytes.Buffer
	code := run([]string{"-f", "0-"}, in, &out, &errBuf)
	if code != 2 {
		t.Fatalf("expected exit 2, got %d", code)
	}
}

func TestRun_SeparatedOnly(t *testing.T) {
	in := strings.NewReader("a,b,c\nnode\n1,2\n")
	var out, errBuf bytes.Buffer
	code := run([]string{"-f", "1", "-d", ",", "-s"}, in, &out, &errBuf)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if got, want := out.String(), "a\n1\n"; got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}
