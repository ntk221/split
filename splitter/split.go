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
	ConvertToInt() int
}

type Splitter struct {
	option       CommandOption
	outputPrefix string
	file         *os.File
}

func (s *Splitter) Split() {
	const LineCount = "LineCount"
	const ChunkCount = "ChunkCount"
	const ByteCount = "ByteCount"

	option := s.option
	_ = option
	if option.OptionType() == LineCount {
		s.SplitUsingLineCount()
		return
	}
	if option.OptionType() == ChunkCount {
		s.SplitUsingChunkCount()
		return
	}
	panic("TODO: option の 種類を判定する実装をする")
	// option := s.option
	// if option is chunkCount...
	// chunkCountが設定されてたら -> splitUsingChunkCount
	// byteCountが設定されてたら -> splitUsingByteCount

}

// SplitUsingLineCount lineCount分だけ、fileから読み込み、他のファイルに出力する
// 事前条件: CommandOptionの種類はlineCountでなくてはならない
func (s *Splitter) SplitUsingLineCount() {
	lineCount := s.option
	outputPrefix := s.outputPrefix

	if lineCount.OptionType() != "LineCount" {
		panic("SplitUsingLineCountがLineCount以外のCommandOptionで呼ばれている")
	}
	file := s.file

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

		for i := 0; i < lineCount.ConvertToInt(); i++ {
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
					fmt.Println("行を読み込めませんでした")
					log.Fatal(err)
				}
			}

			_, _ = partFile.WriteString(line)
		}
		_ = partFile.Close()

		partNum++
	}
}

func (s *Splitter) SplitUsingChunkCount() {
	panic("TODO: SplitUsingChunkCountを実装する")
	return
}

func NewSplitter(option CommandOption, outputPrefix string, file *os.File) *Splitter {
	return &Splitter{
		option,
		outputPrefix,
		file,
	}
}
