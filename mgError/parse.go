package mgError

import (
	"github.com/d3v-friends/go-tools/fnError"
)

func ChangeError(err error, parser map[error]error) error {
	for i := range parser {
		if parser[i].Error() == err.Error() {
			return fnError.NewF(parser[i].Error())
		}
	}
	return fnError.New(err.Error())
}
