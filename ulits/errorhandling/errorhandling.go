package errorhandling

import (
	"fmt"
	"runtime"
)

func HandleError(errorString string, err error) (bool, error) {
	if err != nil {
		_, fn, line, _ := runtime.Caller(1)
		return true, fmt.Errorf("%s\n[error] %s:%d\n%v", errorString, fn, line, err)
	}
	return false, nil
}

func FormatError(errorString string) error {
	_, fn, line, _ := runtime.Caller(1)
	return fmt.Errorf("%s\n[error] %s:%d", errorString, fn, line)
}
