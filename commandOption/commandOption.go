package commandOption

import (
	"log"
	"regexp"
	"strconv"
	"strings"
)

type LineCountOption int

func NewLineCountOption(i int) LineCountOption { return LineCountOption(i) }
func (LineCountOption) OptionType() int        { return LineCountType }
func (l LineCountOption) IsDefaultValue() bool { return l == DefaultLineCount }
func (l LineCountOption) ConvertToNum() uint64 { return uint64(l) }

type ChunkCountOption int

func NewChunkCountOption(i int) ChunkCountOption { return ChunkCountOption(i) }
func (ChunkCountOption) OptionType() int         { return ChunkCountType }
func (c ChunkCountOption) IsDefaultValue() bool  { return c == DefaultChunkCount }
func (c ChunkCountOption) ConvertToNum() uint64  { return uint64(c) }

type ByteCountOption int

func NewByteCountOption(s string) ByteCountOption { return parseByteCount(s) }
func (ByteCountOption) OptionType() int           { return ByteCountType }
func (b ByteCountOption) IsDefaultValue() bool    { return b == DefaultByteCount }
func (b ByteCountOption) ConvertToNum() uint64    { return uint64(b) }

func parseByteCount(s string) ByteCountOption {
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

	return ByteCountOption(value)
}
