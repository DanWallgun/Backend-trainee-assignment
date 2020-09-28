package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
)

var (
	// ErrBadShortURL ...
	ErrBadShortURL = errors.New("Bad Request. Bad short URL")
	// ErrURLNotFound ...
	ErrURLNotFound = errors.New("Bad Request. URL not found")
	// ErrInternalServerError ...
	ErrInternalServerError = errors.New("Internal Server Error")
	// ErrInvalidJSONFormat ...
	ErrInvalidJSONFormat = errors.New("Bad Request. Invalid	JSON format")
	// ErrIncorrectURL ...
	ErrIncorrectURL = errors.New("Bad Request. Incorrect URL")
	// ErrAlreadyUsed ...
	ErrAlreadyUsed = errors.New(
		"Bad request. This short URL is already in use in another collation",
	)
)

// ErrorInfo contains error info for serialization
type ErrorInfo struct {
	StatusCode int    `json:"status"`
	Detail     string `json:"detail"`
}

// ErrorHandler writes error info to Response
func ErrorHandler(w http.ResponseWriter, err error, code int) {
	jsonBytes, _ := json.Marshal(struct {
		ErrorsInfo []*ErrorInfo `json:"errors"`
	}{
		ErrorsInfo: []*ErrorInfo{
			&ErrorInfo{
				StatusCode: code,
				Detail:     err.Error(),
			},
		},
	})
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
}
