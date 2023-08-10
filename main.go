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
)

var (
	lineCountOption = flag.Int("l", 1000, "行数を指定してください")
)

func split(args []string, file *os.File) {
	// optionによって処理を分岐する
	flag.Parse()
	lineCount := *lineCountOption
	// chunkCount := *chunkCountOption
	// byteCount := *byteCountOption

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

	if len(args) == 0 {
		log.Fatal("十分な数のコマンドライン引数が与えられていない")
	}

	fileName := args[0]
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	split(args, file)
}
