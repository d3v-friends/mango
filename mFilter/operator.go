package mFilter

type Operator string

func (x Operator) String() string {
	return string(x)
}

const (
	OperatorIn Operator = "$in"
)
