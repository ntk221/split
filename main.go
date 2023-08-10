package main

import (
	"fmt"
	"log"
	"os"
	"bufio"
	"flag"
	"io"
)

var (
	lineCount = flag.Int("l", 1000, "行数を指定してください")
)

func main() {
	flag.Parse()

	fmt.Printf("line_count is %d\n", *lineCount)
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

	partNum := 0
	outputPrefix := "x"

	reader := bufio.NewReader(file)
	for {
		partName := fmt.Sprintf("%s%03d", outputPrefix, partNum)
		partFileName := fmt.Sprintf("%s.txt", partName)
		partFile, err := os.Create(partFileName)
		if err != nil {
			log.Fatal(err)
		}

		for i := 0; i < *lineCount; i++ {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					os.Remove(partFileName)
					fmt.Println("ファイル分割が終了しました")
					return
				}
				fmt.Println("行を読み込めませんでした:", err)
				return
			}
			
			partFile.WriteString(line)
		}
		defer partFile.Close()

		partNum++
	}
}

