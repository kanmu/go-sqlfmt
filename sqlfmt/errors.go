package sqlfmt

import (
	"fmt"
)

// FormatError is an error that occurred while sqlfmt.Process
type FormatError struct {
	err error
}

func (e *FormatError) Error() string {
	return fmt.Sprint(e.err)
}
