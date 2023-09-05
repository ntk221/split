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

func TestSplit(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		option       option.Command
		outputPrefix string
		wantData     string
	}{
		{
			name:         "改行含む2行のテスト",
			input:        "Line1\nLine2\n",
			option:       option.NewLineCount(1),
			outputPrefix: "x",
			wantData:     "twoLines",
		},
		/*{
			name:         "1行",
			input:        "test",
			option:       option.NewLineCount(1),
			outputPrefix: "x",
			wantData:     "noNewLine",
		},*/
	}

	var s *splitter.Splitter
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			option := tt.option
			input := strings.NewReader(tt.input)
			s = splitter.New(option, tt.outputPrefix)

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
