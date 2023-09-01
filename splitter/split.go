package splitter

// splitter は splitコマンドの仕様にしたがって、file を分割する処理を担当する
// Splitメソッドの実装中に現れる個々のヘルパーメソッドに関しては本パッケージの下部を参照してください

import (
	"bufio"
	"fmt"
	. "github.com/ntk221/split/commandOption"
	"io"
	"log"
	"os"
)

const (
	FileLimit = "zz"
)

// 書き込み先のfileを抽象化したinterface
type StringWriteCloser interface {
	io.WriteCloser
	io.StringWriter
}

// 書き込み先のfileを生成する処理を抽象化したinteface
type Creator interface {
	Create(name string) (StringWriteCloser, error)
}

type Splitter struct {
	option       CommandOption
	outputPrefix string
	fileCreator  Creator
}

func (s *Splitter) Split(file io.Reader) {
	if file == nil {
		panic("Splitの呼び出し時のfileにnilが入ってきている")
	}

	option := s.option
	if option.OptionType() == LineCountType {
		s.splitUsingLineCount(file)
		return
	}
	if option.OptionType() == ChunkCountType {
		s.splitUsingChunkCount(file)
		return
	}
	if option.OptionType() == ByteCountType {
		s.splitUsingByteCount(file)
		return
	}
	panic("意図しないOptionTypeです")
}

// SplitUsingLineCount lineCount分だけ、fileから読み込み、他のファイルに出力する
// 事前条件: CommandOptionの種類はlineCountでなくてはならない
func (s *Splitter) splitUsingLineCount(file io.Reader) {
	lineCountOption := s.option
	outputPrefix := s.outputPrefix

	if lineCountOption.OptionType() != LineCountType {
		panic("SplitUsingLineCountがLineCount以外のCommandOptionで呼ばれている")
	}

	partCount := "aa"
	reader := bufio.NewReader(file)
	for {
		if partCount >= FileLimit {
			deletePartFile(outputPrefix)
			log.Fatal("too many files")
		}

		// ファイルの作成処理
		partName := fmt.Sprintf("%s%s", outputPrefix, partCount)
		partFileName := fmt.Sprintf("%s", partName)
		partFile, err := s.fileCreator.Create(partFileName)
		if err != nil {
			log.Fatal(err)
		}

		var i uint64
		lineCount := lineCountOption.ConvertToNum()
		// ファイル内容の分割処理
		for i = 0; i < lineCount; i++ {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					_ = partFile.Close()
					return
				} else {
					fmt.Println("行を読み込めませんでした")
					log.Fatal(err)
				}
			}

			// 分割した内容の書き込み処理
			_, err = partFile.WriteString(line)
			if err != nil {
				log.Fatal(err)
			}
		}
		err = partFile.Close()
		if err != nil {
			log.Fatal(err)
		}

		partCount = incrementString(partCount)
	}
}

func (s *Splitter) splitUsingChunkCount(file io.Reader) {
	chunkCountOption := s.option
	outputPrefix := s.outputPrefix
	_ = outputPrefix

	if chunkCountOption.OptionType() != ChunkCountType {
		panic("SplitUsingChunkCountがLineCount以外のCommandOptionで呼ばれている")
	}

	partCount := "aa"
	reader := bufio.NewReader(file)
	content, err := io.ReadAll(reader)
	if err != nil {
		log.Fatal(err)
	}

	chunkCount := chunkCountOption.ConvertToNum()
	chunkSize := uint64(len(content)) / chunkCount

	var i uint64
	for i = 0; i < chunkCount; i++ {
		if partCount >= FileLimit {
			deletePartFile(outputPrefix)
			log.Fatal("too many files")
		}

		partName := fmt.Sprintf("%s%s", outputPrefix, partCount)
		partFileName := fmt.Sprintf("%s", partName)
		partFile, err := s.fileCreator.Create(partFileName)
		if err != nil {
			log.Fatal(err)
		}

		// i番目のchunkを特定する
		start := i * chunkSize
		end := start + chunkSize

		// i が n-1番目の時(最後のchunkの時)はendをcontentの終端に揃える(manを参照)
		if i == chunkCount-1 {
			end = uint64(len(content))
		}

		chunk := content[start:end]

		// i番目のchunkがすでにからの時は終了する
		if !(len(chunk) > 0) {
			_ = os.Remove(partFileName)
			return
		}

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

func (s *Splitter) splitUsingByteCount(file io.Reader) {
	byteCountOption := s.option

	if byteCountOption.OptionType() != ByteCountType {
		panic("SplitUsingByteCountがByteCount以外のCommandOptionで呼ばれている")
	}

	outputPrefix := s.outputPrefix

	partCount := "aa"
	reader := bufio.NewReader(file)

	for {
		if partCount >= FileLimit {
			deletePartFile(outputPrefix)
			log.Fatal("too many files")
		}

		partName := fmt.Sprintf("%s%s", outputPrefix, partCount)
		partFileName := fmt.Sprintf("%s", partName)
		partFile, err := s.fileCreator.Create(partFileName)
		if err != nil {
			log.Fatal(err)
		}

		var writtenBytes uint64

		byteCount := byteCountOption.ConvertToNum()
		bufSize := getNiceBuffer(byteCount)
		buf := make([]byte, bufSize)
		for writtenBytes < byteCount {
			readSize := bufSize
			// 今回bufferのサイズ分読み込んだらbyteCountをオーバーする時
			// optionで指定されたbyteCount - これまで読み込んだバイト数だけ読めばいい
			if writtenBytes+readSize > byteCount {
				readSize = byteCount - writtenBytes
			}

			n, err := reader.Read(buf[:readSize])
			if err != nil {
				if err == io.EOF {
					// fmt.Println("ファイル分割が終了しました")
					// 1バイトも書き込めなかった場合はファイルを消す
					if writtenBytes == 0 {
						_ = os.Remove(partFileName)
						return
					}
					return
				} else {
					fmt.Println("バイトを読み込めませんでした")
					log.Fatal(err)
				}
			}
			_, err = partFile.Write(buf[:n])
			if err != nil {
				log.Fatal(err)
			}
			writtenBytes += uint64(n)
		}

		err = partFile.Close()
		if err != nil {
			log.Fatal(err)
		}
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

func deletePartFile(outputPrefix string) {
	partCount := "aa"
	for partCount < FileLimit {
		partName := fmt.Sprintf("%s%s", outputPrefix, partCount)
		partFileName := fmt.Sprintf("%s", partName)
		err := os.Remove(partFileName)
		if err != nil {
			log.Fatal(err)
		}
		partCount = incrementString(partCount)
	}
}

func New(option CommandOption, outputPrefix string, fileCreator Creator) *Splitter {
	return &Splitter{
		option,
		outputPrefix,
		fileCreator,
	}
}
