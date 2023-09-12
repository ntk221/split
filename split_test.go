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
		input        string
		option       option.Command
		outputPrefix string
		wantData     string
	}{
		"twoLines":       {"Line1\nLine2\n", lineCount(t, 1), "x", "twoLines"},
		"onlyNewLine":    {"\n\n\n", lineCount(t, 1), "x", "onlyNewLine"},
		"outputPrefix":   {"hello\n", lineCount(t, 1), "HOGE", "outputPrefix"},
		"largeLineCount": {"hello\n", lineCount(t, math.MaxInt), "x", "bigIntCount"},
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
		"simpleCase":  {"HogeHogeHugaHuga", byteCount(t, "4"), "x", "simple"},
		"indivisible": {"Hi,HowAreYou", byteCount(t, "5"), "x", "indivisible"},
		"zeroDivided": {"hello\n", byteCount(t, "0"), "x", "zeroDivided"},
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
		expectErr    error
	}{
		"simpleCase":        {"HogeHogeHugaHuga", chunkCount(t, 4), "x", "simple", nil},
		"indivisible":       {"Hi,HowAreYou", chunkCount(t, 3), "x", "indivisible", nil},
		"tooManyChunkCount": {"hello\n", chunkCount(t, 100), "x", "", splitter.ErrZeroChunk},
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
				if !errors.Is(err, tt.expectErr) {
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

func lineCount(t *testing.T, n int) option.Command {
	t.Helper()

	return option.NewLineCount(n)
}

func chunkCount(t *testing.T, n int) option.Command {
	t.Helper()

	return option.NewChunkCount(n)
}

func byteCount(t *testing.T, b string) option.Command {
	t.Helper()

	return option.NewByteCount(b)
}
