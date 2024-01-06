package typ

type (
	ResultList[T any] struct {
		Page  int64
		Size  int64
		Total int64
		List  []*T
	}
)
