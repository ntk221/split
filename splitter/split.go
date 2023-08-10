package splitter

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

const (
	FileLimit = 677 // 27^2
)

func SplitUsingLineCount(lineCount int, outputPrefix string, partNum int, file *os.File) {

	reader := bufio.NewReader(file)
	for {
		if partNum >= FileLimit {
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
						_ = os.Remove(partFileName)
					}
					fmt.Println("ファイル分割が終了しました")
					return
				} else {
					fmt.Println("行を読み込めませんでした:", err)
					log.Fatal(err)
				}
			}

			_, _ = partFile.WriteString(line)
		}
		_ = partFile.Close()

		partNum++
	}
}
