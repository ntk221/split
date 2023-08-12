package splitter_test

import (
	"github.com/ntk221/split/commandOption"
	"github.com/ntk221/split/splitter"
	"os"
	"testing"
)

const (
	DefaultPrefix = "x"
)

func TestSplitUsingLineCount(t *testing.T) {
	option := commandOption.NewLineCountOption(10)
	outputPrefix := DefaultPrefix

	tempDir := t.TempDir()
	fileContent := []byte("Line 1\nLine 2\nLine 3\nLine 4\nLine 5\nLine 6\nLine 7\nLine 8\nLine 9\nLine 10\nLine 11\n")
	testFile, err := os.CreateTemp(tempDir, "testfile.txt")
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}
	_, err = testFile.Write(fileContent)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	s := splitter.NewSplitter(option, outputPrefix, file, func { return os.Create() })
	s.SplitUsingLineCount()
}

// 文字列用のincrement関数
// ex: incrementString("a") == "b"
// ex: incrementString("az") == "ba"
func incrementString(s string) string {
	runes := []rune(s)
	for i := len(runes) - 1; i >= 0; i-- {
		if runes[i] < 'z' {
			runes[i]++
			return string(runes)
		}
		runes[i] = 'a'
	}

	return "a" + string(runes)
}
