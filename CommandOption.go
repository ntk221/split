package main

import (
	. "github.com/ntk221/split/commandOption"
)

type LineCountOption int

func NewLineCountOption(i int) LineCountOption { return LineCountOption(i) }
func (LineCountOption) OptionType() int        { return LineCountType }
func (l LineCountOption) IsDefaultValue() bool { return l == DefaultLineCount }
func (l LineCountOption) ConvertToInt() int    { return int(l) }

type ChunkCountOption int

func NewChunkCountOption(i int) ChunkCountOption { return ChunkCountOption(i) }
func (ChunkCountOption) OptionType() int         { return ChunkCountType }
func (c ChunkCountOption) IsDefaultValue() bool  { return c == DefaultChunkCount }
func (c ChunkCountOption) ConvertToInt() int     { return int(c) }
