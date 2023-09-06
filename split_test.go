package main_test

import (
	"flag"
	"github.com/ntk221/split/option"
	"github.com/ntk221/split/splitter"
	"github.com/tenntenn/golden"
	"math"
	"strings"
	"testing"
)

var (
	flagUpdate bool
)

func init() {
	flag.BoolVar(&flagUpdate, "update", false, "update golden files")
}

func TestSplitUsingLineCount(t *testing.T) {
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
		{
			name:         "改行のみ",
			input:        "\n\n\n",
			option:       option.NewLineCount(1),
			outputPrefix: "x",
			wantData:     "onlyNewLine",
		},
		{
			name:         "output prefix指定",
			input:        "hello\n",
			option:       option.NewLineCount(1),
			outputPrefix: "HOGE",
			wantData:     "outputPrefix",
		},
		{
			name:         "line countをINT_MAXで指定",
			input:        "hello\n",
			option:       option.NewLineCount(math.MaxInt),
			outputPrefix: "x",
			wantData:     "bigIntCount",
		},
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

			if diff := golden.Check(t, flagUpdate, "testdata/lineCount", tt.wantData, got); diff != "" {
				t.Errorf("Test case %s failed:\n%s", tt.name, diff)
			}
		})
	}
}

func TestSplitUsingByteCount(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		option       option.Command
		outputPrefix string
		wantData     string
	}{
		{
			name:         "改行含む2行のテスト",
			input:        "HogeHogeHugaHuga",
			option:       option.NewByteCount("4"),
			outputPrefix: "x",
			wantData:     "simple",
		},
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

			if diff := golden.Check(t, flagUpdate, "testdata/byteCount", tt.wantData, got); diff != "" {
				t.Errorf("Test case %s failed:\n%s", tt.name, diff)
			}
		})
	}
}
