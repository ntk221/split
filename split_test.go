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

func Test(t *testing.T) {
	dir := t.TempDir()
	option := option.LineCount(1)
	input := strings.NewReader("Line1\nLine2\n")
	s := splitter.New(option, "x")

	cli := &splitter.CLI{
		Input:     input,
		OutputDir: dir,
		Splitter:  s,
	}

	cli.Run()

	got := golden.Txtar(t, dir) // txtar形式に固める

	if diff := golden.Check(t, flagUpdate, "testdata", "mytest", got); diff != "" { // testdata 以下のmytestとdiffをとる
		t.Error(diff)
	}
}
