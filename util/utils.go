package util

import (
	"encoding/json"
	"net/http"
	"time"
)

func Panicerr(err error) {
	if err != nil {
		panic(err)
	}
}

func CreateDataStringResponse(w http.ResponseWriter, code int, data string) {
	response := DataResponse{Code: code, Data: data}
	responseJson, e := json.Marshal(response)
	Panicerr(e)

	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(code)
	w.Write(responseJson)
}

func CreateErrorResponse(w http.ResponseWriter, code int, errors []string) {
	//debug.PrintStack()
	response := ErrorResponse{Code: code, Errors: errors}
	responseJson, e := json.Marshal(response)
	Panicerr(e)

	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(code)
	w.Write(responseJson)
}

func StringNilOrEmpty(str *string) bool {
	return str == nil || len(*str) == 0
}

func NowUtc() string {
	return time.Now().UTC().Format("2006-01-02T15:04:05.999Z07:00")
}
