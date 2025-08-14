package cut

import (
	"bytes"
	"strings"
	"testing"
)

func TestProcess_Basic(t *testing.T) {
	fs, _ := ParseFields("1,3")
	p := Processor{Delimiter: []byte("\t"), FieldSelection: fs}
	in := strings.NewReader("a\tb\tc\n1\t2\t3\nno_delim\n")
	var out bytes.Buffer
	if err := p.Process(in, &out); err != nil {
		t.Fatalf("Process error: %v", err)
	}
	got := out.String()
	want := "a\tc\n1\t3\nno_delim\n"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestProcess_SeparatedOnly(t *testing.T) {
	fs, _ := ParseFields("1")
	p := Processor{Delimiter: []byte(","), FieldSelection: fs, SeparatedOnly: true}
	in := strings.NewReader("a,b,c\nnode\n1,2\n")
	var out bytes.Buffer
	if err := p.Process(in, &out); err != nil {
		t.Fatalf("Process error: %v", err)
	}
	got := out.String()
	want := "a\n1\n"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestProcess_Ranges(t *testing.T) {
	fs, _ := ParseFields("2-3,5-")
	p := Processor{Delimiter: []byte("|"), FieldSelection: fs}
	in := strings.NewReader("a|b|c|d|e|f\n1|2|3|4\n")
	var out bytes.Buffer
	if err := p.Process(in, &out); err != nil {
		t.Fatalf("Process error: %v", err)
	}
	got := out.String()
	want := "b|c|e|f\n2|3\n"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestProcess_MultiByteDelimiter(t *testing.T) {
	fs, _ := ParseFields("1,3")
	p := Processor{Delimiter: []byte("<->"), FieldSelection: fs}
	in := strings.NewReader("a<->b<->c\n1<->2<->3\n")
	var out bytes.Buffer
	if err := p.Process(in, &out); err != nil {
		t.Fatalf("Process error: %v", err)
	}
	got := out.String()
	want := "a<->c\n1<->3\n"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}
