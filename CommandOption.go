package main

type LineCountOption int

func NewLineCountOption(i int) LineCountOption { return LineCountOption(i) }
func (LineCountOption) OptionType() string     { return "LineCount" }
func (l LineCountOption) IsDefaultValue() bool { return l == DefaultLineCount }
func (l LineCountOption) ConvertToInt() int    { return int(l) }

type ChunkCountOption int

func NewChunkCountOption(i int) ChunkCountOption { return ChunkCountOption(i) }
func (ChunkCountOption) OptionType() string      { return "ChunkCount" }
func (c ChunkCountOption) IsDefaultValue() bool  { return c == DefaultChunkCount }
func (c ChunkCountOption) ConvertToInt() int     { return int(c) }
