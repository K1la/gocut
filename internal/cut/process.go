package cut

import (
	"bufio"
	"bytes"
	"io"
)

type Processor struct {
	Delimiter      []byte
	SeparatedOnly  bool
	FieldSelection FieldSelection
}

func (p Processor) Process(r io.Reader, w io.Writer) error {
	if len(p.Delimiter) == 0 {
		p.Delimiter = []byte{'\t'}
	}
	scanner := bufio.NewScanner(r)
	// Increase buffer size for long lines
	const maxCapacity = 10 * 1024 * 1024
	buf := make([]byte, 64*1024)
	scanner.Buffer(buf, maxCapacity)

	for scanner.Scan() {
		line := scanner.Bytes()
		// Keep a copy because scanner.Bytes() is reused
		lineCopy := make([]byte, len(line))
		copy(lineCopy, line)

		// separatedOnly: skip if delimiter absent
		if p.SeparatedOnly && !bytes.Contains(lineCopy, p.Delimiter) {
			continue
		}

		// Split and select fields
		fields := splitFields(lineCopy, p.Delimiter)

		// Collect selected fields in order of appearance by index
		// We keep original field order, selecting those included by FieldSelection.
		selected := make([][]byte, 0, len(fields))
		for i, f := range fields {
			if p.FieldSelection.Includes(i + 1) {
				// Include empty fields as empty tokens
				selected = append(selected, f)
			}
		}

		// If no fields selected, print empty line
		var out []byte
		if len(selected) == 0 {
			out = []byte{}
		} else {
			out = bytes.Join(selected, p.Delimiter)
		}
		if _, err := w.Write(out); err != nil {
			return err
		}
		if _, err := w.Write([]byte{'\n'}); err != nil {
			return err
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func splitFields(line []byte, delim []byte) [][]byte {
	if len(delim) == 1 {
		return bytes.Split(line, delim)
	}
	// For multi-byte delimiter, perform manual split preserving empty fields.
	var res [][]byte
	from := 0
	for {
		idx := bytes.Index(line[from:], delim)
		if idx < 0 {
			res = append(res, line[from:])
			break
		}
		idx += from
		res = append(res, line[from:idx])
		from = idx + len(delim)
		if from > len(line) {
			res = append(res, []byte{})
			break
		}
	}
	return res
}
