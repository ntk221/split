// The split utility reads the given file and breaks it up
// into files of 1000 lines each (if no options are specified), leaving the file unchanged.
// If file is a single dash (‘-’) or absent, split reads from the standard input.
//
// The options are as follows:
//
//		-b byte_count[K|k|M|m|G|g]
//		 Create split files byte_count bytes in length.
//		 If k or K is appended to the number,
//		 the file is split into byte_count kilobyte pieces.
//		 If m or M is appended to the number,
//		 the file is split into byte_count megabyte pieces.
//		 If g or G is appended to the number,
//		 the file is split into byte_count gigabyte pieces.
//
//		-l line_count
//		 Create split files line_count lines in length.
//
//	 	-n chunk_count
//	     Split file into chunk_count smaller files.
//	     The first n - 1 files will be of size (size of file / chunk_count ) and
//	     the last file will contain the remaining bytes.
//
// プログラムの実行例: ./split -l 2 test.txt
//
// flag packageを使った際のoptionの指定方法が option + space + value という形式しか発見できなかった
// よって ./split -l2 test.txt のように space を開けない実行が未実装にした
//
// textファイルではない入力に関する挙動について man にはそれについて説明がなかった
// よって、本プログラムの仕様として、text ファイル以外の入力を受け取らないようにした
//
// text ファイルの判定に file コマンドを使用しているのでこのプログラムはfileコマンドが使える環境でなくては動作しない
package main

import (
	"flag"
	"fmt"
	"github.com/ntk221/split/splitter"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/ntk221/split/option"
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

func main() {
	flag.Parse()
	args := flag.Args()

	file, closeFile := readyFile(args)
	defer closeFile()
	if ok := detectFileType(file); !ok {
		log.Fatal("指定されたファイルはtextファイルではありません")
	}

	// 1. ファイル名が指定されている
	// 2. オプション指定されている
	// 3. 出力ファイルのprefixは指定されていない
	// 以上の条件を満たす時、コマンドライン引数の先頭はオプションであるべきである
	if commandLineArgs := os.Args; len(commandLineArgs) > 2 && len(args) < 2 {
		first := commandLineArgs[1]
		if first != "-l" && first != "-n" && first != "-b" {
			log.Fatal(Synopsys)
		}
		if ok := validateOptions(); !ok {
			log.Fatal(Synopsys)
		}
	}

	options := []option.Command{option.NewLineCount(*lineCountOption), option.NewChunkCount(*chunkCountOption), option.NewByteCount(*byteCountOption)}
	option := selectOption(options)

	outputPrefix := DefaultPrefix
	// 引数でprefixが指定されている場合はそれを使う
	if len(args) > 1 {
		outputPrefix = args[1]
	}

	outputDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		return
	}

	s := splitter.New(option, outputPrefix)

	cli := &splitter.CLI{
		Input:     file,
		OutputDir: outputDir,
		Splitter:  s,
	}

	cli.Run()
	return
}

// コマンドライン引数でファイル名を指定された場合はそれをオープンして返す
// コマンドライン引数が指定されない場合は、標準入力から受け取る
func readyFile(args []string) (file *os.File, close func()) {
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
	}

	close = func() { file.Close() }
	return file, close
}

// file コマンドを使ってsplitするファイルがtextファイルであるか否かを判定する
func detectFileType(file *os.File) bool {
	cmd := exec.Command("file", "-")
	cmd.Stdin = file

	output, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	if !strings.Contains(string(output), "text") {
		return false
	}
	return true
}

// 対応するoptionを-n, -l, -bのみに制限した場合、optionのvalidationについては複数指定されているか否かで判定できる
// この条件を満たさない場合はvalidationについての実装が異なるので拡張する場合は要変更
func validateOptions() bool {
	optionCount := 0
	flag.VisitAll(func(f *flag.Flag) {
		if f.Value.String() != f.DefValue {
			optionCount++
		}
	})

	if optionCount > 1 {
		return false
	}

	return true
}

// プログラムの引数として指定されたoptionを返す
// 事前条件: すでにoptionsは適切なものが残っていることが保証されている
func selectOption(options []option.Command) option.Command {
	defaultOption := option.DefaultLineCount

	var selected option.Command
	selected = option.NewLineCount(defaultOption)
	for _, o := range options {
		if o.IsDefaultValue() {
			continue
		}
		selected = o
	}
	return selected
}
