package splitter

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

	. "github.com/ntk221/split/commandOption"
)

const (
	FileLimit = `zz` // 27^2
)

type Splitter struct {
	option       CommandOption
	outputPrefix string
	file         *os.File
}

func (s *Splitter) Split() {

	option := s.option
	_ = option
	if option.OptionType() == LineCountType {
		s.SplitUsingLineCount()
		return
	}
	if option.OptionType() == ChunkCountType {
		s.SplitUsingChunkCount()
		return
	}
	if option.OptionType() == ByteCountType {
		s.SplitUsingByteCount()
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

	if lineCount.OptionType() != LineCountType {
		panic("SplitUsingLineCountがLineCount以外のCommandOptionで呼ばれている")
	}

	file := s.file
	partNum := `aa`
	reader := bufio.NewReader(file)
	for {
		if partNum >= FileLimit {
			log.Fatal("too many files")
		}
		partName := fmt.Sprintf("%s%s", outputPrefix, partNum)
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

		partNum = incrementString(partNum)
	}
}

func (s *Splitter) SplitUsingChunkCount() {
	chunkCount := s.option
	outputPrefix := s.outputPrefix
	_ = outputPrefix

	if chunkCount.OptionType() != ChunkCountType {
		panic("SplitUsingChunkCountがLineCount以外のCommandOptionで呼ばれている")
	}

	file := s.file
	partNum := `aa`
	reader := bufio.NewReader(file)
	// 全てのfile内容([]byte)を読み込む
	content, err := io.ReadAll(reader)
	if err != nil {
		log.Fatal(err)
	}

	chunkSize := len(content) / chunkCount.ConvertToInt()

	for i := 0; i < chunkCount.ConvertToInt(); i++ {
		if partNum >= FileLimit {
			log.Fatal("too many files")
		}
		partName := fmt.Sprintf("%s%s", outputPrefix, partNum)
		partFileName := fmt.Sprintf("%s.txt", partName)
		partFile, err := os.Create(partFileName)
		if err != nil {
			log.Fatal(err)
		}

		// i番目のchunkを特定する
		start := i * chunkSize
		end := start + chunkSize
		// i が n-1番目の時はendをcontentの終端に揃える(manを参照)
		if i == chunkCount.ConvertToInt()-1 {
			end = len(content)
		}

		chunk := content[start:end]

		_, err = partFile.Write(chunk)
		if err != nil {
			log.Fatal(err)
		}

		err = partFile.Close()
		if err != nil {
			log.Fatal(err)
		}

		partNum = incrementString(partNum)
	}

	return
}

func (s *Splitter) SplitUsingByteCount() {
	panic("TODO: SplitUsingByteCountを実装する")
}

// 文字列用のincrement関数
// ex: incrementString("a") == "b"
// ex: incrementString("az") == "ba"
func incrementString(s string) string {
	runes := []rune(s)
	for i := len(runes) - 1; i >= 0; i-- {
		if runes[i] < 'z' {
			runes[i]++
			return string(runes)
		}
		runes[i] = 'a'
	}

	return "a" + string(runes)
}

func NewSplitter(option CommandOption, outputPrefix string, file *os.File) *Splitter {
	return &Splitter{
		option,
		outputPrefix,
		file,
	}
}
