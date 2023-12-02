package route

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
)

type StatusError struct {
	Code   int
	Msg    string
	ErrLog string
}

func newError(msg string, err error, code int) StatusError {
	pc, filename, line, _ := runtime.Caller(1)

	return StatusError{
		Code:   code,
		Msg:    msg,
		ErrLog: fmt.Sprintf("[error] in %s[%s:%d] %v", runtime.FuncForPC(pc).Name(), filename, line, err),
	}
}

func ErrorBadRequest(msg string, err error) StatusError {
	return newError(msg, err, http.StatusBadRequest)
}

func ErrorNotFound(msg string, err error) StatusError {
	return newError(msg, err, http.StatusNotFound)
}

func ErrorInternal(msg string, err error) StatusError {
	return newError(msg, err, http.StatusInternalServerError)
}

func (sc StatusError) LogError() {
	fmt.Printf("%s\n", sc.ErrLog)
}

func (sc StatusError) JSONError(w http.ResponseWriter) {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(sc.Code)
	errResp := struct {
		Message string
	}{
		Message: sc.Msg,
	}
	json.NewEncoder(w).Encode(errResp)
}
