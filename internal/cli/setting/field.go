package setting

type Kind uint8

const (
	KindString Kind = iota
	KindDir
	KindFile
	KindInt
	KindBool
)

type Field struct {
	Key         string
	Title       string
	Description string
	Kind        Kind
	Min         int
	Max         int
}
