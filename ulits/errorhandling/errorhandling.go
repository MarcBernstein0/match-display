package errorhandling

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
)

func HandleError(errorString string, err error) error {
	if err != nil {
		_, fn, line, _ := runtime.Caller(1)
		return fmt.Errorf("%s\n[error] %s:%d\n%v", errorString, fn, line, err)
	}
	return nil
}

func FormatError(errorString string) error {
	_, fn, line, _ := runtime.Caller(1)
	return fmt.Errorf("%s\n[error] %s:%d", errorString, fn, line)
}

func ErrorResponse(w http.ResponseWriter, message string, httpStatusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	resp := make(map[string]string)
	resp["message"] = message
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}
