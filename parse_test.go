package tfjson

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"
)

const testFixtureDir = "test-fixtures"
const testGoldenFileName = "plan.json"

func TestParse(t *testing.T) {
	entries, err := ioutil.ReadDir(testFixtureDir)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	for _, e := range entries {
		if !e.IsDir() {
			continue
		}

		t.Run(e.Name(), func(t *testing.T) {
			expected, err := ioutil.ReadFile(filepath.Join(testFixtureDir, e.Name(), testGoldenFileName))
			if err != nil {
				t.Fatal(err)
			}

			var parsed Plan
			if err = json.Unmarshal(expected, &parsed); err != nil {
				t.Fatal(err)
			}

			actual, err := json.MarshalIndent(parsed, "", "  ")
			if err != nil {
				t.Fatal(err)
			}

			// Add a newline at the end
			actual = append(actual, byte('\n'))

			if err := testDiff(actual, expected); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func testDiff(out, gld []byte) error {
	var b strings.Builder // holding long error message

	// compare lengths
	if len(out) != len(gld) {
		fmt.Fprintf(&b, "\nlength changed: len(output) = %d, len(golden) = %d", len(out), len(gld))
	}

	// compare contents
	line := 1
	offs := 1
	for i := 0; i < len(out) && i < len(gld); i++ {
		ch := out[i]
		if ch != gld[i] {
			fmt.Fprintf(&b, "\noutput:%d:%d: %s", line, i-offs+1, lineAt(out, offs))
			fmt.Fprintf(&b, "\ngolden:%d:%d: %s", line, i-offs+1, lineAt(gld, offs))
			fmt.Fprintf(&b, "\n\n")
			break
		}
		if ch == '\n' {
			line++
			offs = i + 1
		}
	}

	if b.Len() > 0 {
		return errors.New(b.String())
	}
	return nil
}

func lineAt(text []byte, offs int) []byte {
	i := offs
	for i < len(text) && text[i] != '\n' {
		i++
	}
	return text[offs:i]
}