package mgOp

import "fmt"

type Field string

func (x Field) Ref() string {
	return fmt.Sprintf("$%s", x)
}

func (x Field) SystemRef() string {
	return fmt.Sprintf("$$%s", x)
}

func (x Field) String() string {
	return string(x)
}
