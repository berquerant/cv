package main_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEndToEnd(t *testing.T) {
	e := newExecutor(t)
	defer e.close()

	if err := run(e.cmd, "-h"); err != nil {
		t.Fatalf("%s help %v", e.cmd, err)
	}

	conv := func(origin, in, out string) (string, error) {
		cmd := exec.Command(e.cmd, in, out)
		cmd.Stdin = bytes.NewBufferString(origin)
		b, err := cmd.Output()
		return string(b), err
	}

	const (
		mapOrigin     = `{"a":1,"b":["b1","b2"],"c":{"d":2}}`
		mapOriginYaml = `a: 1
b:
    - b1
    - b2
c:
    d: 2
`
		mapOriginToml = `a = 1.0
b = ['b1', 'b2']

[c]
d = 2.0
`

		simpleMapJSON      = `{"a":"1","b":"bv","c":"100"}`
		simpleMapLTSV      = "a:1\tb:bv\tc:100\n"
		simpleSliceMapJSON = `[{"a":"1","b":"bv","c":"100"}]`

		simpleSliceJSON       = `["a","b","c"]`
		simpleCSV             = "a,b,c\n"
		simpleNestedSliceJSON = `[["a","b","c"]]`
	)
	for _, tc := range []struct {
		name     string
		input    string
		in       string
		out      string
		wantConv string
		wantInv  string
	}{
		{
			name:     "simple-slice-json-csv",
			input:    simpleSliceJSON,
			in:       "json",
			out:      "csv",
			wantConv: simpleCSV,
			wantInv:  simpleNestedSliceJSON,
		},
		{
			name:     "simple-nested-slice-json-csv",
			input:    simpleNestedSliceJSON,
			in:       "json",
			out:      "csv",
			wantConv: simpleCSV,
			wantInv:  simpleNestedSliceJSON,
		},
		{
			name:     "simple-json-ltsv",
			input:    simpleMapJSON,
			in:       "json",
			out:      "ltsv",
			wantConv: simpleMapLTSV,
			wantInv:  simpleSliceMapJSON,
		},
		{
			name:     "simple-slice-json-ltsv",
			input:    simpleSliceMapJSON,
			in:       "json",
			out:      "ltsv",
			wantConv: simpleMapLTSV,
			wantInv:  simpleSliceMapJSON,
		},
		{
			name:     "json-json-map",
			input:    mapOrigin,
			in:       "json",
			out:      "json",
			wantConv: mapOrigin,
			wantInv:  mapOrigin,
		},
		{
			name:     "json-yaml-map",
			input:    mapOrigin,
			in:       "json",
			out:      "yaml",
			wantConv: mapOriginYaml,
			wantInv:  mapOrigin,
		},
		{
			name:     "json-toml-map",
			input:    mapOrigin,
			in:       "json",
			out:      "toml",
			wantConv: mapOriginToml,
			wantInv:  mapOrigin,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c, err := conv(tc.input, tc.in, tc.out)
			assert.Nil(t, err)
			assert.Equal(t, tc.wantConv, c)
			r, err := conv(tc.wantConv, tc.out, tc.in)
			assert.Nil(t, err)
			assert.Equal(t, tc.wantInv, r)
		})
	}
}

func run(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	cmd.Dir = "."
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

type executor struct {
	dir string
	cmd string
}

func newExecutor(t *testing.T) *executor {
	t.Helper()
	e := &executor{}
	e.init(t)
	return e
}

func (e *executor) init(t *testing.T) {
	t.Helper()
	dir, err := os.MkdirTemp("", "cv")
	if err != nil {
		t.Fatal(err)
	}
	cmd := filepath.Join(dir, "cv")
	// build grinfo command
	if err := run("go", "build", "-o", cmd); err != nil {
		t.Fatal(err)
	}
	e.dir = dir
	e.cmd = cmd
}

func (e *executor) close() {
	os.RemoveAll(e.dir)
}
