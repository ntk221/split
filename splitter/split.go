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
	FileLimit = "zz"
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
	panic("意図しないOptionTypeです")
}

// SplitUsingLineCount lineCount分だけ、fileから読み込み、他のファイルに出力する
// 事前条件: CommandOptionの種類はlineCountでなくてはならない
func (s *Splitter) SplitUsingLineCount() {
	lineCountOption := s.option
	outputPrefix := s.outputPrefix

	if lineCountOption.OptionType() != LineCountType {
		panic("SplitUsingLineCountがLineCount以外のCommandOptionで呼ばれている")
	}

	file := s.file
	partCount := "aa"
	reader := bufio.NewReader(file)
	for {
		if partCount >= FileLimit {
			log.Fatal("too many files")
		}
		partName := fmt.Sprintf("%s%s", outputPrefix, partCount)
		partFileName := fmt.Sprintf("%s.txt", partName)
		partFile, err := os.Create(partFileName)
		if err != nil {
			log.Fatal(err)
		}

		var i uint64
		lineCount := lineCountOption.ConvertToNum()
		for i = 0; i < lineCount; i++ {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					// 一度も読み込めない時はCreateしたファイルを消す
					if i == 0 {
						_ = os.Remove(partFileName)
					}
					// fmt.Println("ファイル分割が終了しました")
					return
				} else {
					fmt.Println("行を読み込めませんでした")
					log.Fatal(err)
				}
			}

			_, _ = partFile.WriteString(line)
		}
		_ = partFile.Close()

		partCount = incrementString(partCount)
	}
}

func (s *Splitter) SplitUsingChunkCount() {
	chunkCountOption := s.option
	outputPrefix := s.outputPrefix
	_ = outputPrefix

	if chunkCountOption.OptionType() != ChunkCountType {
		panic("SplitUsingChunkCountがLineCount以外のCommandOptionで呼ばれている")
	}

	file := s.file
	partCount := "aa"
	reader := bufio.NewReader(file)
	// 全てのfile内容([]byte)を読み込む
	content, err := io.ReadAll(reader)
	if err != nil {
		log.Fatal(err)
	}

	chunkCount := chunkCountOption.ConvertToNum()
	chunkSize := uint64(len(content)) / chunkCount

	var i uint64
	for i = 0; i < chunkCount; i++ {
		if partCount >= FileLimit {
			log.Fatal("too many files")
		}
		partName := fmt.Sprintf("%s%s", outputPrefix, partCount)
		partFileName := fmt.Sprintf("%s.txt", partName)
		partFile, err := os.Create(partFileName)
		if err != nil {
			log.Fatal(err)
		}

		// i番目のchunkを特定する
		start := i * chunkSize
		end := start + chunkSize
		// i が n-1番目の時はendをcontentの終端に揃える(manを参照)
		if i == chunkCount-1 {
			end = uint64(len(content))
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

		partCount = incrementString(partCount)
	}

	return
}

func (s *Splitter) SplitUsingByteCount() {
	byteCountOption := s.option

	if byteCountOption.OptionType() != ByteCountType {
		panic("SplitUsingByteCountがByteCount以外のCommandOptionで呼ばれている")
	}

	file := s.file
	outputPrefix := s.outputPrefix

	partCount := "aa"
	reader := bufio.NewReader(file)

	for {
		if partCount >= FileLimit {
			log.Fatal("too many files")
		}

		partName := fmt.Sprintf("%s%s", outputPrefix, partCount)
		partFileName := fmt.Sprintf("%s.txt", partName)
		partFile, err := os.Create(partFileName)
		if err != nil {
			log.Fatal(err)
		}

		var writtenBytes uint64

		byteCount := byteCountOption.ConvertToNum()
		bufSize := getNiceBuffer(byteCount)
		buf := make([]byte, bufSize)
		for writtenBytes < byteCount {
			readSize := bufSize
			if writtenBytes+readSize > byteCount {
				readSize = byteCount - writtenBytes
			}

			n, err := reader.Read(buf[:readSize])
			if err != nil {
				if err == io.EOF {
					// fmt.Println("ファイル分割が終了しました")
					return
				} else {
					fmt.Println("バイトを読み込めませんでした")
					log.Fatal(err)
				}
			}
			_, _ = partFile.Write(buf[:n])
			writtenBytes += uint64(n)
		}

		_ = partFile.Close()
		partCount = incrementString(partCount)
	}
}

func getNiceBuffer(byteCount uint64) uint64 {
	if byteCount > 1024*1024*1024 {
		return 32 * 1024 * 1024 // 32MB バッファ
	} else if byteCount > 1024*1024 {
		return 4 * 1024 * 1024 // 4MB バッファ
	} else if byteCount > 1024 {
		return 64 * 1024 // 64KB バッファ
	}
	return 4096 // デフォルト 4KB バッファ
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
