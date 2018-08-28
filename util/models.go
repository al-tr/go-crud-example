package util

type DataResponse struct {
	Code int    `json:"code"`
	Data string `json:"data"`
}

type ErrorResponse struct {
	Code   int      `json:"code"`
	Errors []string `json:"errors"`
}
