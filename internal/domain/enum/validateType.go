package enum

type ValidateType int

const (
	TYPE_INT ValidateType = iota
	TYPE_STRING
	TYPE_DATA
)

func (s ValidateType) String() string {
	return [...]string{"int", "string", "datetime"}[s]
}
