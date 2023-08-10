package main

import (
	"fmt"
	"log"
	"os"
	"bufio"
	"flag"
	"io"
)

const (
	FILE_LIMIT = 677 // 27^2
	DEFAULT_PREFIX = "x"
)

var (
	lineCountOption = flag.Int("l", 1000, "行数を指定してください")
)

func splitUsingLineCount(lineCount int, outputPrefix string, partNum int, file *os.File) {

	reader := bufio.NewReader(file)
	for {
		if partNum >= FILE_LIMIT {
			log.Fatal("too many files")
		}
		partName := fmt.Sprintf("%s%03d", outputPrefix, partNum)
		partFileName := fmt.Sprintf("%s.txt", partName)
		partFile, err := os.Create(partFileName)
		if err != nil {
			log.Fatal(err)
		}

		for i := 0; i < lineCount; i++ {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					// 一度も読み込めない時はCreateしたファイルを消す
					if i == 0 {
						os.Remove(partFileName)
					}
					fmt.Println("ファイル分割が終了しました")
					return 
				} else {
					fmt.Println("行を読み込めませんでした:", err)
					log.Fatal(err)
				}
			}

			partFile.WriteString(line)
		}
		partFile.Close()

		partNum++
	}
}

func split(args []string, file *os.File) {
	// optionによって処理を分岐する
	flag.Parse()
	lineCount := *lineCountOption
	// chunkCount := *chunkCountOption
	// byteCount := *byteCountOption

	partNum := 0
	outputPrefix := DEFAULT_PREFIX
	// 引数でprefixが指定されている場合はそれを使う
	if len(args) > 1 {
		outputPrefix = args[1]
	}
	
	// chunkCountが設定されてたら -> splitUsingChunkCount
	// byteCountが設定されてたら -> splitUsingByteCount

	splitUsingLineCount(lineCount, outputPrefix, partNum, file)
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

