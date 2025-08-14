package cut

import (
	"fmt"
	"strconv"
	"strings"
)

// FieldSelection represents a parsed -f flag that can check
// whether a 1-based field index should be included.
type FieldSelection struct {
	exact  map[int]struct{}
	ranges []fieldRange
}

type fieldRange struct {
	start int // 0 means open start
	end   int // 0 means open end
}

// ParseFields parses specifications like "1,3-5,7-, -4".
// Indices are 1-based. Zero and negative indices are invalid.
func ParseFields(spec string) (FieldSelection, error) {
	fs := FieldSelection{exact: make(map[int]struct{})}

	if strings.TrimSpace(spec) == "" {
		return fs, fmt.Errorf("empty fields spec")
	}

	parts := strings.Split(spec, ",")
	for _, p := range parts {
		token := strings.TrimSpace(p)
		if token == "" {
			return fs, fmt.Errorf("empty token in fields spec")
		}

		if strings.Contains(token, "-") {
			// range: start-end where start or end may be empty
			se := strings.SplitN(token, "-", 2)
			if len(se) != 2 {
				return fs, fmt.Errorf("invalid range: %q", token)
			}
			var r fieldRange

			// start
			if strings.TrimSpace(se[0]) == "" {
				r.start = 0
			} else {
				v, err := parsePositiveInt(se[0])
				if err != nil {
					return fs, fmt.Errorf("invalid range start in %q: %w", token, err)
				}
				r.start = v
			}

			// end
			if strings.TrimSpace(se[1]) == "" {
				r.end = 0
			} else {
				v, err := parsePositiveInt(se[1])
				if err != nil {
					return fs, fmt.Errorf("invalid range end in %q: %w", token, err)
				}
				if r.start != 0 && v < r.start {
					return fs, fmt.Errorf("range end < start in %q", token)
				}
				r.end = v
			}

			if r.start == 0 && r.end == 0 {
				return fs, fmt.Errorf("invalid open range '-' in %q", token)
			}
			fs.ranges = append(fs.ranges, r)
			continue
		}

		// single index
		v, err := parsePositiveInt(token)
		if err != nil {
			return fs, fmt.Errorf("invalid field %q: %w", token, err)
		}
		fs.exact[v] = struct{}{}
	}

	return fs, nil
}

func parsePositiveInt(s string) (int, error) {
	v, err := strconv.Atoi(strings.TrimSpace(s))
	if err != nil || v <= 0 {
		if err == nil {
			return 0, fmt.Errorf("must be positive integer")
		}
		return 0, err
	}
	return v, nil
}

// Includes returns true if the given 1-based index is selected.
func (fs FieldSelection) Includes(index int) bool {
	if index <= 0 {
		return false
	}
	if _, ok := fs.exact[index]; ok {
		return true
	}
	for _, r := range fs.ranges {
		if r.start == 0 && r.end != 0 {
			if index <= r.end {
				return true
			}
			continue
		}
		if r.end == 0 && r.start != 0 {
			if index >= r.start {
				return true
			}
			continue
		}
		if r.start != 0 && r.end != 0 {
			if index >= r.start && index <= r.end {
				return true
			}
		}
	}
	return false
}
