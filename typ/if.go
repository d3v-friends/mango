package typ

type (
	Model interface {
		GetColNm() string
		GetMigrate() []FnMigrate
	}

	Pager interface {
		GetSize() int64
		GetPage() int64
	}

	Filter interface {
		GetFilter() (filter any, err error)
		GetColNm() string
	}

	Sorter interface {
		GetSort() (filter any, err error)
	}
)
