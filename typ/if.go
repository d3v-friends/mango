package typ

type (
	Model interface {
		GetColNm() string
		GetMigrate() []FnMigrate
	}

	PageArgs interface {
		GetSize() int64
		GetPage() int64
	}

	Filter interface {
		Filter() (filter any, err error)
	}

	Sort interface {
		Sort() (filter any, err error)
	}
)
