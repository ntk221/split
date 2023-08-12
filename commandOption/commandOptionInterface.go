package commandOption

const (
	ChunkCountType = iota
	LineCountType
	ByteCountType

	DefaultChunkCount = 0
	DefaultByteCount  = 0
	DefaultLineCount  = 1000
)

type CommandOption interface {
	OptionType() int
	IsDefaultValue() bool
	ConvertToNum() uint64
}
