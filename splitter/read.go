package splitter

// 渡されるoptionに応じて *bufio.Reader からバイト列を読み込む関数群

import (
	"bufio"
	"fmt"
	"github.com/ntk221/split/option"
	"io"
	"log"
	"os"
)

func readLines(lineCount uint64, reader *bufio.Reader) ([]string, error) {
	var lines []string

	var i uint64
	for i = 0; i < lineCount; i++ {
		line, err := reader.ReadString('\n')
		if err != nil {
			if len(line) > 0 {
				lines = append(lines, line)
			}
			if err == io.EOF {
				return lines, fmt.Errorf("readLines(): %w", err)
			}
			log.Fatalf("readLines()で読み込みに失敗しました: %v", err)
		}
		lines = append(lines, line)
	}

	return lines, nil
}

func readChunk(index uint64, chunkSize uint64, chunkCount uint64, content []byte) ([]byte, bool) {
	// index番目のchunkを特定する
	start := index * chunkSize
	end := start + chunkSize
	// index が n-1番目の時(最後のchunk1の時)はendをcontentの終端に揃える(manを参照)
	if index == chunkCount-1 {
		end = uint64(len(content))
	}
	chunk := content[start:end]
	// i番目のchunkがすでに空の時は終了する
	if !(len(chunk) > 0) {
		return nil, false
	}

	return chunk, true
}

func readBytes(byteCountOption option.ByteCount, reader *bufio.Reader, outputFile *os.File) ([]byte, error) {
	// 読み込んだ合計バイト数
	var readBytes uint64

	byteCount := byteCountOption.ConvertToNum()
	bufSize := getNiceBuffer(byteCount)
	buf := make([]byte, bufSize)

	// optionに0が指定されている場合は全部読みこんで返す
	if byteCountOption == 0 {
		n, err := reader.Read(buf)
		if err != nil {
			if err == io.EOF {
				// fmt.Println("ファイル分割が終了しました")
				// 1バイトも書き込めなかった場合はファイルを消す
				if readBytes == 0 {
					_ = os.Remove(outputFile.Name())
					return nil, fmt.Errorf("readBytes(): %w", err)
				}
			}
			log.Fatalf("readBytes()で読み込みに失敗しました: %v", err)
		}
		return buf[:n], nil
	}

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
					return nil, fmt.Errorf("readBytes(): %w", err)
				}
				// 最後に読み込んだ分は書き込んでおく
				_, err = outputFile.Write(buf[:readBytes+uint64(n)])
				if err != nil {
					return nil, fmt.Errorf("readBytes(): %w", err)
				}
				err = outputFile.Close()
				if err != nil {
					return nil, fmt.Errorf("readBytes(): %w", err)
				}
				return nil, ErrFinishWrite
			}
			// fmt.Println("バイトを読み込めませんでした")
			log.Fatalf("readBytes()で読み込みに失敗しました: %v", err)
		}

		readBytes += uint64(n)
	}

	return buf[:readBytes], nil
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
