package main_test

import (
	"flag"
	"github.com/ntk221/split/option"
	"github.com/ntk221/split/splitter"
	"github.com/tenntenn/golden"
	"strings"
	"testing"
)

var (
	flagUpdate bool
)

func init() {
	flag.BoolVar(&flagUpdate, "update", false, "update golden files")
}

func TestSplitterCLI(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		lineCount int
		expected  string
	}{
		{
			name:      "Test Case 1",
			input:     "Line1\nLine2\n",
			lineCount: 1,
			expected:  "expected_output_1.txtar",
		},
		// Add more test cases here as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			option := option.LineCount(tt.lineCount)
			input := strings.NewReader(tt.input)
			s := splitter.New(option, "x")

			cli := &splitter.CLI{
				Input:     input,
				OutputDir: dir,
				Splitter:  s,
			}

			cli.Run()

			got := golden.Txtar(t, dir)

			if diff := golden.Check(t, flagUpdate, "testdata", tt.expected, got); diff != "" {
				t.Errorf("Test case %s failed:\n%s", tt.name, diff)
			}
		})
	}
}

func Test(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		option       option.Command
		outputPrefix string
		wantData     string
	}{
		{
			name:         "Test Case 1",
			input:        "Line1\nLine2\n",
			option:       option.NewLineCount(1),
			outputPrefix: "x",
			wantData:     "mytest",
		},
		// Add more test cases here as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			option := tt.option
			input := strings.NewReader(tt.input)
			s := splitter.New(option, tt.outputPrefix)

			cli := &splitter.CLI{
				Input:     input,
				OutputDir: dir,
				Splitter:  s,
			}

			cli.Run()

			got := golden.Txtar(t, dir)

			if diff := golden.Check(t, flagUpdate, "testdata", tt.wantData, got); diff != "" {
				t.Errorf("Test case %s failed:\n%s", tt.name, diff)
			}
		})
	}
}
