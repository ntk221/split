package main

import (
	"flag"
	"fmt"
	"github.com/ntk221/split/splitter"
	"log"
	"net/http"
	"os"

	. "github.com/ntk221/split/commandOption"
)

const (
	DefaultPrefix = "x"
	Synopsys      = `
	usage:	split [-l line_count] [file [prefix]]
		split -b byte_count[K|k|M|m|G|g] [file [prefix]]
		split -n chunk_count [file [prefix]]`
)

var (
	lineCountOption  = flag.Int("l", 1000, "行数を指定してください")
	chunkCountOption = flag.Int("n", 0, "chunk数を指定してください")
	byteCountOption  = flag.String("b", "", "バイト数を指定してください（例: 10K, 2M, 3G）")
)

// プログラムの実行例: ./split -l 2 test.txt
// flag packageを使った際のoptionの指定方法が option + space + value という形式しか発見できなかったため、 ./split -l2 test.txt のように space を開けない実行が未実装
// textファイルではない入力に関して実験をしたが、splitコマンドと挙動が揃わず、man にもそれについて説明がなかったので
// 本プログラムの仕様として、text ファイル以外の入力を受け取らないようにした
func main() {
	flag.Parse()
	args := flag.Args()

	file, closeFile := readyFile(args)
	defer closeFile()
	detectFileType(file)

	// 1. ファイル名が指定されていて、かつ、オプション指定されている時には
	// コマンドライン引数の先頭はオプションであるべきである
	// 2. ファイル名が指定されていて、かつ、オプション指定されている時には
	if commandLineArgs := os.Args; len(commandLineArgs) > 2 && len(args) < 2 {
		first := commandLineArgs[1]
		if first != "-l" && first != "-n" && first != "-b" {
			log.Fatal(Synopsys)
		}
		validateOptions()
	}

	// optionの取得
	lineCount := NewLineCountOption(*lineCountOption)
	chunkCount := NewChunkCountOption(*chunkCountOption)
	byteCount := NewByteCountOption(*byteCountOption)

	options := make([]CommandOption, 0)
	options = append(options, lineCount)
	options = append(options, chunkCount)
	options = append(options, byteCount)

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

// コマンドライン引数でファイル名を指定された場合はそれをオープンして返す
// コマンドライン引数が指定されない場合は、標準入力から受け取る
func readyFile(args []string) (*os.File, func()) {
	if len(args) == 0 {
		return os.Stdin, func() {}
	}

	fileName := args[0]
	file, err := os.OpenFile(fileName, os.O_RDONLY, 0666)
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatal(fmt.Sprintf("split: %s: no such file or directory", fileName))
		}
		log.Fatal(err)
		return nil, func() {}
	}

	return file, func() { file.Close() }
}

// http.DetectContentType を呼び出してファイルの種類を特定し、
// textファイルかapplication/octet-stream(DetectContentTypeのデフォルトの値)でない時にプログラムを終了する
func detectFileType(file *os.File) {
	defer func() {
		// 処理が終了したら元の位置に戻す
		_, err := file.Seek(0, 0)
		if err != nil {
			log.Fatal(err)
		}
	}()

	// http.DetectContentType で file type を特定するには, 最大で512バイト必要になる
	// 参考url: https://pkg.go.dev/net/http#DetectContentType
	buf := make([]byte, 512)
	_, err := file.Read(buf)
	if err != nil {
		log.Fatal(err)
	}

	fileType := http.DetectContentType(buf)

	if fileType != "application/octet-stream" && fileType[:5] != "text/" {
		log.Fatal("このプログラムは、textファイルのみを入力として受け取ります")
	}
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
