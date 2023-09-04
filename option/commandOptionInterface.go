package option

const (
	ChunkCountType = iota
	LineCountType
	ByteCountType

	DefaultChunkCount = 0
	DefaultByteCount  = 0
	DefaultLineCount  = 1000
)

type Command interface {
	Type() int
	IsDefaultValue() bool
	ConvertToNum() uint64
}
