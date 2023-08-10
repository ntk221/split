package main

import (
	"flag"
	"fmt"
	"github.com/ntk221/split/splitter"
	"log"
	"os"
)

const (
	DefaultPrefix = "x"
	Synopsys      = `
	usage:	split -d [-l line_count] [-a suffix_length] [file [prefix]]
		split -d -b byte_count[K|k|M|m|G|g] [-a suffix_length] [file [prefix]]
		split -d -n chunk_count [-a suffix_length] [file [prefix]]
		split -d -p pattern [-a suffix_length] [file [prefix]]
	`
)

var (
	lineCountOption  = flag.Int("l", 1000, "行数を指定してください")
	chunkCountOption = flag.Int("n", 0, "chunkCountを指定してください")
	// byteCountOption  = flag.String("b", "", "バイト数を指定してください（例: 10K, 2M, 3G）")
)

func split(args []string, file *os.File) {
	// optionによって処理を分岐する
	flag.Parse()
	lineCount := *lineCountOption
	// chunkCount := *chunkCountOption
	// byteCount := *byteCountOption

	optionCount := 0
	flag.VisitAll(func(f *flag.Flag) {
		if f.Value.String() != f.DefValue {
			optionCount++
		}
	})

	if optionCount > 1 {
		log.Fatal(Synopsys)
	}

	partNum := 0
	outputPrefix := DefaultPrefix
	// 引数でprefixが指定されている場合はそれを使う
	if len(args) > 1 {
		outputPrefix = args[1]
	}

	// chunkCountが設定されてたら -> splitUsingChunkCount
	// byteCountが設定されてたら -> splitUsingByteCount

	splitter.SplitUsingLineCount(lineCount, outputPrefix, partNum, file)
}

func main() {
	flag.Parse()
	args := flag.Args()
	fmt.Println(args)

	// 今は標準入力から受け取る機能がないが、実装する
	if len(args) == 0 {
		log.Fatal("TODO")
	}

	fileName := args[0]
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	split(args, file)
}
