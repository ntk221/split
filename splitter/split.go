// Package splitter は splitコマンドの仕様に従って、file を分割する処理を担当する
// split メソッドは主に以下の処理を行う
//
//  1. optionに従って以下の処理を繰り返す
//  2. 書き込み用のファイルを生成する
//  3. 引数として受け取った読み込み用ファイルから読み込む
//  4. 2で用意したファイルに書き込む
//  5. 4で書き込んだファイルをクローズする
//  6. 1に戻る
//
// CLI　構造体は Splitter 構造体をラップしていて、入力元ファイル(io.Reader)と出力ファイル名を切り替えられる
// これによって、テスト時に柔軟性を持たせることが可能
package splitter

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/ntk221/split/option"
	"io"
	"log"
	"os"
)

const (
	FileLimit string = "zz"
)

var (
	ErrFinishWrite = errors.New("ファイルの書き込みが終了しました")
	ErrTooManyFile = errors.New("ファイルが生成できる上限を超えました")
	ErrZeroChunk   = errors.New("chunkが分割可能な上限を超えています")
)

// CLI はSplitter構造体のラッパー
// Input, Output をCLIで制御することができる
type CLI struct {
	Input     io.Reader
	OutputDir string
	Splitter  *Splitter
}

// Run は Splitter の split メソッドを呼び出す
// split　がエラー情報を返すのでそれをそのまま呼び出しもとに返す
func (cli *CLI) Run(option option.Command) error {
	input := cli.Input
	outputDir := cli.OutputDir
	err := cli.Splitter.split(input, outputDir, option)
	return err
}

type Splitter struct {
	outputPrefix string
}

func (s *Splitter) split(input io.Reader, outputDir string, opt option.Command) error {
	if input == nil {
		panic("splitの呼び出し時のfileにnilが入ってきている")
	}

	var err error
	switch opt.(type) {
	case option.LineCount:
		err = s.splitUsingLineCount(input, outputDir, opt)
	case option.ChunkCount:
		err = s.splitUsingChunkCount(input, outputDir, opt)
	case option.ByteCount:
		err = s.splitUsingByteCount(input, outputDir, opt)
	default:
		panic("意図しないOptionTypeです")
	}

	return err
}

// SplitUsingLineCount はlineCount分だけ、fileから読み込み、他のファイルに出力する
// 事前条件: CommandOptionの種類はlineCountでなくてはならない
func (s *Splitter) splitUsingLineCount(file io.Reader, outputDir string, lineCount option.Command) error {
	outputSuffix := "aa"
	outputPrefix := s.outputPrefix

	if _, ok := lineCount.(option.LineCount); !ok {
		panic("SplitUsingLineCountがLineCount以外のCommandOptionで呼ばれている")
	}

	reader := bufio.NewReader(file)
	for {
		if outputSuffix >= FileLimit {
			deletePartFile(outputPrefix)
			return ErrTooManyFile
		}

		// 書き込み先のファイルを生成
		outputFile, err := os.OpenFile(outputDir+"/"+outputPrefix+outputSuffix, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("splitUsingLineCount(): %w", err)
		}

		// 読み込み元のファイルから読み込む
		lineCount := lineCount.ConvertToNum()
		lines, err := readLines(lineCount, reader)
		if err != nil {
			if errors.Is(err, io.EOF) {
				if len(lines) == 0 {
					_ = os.Remove(outputFile.Name())
					return nil
				}
				for _, line := range lines {
					_, err = outputFile.WriteString(line)
					if err != nil {
						log.Fatalf("splitUsingLineCount()でファイル書き込みに失敗しました: %v", err)
					}
				}
				return nil
			}
			log.Fatal(err)
		}

		// 書き込み先のファイルに書き込む
		for _, line := range lines {
			_, err = outputFile.WriteString(line)
			if err != nil {
				return fmt.Errorf("splitUsingLineCount(): %w", err)
			}
		}

		// 書き込んだファイルを閉じる
		err = outputFile.Close()
		if err != nil {
			return fmt.Errorf("splitUsingLineCount(): %w", err)
		}

		outputSuffix = incrementString(outputSuffix)
	}
}

func (s *Splitter) splitUsingChunkCount(file io.Reader, outputDir string, chunkCountOption option.Command) error {
	outputSuffix := "aa"
	outputPrefix := s.outputPrefix
	_ = outputPrefix

	if _, ok := chunkCountOption.(option.ChunkCount); !ok {
		panic("SplitUsingChunkCountがLineCount以外のCommandOptionで呼ばれている")
	}

	reader := bufio.NewReader(file)
	content, err := io.ReadAll(reader)
	if err != nil {
		log.Fatalf("splitUsingChunkCount()でファイル読み込みに失敗しました: %v", err)
	}

	chunkCount := chunkCountOption.ConvertToNum()
	chunkSize := uint64(len(content)) / chunkCount

	if chunkSize == 0 {
		return ErrZeroChunk
	}

	var i uint64
	for i = 0; i < chunkCount; i++ {
		if outputSuffix >= FileLimit {
			deletePartFile(outputPrefix)
			return ErrTooManyFile
		}

		// ファイルの作成またはオープン（存在しなければ新規作成、存在すれば上書き）
		outputFile, err := os.OpenFile(outputDir+"/"+outputPrefix+outputSuffix, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("splitUsingChunkCount()でファイルオープンに失敗しました: %v", err)
		}

		// 読み込みファイルから読み込む
		// iは分割したchunkに割り振ったindex
		chunk, ok := readChunk(i, chunkSize, chunkCount, content)
		if !ok {
			_ = os.Remove(outputFile.Name())
			return fmt.Errorf("splitUsingChunkCount(): %w", err)
		}

		// 書き込み先のファイルに書き込む
		_, err = outputFile.Write(chunk)
		if err != nil {
			log.Fatalf("splitUsingChunkCount()でファイル書き込みに失敗しました: %v", err)
		}

		// 書き込んだファイルを閉じる
		err = outputFile.Close()
		if err != nil {
			log.Fatalf("splitUsingChunkCount()でファイルのクローズに失敗しました: %v", err)
		}

		outputSuffix = incrementString(outputSuffix)
	}

	return nil
}

func (s *Splitter) splitUsingByteCount(file io.Reader, outputDir string, byteCountOption option.Command) error {
	outputSuffix := "aa"

	var byteCount option.ByteCount
	var ok bool
	if byteCount, ok = byteCountOption.(option.ByteCount); !ok {
		panic("SplitUsingByteCountがByteCount以外のCommandOptionで呼ばれている")
	}

	outputPrefix := s.outputPrefix

	reader := bufio.NewReader(file)

	for {
		if outputSuffix >= FileLimit {
			deletePartFile(outputPrefix)
			return ErrTooManyFile
		}

		// ファイルの作成またはオープン（存在しなければ新規作成、存在すれば上書き）
		outputFile, err := os.OpenFile(outputDir+"/"+outputPrefix+outputSuffix, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("splitUsingByteCount()でファイルのオープンに失敗しました: %v", err)
		}

		buf, err := readBytes(byteCount, reader, outputFile)
		if err != nil {
			if errors.Is(err, io.EOF) || errors.Is(err, ErrFinishWrite) {
				return fmt.Errorf("splitUsingByteCount(): %w", err)
			}
			return fmt.Errorf("splitUsingByteCount(): %w", err)
		}

		// 書き込み処理
		_, err = outputFile.Write(buf)
		if err != nil {
			log.Fatalf("splitUsingByteCount()でファイルの書き込み処理に失敗しました: %v", err)
		}

		err = outputFile.Close()
		if err != nil {
			log.Fatalf("splitUsingByteCount()でファイルのクローズに失敗しました: %v", err)
		}

		outputSuffix = incrementString(outputSuffix)
	}
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
	outputSuffix := "aa"
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

func New(outputPrefix string) *Splitter {
	return &Splitter{
		outputPrefix,
	}
}
