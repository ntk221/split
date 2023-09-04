package option

import (
	"log"
	"regexp"
	"strconv"
	"strings"
)

type LineCount int

func NewLineCount(i int) LineCount       { return LineCount(i) }
func (LineCount) Type() int              { return LineCountType }
func (l LineCount) IsDefaultValue() bool { return l == DefaultLineCount }
func (l LineCount) ConvertToNum() uint64 { return uint64(l) }

type ChunkCount int

func NewChunkCount(i int) ChunkCount      { return ChunkCount(i) }
func (ChunkCount) Type() int              { return ChunkCountType }
func (c ChunkCount) IsDefaultValue() bool { return c == DefaultChunkCount }
func (c ChunkCount) ConvertToNum() uint64 { return uint64(c) }

type ByteCount int

func NewByteCount(s string) ByteCount    { return parseByteCount(s) }
func (ByteCount) Type() int              { return ByteCountType }
func (b ByteCount) IsDefaultValue() bool { return b == DefaultByteCount }
func (b ByteCount) ConvertToNum() uint64 { return uint64(b) }

func parseByteCount(s string) ByteCount {
	pattern := regexp.MustCompile(`^(\d+)([KkMmGg]?)$`)
	match := pattern.FindStringSubmatch(s)
	if match == nil {
		return DefaultByteCount
	}

	valueStr := match[1]
	unit := strings.ToLower(match[2])

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Fatal(err)
	}

	switch unit {
	case "k":
		value *= 1024
	case "m":
		value *= 1024 * 1024
	case "g":
		value *= 1024 * 1024 * 1024
	}

	return ByteCount(value)
}
