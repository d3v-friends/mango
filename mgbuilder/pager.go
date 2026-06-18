package mgbuilder

type PagerArgs interface {
	GetPage() int64
	GetSize() int64
	GetSkip() int64
	GetLimit() int64
}
