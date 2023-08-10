package main

import (
	"flag"
	"github.com/ntk221/split/splitter"
	"log"
	"os"

	. "github.com/ntk221/split/commandOption"
)

const (
	DefaultPrefix = "x"
	Synopsys      = `
	usage:	split [-l line_count] [file [prefix]]
		split -b byte_count[K|k|M|m|G|g] [file [prefix]]
		split -n chunk_count [file [prefix]]
	`
	DefaultChunkCount = 0
	DefaultByteCount  = ""
	DefaultLineCount  = 1000
)

var (
	lineCountOption  = flag.Int("l", 1000, "行数を指定してください")
	chunkCountOption = flag.Int("n", 0, "chunk数を指定してください")
	byteCountOption  = flag.String("b", "", "バイト数を指定してください（例: 10K, 2M, 3G）")
)

func main() {
	flag.Parse()
	args := flag.Args()

	// 今は標準入力から受け取る機能がないが、実装する
	// var file *os.File
	if len(args) == 0 {
		log.Fatal("TODO: 標準入力から読み込む機能を実装する")
	}

	fileName := args[0]
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// optionによって処理を分岐する
	lineCount := NewLineCountOption(*lineCountOption)
	chunkCount := NewChunkCountOption(*chunkCountOption)
	// byteCount := *byteCountOption

	validateOptions()
	options := make([]CommandOption, 0)
	options = append(options, lineCount)
	options = append(options, chunkCount)
	// options = append(options, byteCount)

	// プログラムの引数で指定されたものを選ぶ
	// validateで適切なoptionだけが残っていることを保証している
	option := selectOption(options)
	outputPrefix := DefaultPrefix
	// 引数でprefixが指定されている場合はそれを使う
	if len(args) > 1 {
		outputPrefix = args[1]
	}

	s := splitter.NewSplitter(option, outputPrefix, file)
	s.Split()
	return
}

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

// プログラムの引数として指定されたoptionを返す
// 事前条件: すでにoptionsは適切なものが残っていることが保証されている
func selectOption(options []CommandOption) CommandOption {
	var ret CommandOption = NewLineCountOption(DefaultLineCount)
	for _, o := range options {
		if o.IsDefaultValue() {
			continue
		}
		ret = o
	}
	return ret
}
