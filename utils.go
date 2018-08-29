package main

import (
	"encoding/json"
	"net/http"
	"time"
)

// StringResponse - Json model response
type StringResponse struct {
	Code int    `json:"code"`
	Data string `json:"data"`
}

// ErrorResponse - Json model response
type ErrorResponse struct {
	Code   int      `json:"code"`
	Errors []string `json:"errors"`
}

func panicerr(err error) {
	if err != nil {
		panic(err)
	}
}

func createDataStringResponse(w http.ResponseWriter, code int, data string) {
	response := StringResponse{Code: code, Data: data}
	responseJson, e := json.Marshal(response)
	panicerr(e)

	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(code)
	w.Write(responseJson)
}

func createErrorResponse(w http.ResponseWriter, code int, errors []string) {
	//debug.PrintStack()
	response := ErrorResponse{Code: code, Errors: errors}
	responseJson, e := json.Marshal(response)
	panicerr(e)

	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(code)
	w.Write(responseJson)
}

func stringNilOrEmpty(str *string) bool {
	return str == nil || len(*str) == 0
}

func nowUtc() string {
	return time.Now().UTC().Format("2006-01-02T15:04:05.999Z07:00")
}
