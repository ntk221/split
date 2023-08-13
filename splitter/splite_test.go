package splitter_test

import (
	"bytes"
	"fmt"
	"github.com/ntk221/split/commandOption"
	"github.com/ntk221/split/splitter"
	"testing"
)

const (
	DefaultPrefix = "x"
)

type FakeFileCreator struct {
	fakeWriteCloser *FakeWriteCloser
	createCount     int
}

func NewFakeFileCreator(closer *FakeWriteCloser) *FakeFileCreator {
	return &FakeFileCreator{
		closer,
		0,
	}
}

func (fc *FakeFileCreator) Create(name string) (splitter.StringWriteCloser, error) {
	fc.createCount++
	return fc.fakeWriteCloser, nil
}

type FakeWriteCloser struct {
	byteContent   [][]byte
	stringContent []string
	writeCount    int
	closeCount    int
}

func NewFakeWriteCloser() *FakeWriteCloser {
	return &FakeWriteCloser{
		byteContent:   make([][]byte, 0),
		stringContent: make([]string, 0),
		writeCount:    0,
		closeCount:    0,
	}
}

func (fw *FakeWriteCloser) Write(p []byte) (int, error) {
	fw.byteContent = append(fw.byteContent, p)
	return len(p), nil

}

func (fw *FakeWriteCloser) WriteString(s string) (int, error) {
	fw.stringContent = append(fw.stringContent, s)
	return len(s), nil
}

func (fw *FakeWriteCloser) Close() error {
	fw.closeCount++
	return nil
}

func TestSplitUsingLineCountNormalCase(t *testing.T) {
	option := commandOption.NewLineCountOption(3)
	outputPrefix := DefaultPrefix

	fileContent := []byte("Line 1\nLine 2\nLine 3\nLine 4\nLine 5\nLine 6\nLine 7\nLine 8\nLine 9\nLine 10\nLine 11\n")
	fileReader := bytes.NewReader(fileContent)

	fakeWriteCloser := NewFakeWriteCloser()
	fakeFileCreator := NewFakeFileCreator(fakeWriteCloser)

	s := splitter.NewSplitter(option, outputPrefix, fileReader, fakeFileCreator)
	s.SplitUsingLineCount()

	// ちゃんと全部の行を読んでる
	i := 1
	for _, string := range fakeWriteCloser.stringContent {
		expect := fmt.Sprintf("Line %d\n", i)
		if expect != string {
			t.Fatalf("exepct :%s, but actual :%s", expect, string)
		}
		i++
	}

	// 12 / 3 == 4 回 Createしてる
	if fakeFileCreator.createCount != 4 {
		t.Fatalf("createCount should be %d, but actual %d\n", 4, fakeFileCreator.createCount)
	}

	// 4回Closeしている
	if fakeWriteCloser.closeCount != 4 {
		t.Fatalf("closeCount should be %d, but actual %d\n", 4, fakeWriteCloser.closeCount)
	}

}

func TestSplitUsingLineCountNoNewLine(t *testing.T) {
	option := commandOption.NewLineCountOption(3)
	outputPrefix := DefaultPrefix

	fileContent := []byte("Line 1Line 2Line 3Line 4Line 5Line 6Line 7Line 8Line 9Line 10Line 11")
	fileReader := bytes.NewReader(fileContent)

	fakeWriteCloser := NewFakeWriteCloser()
	fakeFileCreator := NewFakeFileCreator(fakeWriteCloser)

	s := splitter.NewSplitter(option, outputPrefix, fileReader, fakeFileCreator)
	s.SplitUsingLineCount()

	if fakeWriteCloser.writeCount != 0 {
		t.Fatalf("if file content don't include newline, write count should be %d, but actual %d", 0, fakeWriteCloser.writeCount)
	}

	// Createは一回
	if fakeFileCreator.createCount != 1 {
		t.Fatalf("createCount should be %d, but actual %d\n", 1, fakeFileCreator.createCount)
	}

	// Closeも一回
	if fakeWriteCloser.closeCount != 1 {
		t.Fatalf("closeCount should be %d, but actual %d\n", 1, fakeWriteCloser.closeCount)
	}
}

func TestSplitUsingLineCountEmptyFile(t *testing.T) {
	option := commandOption.NewLineCountOption(3)
	outputPrefix := DefaultPrefix

	fileContent := []byte("")
	fileReader := bytes.NewReader(fileContent)

	fakeWriteCloser := NewFakeWriteCloser()
	fakeFileCreator := NewFakeFileCreator(fakeWriteCloser)

	s := splitter.NewSplitter(option, outputPrefix, fileReader, fakeFileCreator)
	s.SplitUsingLineCount()

	// 何も書き込まない
	if len(fakeWriteCloser.stringContent) != 0 {
		t.Fatalf("stringContent should be empty, but it's not")
	}

	// file Createはする
	if fakeFileCreator.createCount != 1 {
		t.Fatalf("createCount should be %d, but actual %d\n", 1, fakeFileCreator.createCount)
	}

	// Closeもちゃんとする
	if fakeWriteCloser.closeCount != 1 {
		t.Fatalf("closeCount should be %d, but actual %d\n", 1, fakeWriteCloser.closeCount)
	}
}

func TestSplitUsingLineCountLessLinesThanOption(t *testing.T) {
	option := commandOption.NewLineCountOption(20) // 行数よりも大きい値
	outputPrefix := DefaultPrefix

	fileContent := []byte("Line 1\nLine 2\nLine 3")
	fileReader := bytes.NewReader(fileContent)

	fakeWriteCloser := NewFakeWriteCloser()
	fakeFileCreator := NewFakeFileCreator(fakeWriteCloser)

	s := splitter.NewSplitter(option, outputPrefix, fileReader, fakeFileCreator)
	s.SplitUsingLineCount()

	i := 1
	for _, string := range fakeWriteCloser.stringContent {
		expect := fmt.Sprintf("Line %d\n", i)
		if expect != string {
			t.Fatalf("exepct :%s, but actual :%s", expect, string)
		}
		i++
	}

	// 作成されるファイル数は 1 つのみ
	if fakeFileCreator.createCount != 1 {
		t.Fatalf("createCount should be %d, but actual %d\n", 1, fakeFileCreator.createCount)
	}

	// Closeも一回
	if fakeWriteCloser.closeCount != 1 {
		t.Fatalf("closeCount should be %d, but actual %d\n", 1, fakeWriteCloser.closeCount)
	}
}

func TestSplitUsingChunkCountSingleChunk(t *testing.T) {
	chunkCount := 1
	fileContent := []byte("Line 1\nLine 2\nLine 3\nLine 4\nLine 5\nLine 6\nLine 7\nLine 8\nLine 9\nLine 10\nLine 11\n")
	expectedChunks := [][]byte{
		[]byte("Line 1\nLine 2\nLine 3\nLine 4\nLine 5\nLine 6\nLine 7\nLine 8\nLine 9\nLine 10\nLine 11\n"),
	}

	option := commandOption.NewChunkCountOption(chunkCount)
	outputPrefix := DefaultPrefix
	fileReader := bytes.NewReader(fileContent)
	fakeWriteCloser := NewFakeWriteCloser()
	fakeFileCreator := NewFakeFileCreator(fakeWriteCloser)
	s := splitter.NewSplitter(option, outputPrefix, fileReader, fakeFileCreator)
	s.SplitUsingChunkCount()

	// Compare the chunks written by fakeWriteCloser
	if len(fakeWriteCloser.byteContent) != len(expectedChunks) {
		t.Fatalf("expected %d chunks, but got %d", len(expectedChunks), len(fakeWriteCloser.byteContent))
	}

	for i, expectedChunk := range expectedChunks {
		if !bytes.Equal(fakeWriteCloser.byteContent[i], expectedChunk) {
			t.Errorf("chunk mismatch for chunk %d\nExpected:\n%s\nActual:\n%s\n",
				i, string(expectedChunk), string(fakeWriteCloser.byteContent[i]))
		}
	}
}

func TestSplitUsingChunkCountSomeChunk(t *testing.T) {
	chunkCount := 3
	fileContent := []byte("aaaaaaaaaaaa")
	expectedChunks := [][]byte{
		[]byte("aaaa"),
		[]byte("aaaa"),
		[]byte("aaaa"),
	}

	option := commandOption.NewChunkCountOption(chunkCount)
	outputPrefix := DefaultPrefix
	fileReader := bytes.NewReader(fileContent)
	fakeWriteCloser := NewFakeWriteCloser()
	fakeFileCreator := NewFakeFileCreator(fakeWriteCloser)

	s := splitter.NewSplitter(option, outputPrefix, fileReader, fakeFileCreator)
	s.SplitUsingChunkCount()

	// Compare the chunks written by fakeWriteCloser
	if len(fakeWriteCloser.byteContent) != len(expectedChunks) {
		t.Fatalf("expected %d chunks, but got %d", len(expectedChunks), len(fakeWriteCloser.byteContent))
	}

	for i, expectedChunk := range expectedChunks {
		if !bytes.Equal(fakeWriteCloser.byteContent[i], expectedChunk) {
			t.Errorf("chunk mismatch for chunk %d\nExpected:\n%s\nActual:\n%s\n",
				i, string(expectedChunk), string(fakeWriteCloser.byteContent[i]))
		}
	}
}

func TestSplitUsingChunkCountCantDivideContent(t *testing.T) {
	chunkCount := 3
	fileContent := []byte("aaaaaaaaaaa") // 11バイト
	// 11 / 3 == 3
	expectedChunks := [][]byte{
		[]byte("aaa"),
		[]byte("aaa"),
		[]byte("aaaaa"), // 最後は残りを全部詰める
	}

	option := commandOption.NewChunkCountOption(chunkCount)
	outputPrefix := DefaultPrefix
	fileReader := bytes.NewReader(fileContent)
	fakeWriteCloser := NewFakeWriteCloser()
	fakeFileCreator := NewFakeFileCreator(fakeWriteCloser)
	s := splitter.NewSplitter(option, outputPrefix, fileReader, fakeFileCreator)
	s.SplitUsingChunkCount()

	// Compare the chunks written by fakeWriteCloser
	if len(fakeWriteCloser.byteContent) != len(expectedChunks) {
		t.Fatalf("expected %d chunks, but got %d", len(expectedChunks), len(fakeWriteCloser.byteContent))
	}

	for i, expectedChunk := range expectedChunks {
		if !bytes.Equal(fakeWriteCloser.byteContent[i], expectedChunk) {
			t.Errorf("chunk mismatch for chunk %d\nExpected:\n%s\nActual:\n%s\n",
				i, string(expectedChunk), string(fakeWriteCloser.byteContent[i]))
		}
	}
}

func TestSplitUsingChunkCountEmptyContent(t *testing.T) {
	chunkCount := 3
	fileContent := []byte("") // 空のcontent
	// 11 / 3 == 3
	expectedChunks := [][]byte{}

	option := commandOption.NewChunkCountOption(chunkCount)
	outputPrefix := DefaultPrefix
	fileReader := bytes.NewReader(fileContent)
	fakeWriteCloser := NewFakeWriteCloser()
	fakeFileCreator := NewFakeFileCreator(fakeWriteCloser)
	s := splitter.NewSplitter(option, outputPrefix, fileReader, fakeFileCreator)
	s.SplitUsingChunkCount()

	// Compare the chunks written by fakeWriteCloser
	if len(fakeWriteCloser.byteContent) != len(expectedChunks) {
		t.Fatalf("expected %d chunks, but got %d", len(expectedChunks), len(fakeWriteCloser.byteContent))
	}

	for i, expectedChunk := range expectedChunks {
		if !bytes.Equal(fakeWriteCloser.byteContent[i], expectedChunk) {
			t.Errorf("chunk mismatch for chunk %d\nExpected:\n%s\nActual:\n%s\n",
				i, string(expectedChunk), string(fakeWriteCloser.byteContent[i]))
		}
	}
}

func TestSplitUsingChunkCountOverContentSize(t *testing.T) {
	chunkCount := 20
	fileContent := []byte("aaaaaa")

	option := commandOption.NewChunkCountOption(chunkCount)
	outputPrefix := DefaultPrefix
	fileReader := bytes.NewReader(fileContent)
	fakeWriteCloser := NewFakeWriteCloser()
	fakeFileCreator := NewFakeFileCreator(fakeWriteCloser)

	s := splitter.NewSplitter(option, outputPrefix, fileReader, fakeFileCreator)
	s.SplitUsingChunkCount()

	if fakeFileCreator.createCount != 1 {
		t.Fatalf("file create should be only 1 time, but actual %d times", fakeFileCreator.createCount)
	}
}

func TestSplitUsingByteCount(t *testing.T) {
	tests := []struct {
		name           string
		byteCount      string
		fileContent    []byte
		expectedChunks [][]byte
	}{
		{
			name:        "SingleChunk",
			byteCount:   "100B",
			fileContent: []byte("This is a test content for splitting by byte count."),
			expectedChunks: [][]byte{
				[]byte("This is a test content for splitting by byte count."),
			},
		},
		{
			name:        "MultipleChunks",
			byteCount:   "20b",
			fileContent: []byte("This is a test content for splitting by byte count."),
			expectedChunks: [][]byte{
				[]byte("This is a test "),
				[]byte("content for spli"),
				[]byte("tting by byte cou"),
				[]byte("nt."),
			},
		},
		// Add more test cases as needed...
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			option := commandOption.NewByteCountOption(test.byteCount)
			outputPrefix := DefaultPrefix

			fileReader := bytes.NewReader(test.fileContent)

			fakeWriteCloser := NewFakeWriteCloser()
			fakeFileCreator := NewFakeFileCreator(fakeWriteCloser)

			s := splitter.NewSplitter(option, outputPrefix, fileReader, fakeFileCreator)
			s.SplitUsingByteCount()

			// Compare the chunks written by fakeWriteCloser
			if len(fakeWriteCloser.byteContent) != len(test.expectedChunks) {
				t.Fatalf("expected %d chunks, but got %d", len(test.expectedChunks), len(fakeWriteCloser.byteContent))
			}

			for i, expectedChunk := range test.expectedChunks {
				if !bytes.Equal(fakeWriteCloser.byteContent[i], expectedChunk) {
					t.Errorf("chunk mismatch for chunk %d\nExpected:\n%s\nActual:\n%s\n",
						i, string(expectedChunk), string(fakeWriteCloser.byteContent[i]))
				}
			}
		})
	}
}
