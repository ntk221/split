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

type CommandOption interface {
	OptionType() string
	IsDefaultValue() bool
}

type Splitter struct {
	option       CommandOption
	outputPrefix string
	file         *os.File
}

func (s *Splitter) Split(args []string, file *os.File) {
	option := s.option
	_ = option
	panic("TODO: option の 種類を判定する実装をする")
	// option := s.option
	// if option is chunkCount...
	// chunkCountが設定されてたら -> splitUsingChunkCount
	// byteCountが設定されてたら -> splitUsingByteCount

}

func (s *Splitter) SplitUsingLineCount(lineCount int, outputPrefix string, file *os.File) {

	partNum := 0
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

func NewSplitter(option CommandOption, outputPrefix string, file *os.File) *Splitter {
	return &Splitter{
		option,
		outputPrefix,
		file,
	}
}
