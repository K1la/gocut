package cut

import "testing"

func TestParseFields_SingleAndRanges(t *testing.T) {
	fs, err := ParseFields("1,3-5,7-, -2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// 1
	if !fs.Includes(1) {
		t.Errorf("expected to include 1")
	}
	// 3-5
	for i := 3; i <= 5; i++ {
		if !fs.Includes(i) {
			t.Errorf("expected to include %d", i)
		}
	}
	// 7-
	for i := 7; i <= 100; i += 13 {
		if !fs.Includes(i) {
			t.Errorf("expected to include %d (7-)", i)
		}
	}
	// -2
	for i := 1; i <= 2; i++ {
		if !fs.Includes(i) {
			t.Errorf("expected to include %d (-2)", i)
		}
	}
	if fs.Includes(6) {
		t.Errorf("did not expect to include 6")
	}
}

func TestParseFields_Errors(t *testing.T) {
	cases := []string{"", ",", "-", "0", "-0", "a-b", "2-1"}
	for _, c := range cases {
		if _, err := ParseFields(c); err == nil {
			t.Errorf("expected error for %q", c)
		}
	}
}
