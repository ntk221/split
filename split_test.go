package main_test

import (
	"errors"
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
	tests := map[string]struct {
		// name         string
		input        string
		option       option.Command
		outputPrefix string
		wantData     string
	}{
		"twoLines":       {input: "Line1\nLine2\n", option: option.NewLineCount(1), outputPrefix: "x", wantData: "twoLines"},
		"onlyNewLine":    {input: "\n\n\n", option: option.NewLineCount(1), outputPrefix: "x", wantData: "onlyNewLine"},
		"outputPrefix":   {input: "hello\n", option: option.NewLineCount(1), outputPrefix: "HOGE", wantData: "outputPrefix"},
		"largeLineCount": {input: "hello\n", option: option.NewLineCount(math.MaxInt), outputPrefix: "x", wantData: "bigIntCount"},
	}

	var s *splitter.Splitter
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			dir := t.TempDir()
			option := tt.option
			input := strings.NewReader(tt.input)
			s = splitter.New(tt.outputPrefix)

			cli := &splitter.CLI{
				Input:     input,
				OutputDir: dir,
				Splitter:  s,
			}

			err := cli.Run(option)
			if err != nil {
				t.Fatal(err)
			}

			got := golden.Txtar(t, dir)

			if diff := golden.Check(t, flagUpdate, "testdata/lineCount", tt.wantData, got); diff != "" {
				t.Errorf("Test case %s failed:\n%s", name, diff)
			}
		})
	}
}

func TestSplitUsingByteCount(t *testing.T) {
	tests := map[string]struct {
		// name         string
		input        string
		option       option.Command
		outputPrefix string
		wantData     string
	}{
		"simpleCase":  {input: "HogeHogeHugaHuga", option: option.NewByteCount("4"), outputPrefix: "x", wantData: "simple"},
		"indivisible": {input: "Hi,HowAreYou", option: option.NewByteCount("5"), outputPrefix: "x", wantData: "indivisible"},
		"zeroDivided": {input: "hello\n", option: option.NewByteCount("0"), outputPrefix: "x", wantData: "zeroDivided"},
	}

	var s *splitter.Splitter
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			dir := t.TempDir()
			option := tt.option
			input := strings.NewReader(tt.input)
			s = splitter.New(tt.outputPrefix)

			cli := &splitter.CLI{
				Input:     input,
				OutputDir: dir,
				Splitter:  s,
			}

			err := cli.Run(option)
			if err != nil {
				t.Fatal(err)
			}

			got := golden.Txtar(t, dir)

			if diff := golden.Check(t, flagUpdate, "testdata/byteCount", tt.wantData, got); diff != "" {
				t.Errorf("Test case %s failed:\n%s", name, diff)
			}
		})
	}
}

func TestSplitUsingChunkCount(t *testing.T) {
	tests := map[string]struct {
		input        string
		option       option.Command
		outputPrefix string
		wantData     string
		wantErr      error
	}{
		"simpleCase":        {input: "HogeHogeHugaHuga", option: option.NewChunkCount(4), outputPrefix: "x", wantData: "simple"},
		"indivisible":       {input: "Hi,HowAreYou", option: option.NewChunkCount(3), outputPrefix: "x", wantData: "indivisible"},
		"tooManyChunkCount": {input: "hello\n", option: option.NewChunkCount(100), outputPrefix: "x", wantErr: splitter.ErrZeroChunk},
	}

	var s *splitter.Splitter
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			dir := t.TempDir()
			option := tt.option
			input := strings.NewReader(tt.input)
			s = splitter.New(tt.outputPrefix)

			cli := &splitter.CLI{
				Input:     input,
				OutputDir: dir,
				Splitter:  s,
			}

			err := cli.Run(option)
			if err != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("test case %s: 想定されたエラーではありませんでした。", name)
				}
				return
			}

			got := golden.Txtar(t, dir)

			if diff := golden.Check(t, flagUpdate, "testdata/chunkCount", tt.wantData, got); diff != "" {
				t.Errorf("Test case %s failed:\n%s", name, diff)
			}
		})
	}
}
