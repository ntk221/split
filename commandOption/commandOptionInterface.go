package commandOption

const (
	ChunkCountType = iota
	LineCountType
	ByteCountType
)

type CommandOption interface {
	OptionType() int
	IsDefaultValue() bool
	ConvertToInt() int
}
