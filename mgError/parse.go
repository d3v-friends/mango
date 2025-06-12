package mgError

import (
	"fmt"
	"github.com/d3v-friends/go-tools/fnError"
)

func Parse(err error) error {
	switch err.Error() {
	case errNotFoundModel:
		return fmt.Errorf(ErrNotFoundModel)
	}
	return err
}

func ChangeError(err error, parser map[error]error) error {
	for i := range parser {
		if parser[i].Error() == err.Error() {
			return fnError.NewF(parser[i].Error())
		}
	}
	return err
}
