package splitter

// splitter は splitコマンドの仕様にしたがって、file を分割する処理を担当する
// Splitメソッドの実装中に現れる個々のヘルパーメソッドに関しては本パッケージの下部を参照してください

import (
	"bufio"
	"errors"
	"fmt"
	. "github.com/ntk221/split/commandOption"
	"io"
	"log"
	"os"
)

const (
	FileLimit string = "zz"
)

var (
	outputSuffix string = "aa" // split 処理が生成する部分ファイルのsuffix。incrementされていく
)

// 書き込み先のfileを抽象化したinterface
type StringWriteCloser interface {
	io.WriteCloser
	io.StringWriter
	Name() string
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
	switch option.OptionType() {
	case LineCountType:
		s.splitUsingLineCount(file)
	case ChunkCountType:
		s.splitUsingChunkCount(file)
	case ByteCountType:
		s.splitUsingByteCount(file)
	default:
		panic("意図しないOptionTypeです")
	}
}

// 入力: outputPrefix, partCount
// 出力: partFile StringWriteCloser
func createFile(c Creator, prefix string, suffix string) (StringWriteCloser, error) {
	partName := fmt.Sprintf("%s%s", prefix, suffix)
	partFileName := fmt.Sprintf("%s", partName)
	partFile, err := c.Create(partFileName)
	if err != nil {
		return nil, err
	}

	return partFile, nil
}

// SplitUsingLineCount lineCount分だけ、fileから読み込み、他のファイルに出力する
// 事前条件: CommandOptionの種類はlineCountでなくてはならない
func (s *Splitter) splitUsingLineCount(file io.Reader) {
	lineCountOption := s.option
	outputPrefix := s.outputPrefix

	if lineCountOption.OptionType() != LineCountType {
		panic("SplitUsingLineCountがLineCount以外のCommandOptionで呼ばれている")
	}

	reader := bufio.NewReader(file)
	// ループが終了するのは
	// 1. 生成するファイルが制限を超える時("aa" ~ "zz" に収まらないとき)
	// 2. 読み込みファイルから読む内容がもうない時
	for {
		if outputSuffix >= FileLimit {
			deletePartFile(outputPrefix)
			log.Fatal("too many files")
		}

		// 書き込み先のファイルの作成
		outputFile, err := createFile(s.fileCreator, outputPrefix, outputSuffix)
		if err != nil {
			log.Fatal(err)
		}

		// 読み込み元のファイルから読み込む
		lineCount := lineCountOption.ConvertToNum()
		lines, err := readLines(lineCount, outputFile, reader)
		if err != nil {
			if errors.Is(err, io.EOF) {
				for _, line := range lines {
					_, err = outputFile.WriteString(line)
					if err != nil {
						log.Fatal(err)
					}
				}
				return
			}
			log.Fatal(err)
		}

		// 書き込み先のファイルに書き込む
		for _, line := range lines {
			_, err = outputFile.WriteString(line)
			if err != nil {
				log.Fatal(err)
			}
		}

		// 書き込んだファイルを閉じる
		err = outputFile.Close()
		if err != nil {
			log.Fatal(err)
		}

		outputSuffix = incrementString(outputSuffix)
	}
}

func (s *Splitter) splitUsingChunkCount(file io.Reader) {
	chunkCountOption := s.option
	outputPrefix := s.outputPrefix
	_ = outputPrefix

	if chunkCountOption.OptionType() != ChunkCountType {
		panic("SplitUsingChunkCountがLineCount以外のCommandOptionで呼ばれている")
	}

	reader := bufio.NewReader(file)
	content, err := io.ReadAll(reader)
	if err != nil {
		log.Fatal(err)
	}

	chunkCount := chunkCountOption.ConvertToNum()
	chunkSize := uint64(len(content)) / chunkCount

	var i uint64
	for i = 0; i < chunkCount; i++ {
		if outputSuffix >= FileLimit {
			deletePartFile(outputPrefix)
			log.Fatal("too many files")
		}

		// 書き込み先のファイルを作る
		outputFile, err := createFile(s.fileCreator, outputPrefix, outputSuffix)
		if err != nil {
			log.Fatal(err)
		}

		// 読み込みファイルから読み込む
		chunk, ok := readChunk(i, chunkSize, chunkCount, content)
		if !ok {
			_ = os.Remove(outputFile.Name())
			return
		}

		// 書き込み先のファイルに書き込む
		_, err = outputFile.Write(chunk)
		if err != nil {
			log.Fatal(err)
		}

		// 書き込んだファイルを閉じる
		err = outputFile.Close()
		if err != nil {
			log.Fatal(err)
		}

		outputSuffix = incrementString(outputSuffix)
	}

	return
}

func (s *Splitter) splitUsingByteCount(file io.Reader) {
	byteCountOption := s.option

	if byteCountOption.OptionType() != ByteCountType {
		panic("SplitUsingByteCountがByteCount以外のCommandOptionで呼ばれている")
	}

	outputPrefix := s.outputPrefix

	reader := bufio.NewReader(file)

	for {
		if outputSuffix >= FileLimit {
			deletePartFile(outputPrefix)
			log.Fatal("too many files")
		}

		// 書き込み先のファイルを作る
		outpuFile, err := createFile(s.fileCreator, outputPrefix, outputSuffix)
		if err != nil {
			log.Fatal(err)
		}

		// 読み込んだ合計バイト数
		var readBytes uint64

		byteCount := byteCountOption.ConvertToNum()
		bufSize := getNiceBuffer(byteCount)
		buf := make([]byte, bufSize)
		/*--- 読み込み処理 ---*/
		for readBytes < byteCount {
			readSize := bufSize
			// 今回bufferのサイズ分読み込んだらbyteCountをオーバーする時
			// optionで指定されたbyteCount - これまで読み込んだバイト数だけ読めばいい
			if readBytes+readSize > byteCount {
				readSize = byteCount - readBytes
			}

			n, err := reader.Read(buf[:readSize])
			if err != nil {
				if err == io.EOF {
					// fmt.Println("ファイル分割が終了しました")
					// 1バイトも書き込めなかった場合はファイルを消す
					if readBytes == 0 {
						_ = os.Remove(outpuFile.Name())
						return
					}
					// 最後に読み込んだ分は書き込んでおく
					_, err = outpuFile.Write(buf[:readBytes+uint64(n)])
					if err != nil {
						log.Fatal(err)
					}
					err = outpuFile.Close()
					if err != nil {
						log.Fatal(err)
					}
					return
				} else {
					fmt.Println("バイトを読み込めませんでした")
					log.Fatal(err)
				}
			}

			readBytes += uint64(n)
		}

		// 書き込み処理
		_, err = outpuFile.Write(buf[:readBytes])
		if err != nil {
			log.Fatal(err)
		}

		err = outpuFile.Close()
		if err != nil {
			log.Fatal(err)
		}

		outputSuffix = incrementString(outputSuffix)
	}
}

// ファイルの読み込み処理
// 入力: lineCount, outputFile
// 出力: lines
func readLines(lineCount uint64, outputFile StringWriteCloser, reader *bufio.Reader) ([]string, error) {
	var lines []string

	var i uint64
	for i = 0; i < lineCount; i++ {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return lines, err
			}
			return nil, err
		}
		lines = append(lines, line)
	}

	return lines, nil
}

// 入力： chunkSize uint64, chunkCount uint64, content []byte
// 出力: selectedChunk []byte
func readChunk(i uint64, chunkSize uint64, chunkCount uint64, content []byte) ([]byte, bool) {
	// i番目のchunkを特定する
	start := i * chunkSize
	end := start + chunkSize
	// i が n-1番目の時(最後のchunkの時)はendをcontentの終端に揃える(manを参照)
	if i == chunkCount-1 {
		end = uint64(len(content))
	}
	chunk := content[start:end]
	// i番目のchunkがすでに空の時は終了する
	if !(len(chunk) > 0) {
		return nil, false
	}

	return chunk, true
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
	for outputSuffix < FileLimit {
		partName := fmt.Sprintf("%s%s", outputPrefix, outputSuffix)
		partFileName := fmt.Sprintf("%s", partName)
		err := os.Remove(partFileName)
		if err != nil {
			log.Fatal(err)
		}
		outputSuffix = incrementString(outputSuffix)
	}
}

func New(option CommandOption, outputPrefix string, fileCreator Creator) *Splitter {
	return &Splitter{
		option,
		outputPrefix,
		fileCreator,
	}
}
