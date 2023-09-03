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

	ErrFinishWrite = errors.New("ファイルの書き込みが終了しました")
)

type CLI struct {
	Input     io.Reader
	OutputDir string
	Splitter  *Splitter
}

func (cli *CLI) Run() {
	input := cli.Input
	outputDir := cli.OutputDir
	cli.Splitter.Split(input, outputDir)
}

// StringWriteCloser は書き込み先のfileを抽象化したinterface
type StringWriteCloser interface {
	io.WriteCloser
	io.StringWriter
	Name() string
}

type Splitter struct {
	option       CommandOption
	outputPrefix string
}

func (s *Splitter) Split(file io.Reader, outputDir string) {
	if file == nil {
		panic("Splitの呼び出し時のfileにnilが入ってきている")
	}

	option := s.option

	switch option.OptionType() {
	case LineCountType:
		s.splitUsingLineCount(file, outputDir)
	case ChunkCountType:
		s.splitUsingChunkCount(file, outputDir)
	case ByteCountType:
		s.splitUsingByteCount(file, outputDir)
	default:
		panic("意図しないOptionTypeです")
	}
}

// SplitUsingLineCount はlineCount分だけ、fileから読み込み、他のファイルに出力する
// 事前条件: CommandOptionの種類はlineCountでなくてはならない
func (s *Splitter) splitUsingLineCount(file io.Reader, outputDir string) {

	lineCountOption := s.option
	outputPrefix := s.outputPrefix

	if _, ok := lineCountOption.(LineCountOption); !ok {
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

		// ファイルの作成またはオープン（存在しなければ新規作成、存在すれば上書き）
		outputFile, err := os.OpenFile(outputDir+"/"+outputPrefix+outputSuffix, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
			return
		}

		// 読み込み元のファイルから読み込む
		lineCount := lineCountOption.ConvertToNum()
		lines, err := readLines(lineCount, reader)
		if err != nil {
			if errors.Is(err, io.EOF) {
				if len(lines) == 0 {
					_ = os.Remove(outputFile.Name())
					return
				}
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

func (s *Splitter) splitUsingChunkCount(file io.Reader, outputDir string) {
	chunkCountOption := s.option
	outputPrefix := s.outputPrefix
	_ = outputPrefix

	if _, ok := chunkCountOption.(ChunkCountOption); !ok {
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

		// ファイルの作成またはオープン（存在しなければ新規作成、存在すれば上書き）
		outputFile, err := os.OpenFile(outputDir+"/"+outputPrefix+outputSuffix, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
			return
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

func (s *Splitter) splitUsingByteCount(file io.Reader, outputDir string) {
	var byteCountOption ByteCountOption

	var ok bool
	if byteCountOption, ok = s.option.(ByteCountOption); !ok {
		panic("SplitUsingByteCountがByteCount以外のCommandOptionで呼ばれている")
	}

	outputPrefix := s.outputPrefix

	reader := bufio.NewReader(file)

	for {
		if outputSuffix >= FileLimit {
			deletePartFile(outputPrefix)
			log.Fatal("too many files")
		}

		// ファイルの作成またはオープン（存在しなければ新規作成、存在すれば上書き）
		outputFile, err := os.OpenFile(outputDir+"/"+outputPrefix+outputSuffix, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
			return
		}

		buf, err := readBytes(byteCountOption, reader, outputFile)
		if err != nil {
			if errors.Is(err, io.EOF) || errors.Is(err, ErrFinishWrite) {
				return
			}
			log.Fatal(err)
		}

		// 書き込み処理
		_, err = outputFile.Write(buf[len(buf):])
		if err != nil {
			log.Fatal(err)
		}

		err = outputFile.Close()
		if err != nil {
			log.Fatal(err)
		}

		outputSuffix = incrementString(outputSuffix)
	}
}

func readLines(lineCount uint64, reader *bufio.Reader) ([]string, error) {
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

func readBytes(byteCountOption ByteCountOption, reader *bufio.Reader, outputFile StringWriteCloser) ([]byte, error) {
	// 読み込んだ合計バイト数
	var readBytes uint64

	byteCount := byteCountOption.ConvertToNum()
	bufSize := getNiceBuffer(byteCount)
	buf := make([]byte, bufSize)

	for readBytes < byteCount {
		readSize := bufSize
		// 今回bufferのサイズ分読み込んだらbyteCountをオーバーする時
		// (optionで指定されたbyteCount - これまで読み込んだバイト数) の分だけ読めばいい
		if readBytes+readSize > byteCount {
			readSize = byteCount - readBytes
		}

		n, err := reader.Read(buf[:readSize])
		if err != nil {
			if err == io.EOF {
				// fmt.Println("ファイル分割が終了しました")
				// 1バイトも書き込めなかった場合はファイルを消す
				if readBytes == 0 {
					_ = os.Remove(outputFile.Name())
					return nil, err
				}
				// 最後に読み込んだ分は書き込んでおく
				_, err = outputFile.Write(buf[:readBytes+uint64(n)])
				if err != nil {
					log.Fatal(err)
				}
				err = outputFile.Close()
				if err != nil {
					log.Fatal(err)
				}
				return nil, ErrFinishWrite
			} else {
				fmt.Println("バイトを読み込めませんでした")
				log.Fatal(err)
			}
		}

		readBytes += uint64(n)
	}

	return buf, nil
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

func New(option CommandOption, outputPrefix string) *Splitter {
	return &Splitter{
		option,
		outputPrefix,
	}
}
