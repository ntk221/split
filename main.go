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
	usage:	split [-l line_count] [file [prefix]]
		split -b byte_count[K|k|M|m|G|g] [file [prefix]]
		split -n chunk_count [file [prefix]]
	`
)

var (
	lineCountOption  = flag.Int("l", 1000, "行数を指定してください")
	chunkCountOption = flag.Int("n", 0, "chunkCountを指定してください")
	// byteCountOption  = flag.String("b", "", "バイト数を指定してください（例: 10K, 2M, 3G）")
)

// 対応するoptionを-n, -l, -bのみに制限した場合、optionのvalidationについては複数指定されているか否かで判定できる
// この条件を満たさない場合はvalidationについての実装が異なるので拡張する場合は要変更
func validateOptions() {
	optionCount := 0
	flag.VisitAll(func(f *flag.Flag) {
		if f.Value.String() != f.DefValue {
			optionCount++
		}
	})

	if optionCount > 1 {
		log.Fatal(Synopsys)
	}
}

func split(args []string, file *os.File) {
	// optionによって処理を分岐する
	flag.Parse()
	lineCount := *lineCountOption
	// chunkCount := *chunkCountOption
	// byteCount := *byteCountOption

	validateOptions()

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
