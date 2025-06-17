package mgOp

import "fmt"

func Ref(v string) string {
	return fmt.Sprintf("$%s", v)
}

func DoubleRef(v string) string {
	return fmt.Sprintf("$$%s", v)
}
