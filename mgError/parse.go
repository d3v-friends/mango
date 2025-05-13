package mgError

import "fmt"

func Parse(err error) error {
	switch err.Error() {
	case errNotFoundModel:
		return fmt.Errorf(ErrNotFoundModel)
	}
	return err
}
