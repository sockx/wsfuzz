package model

import "net/http"

type RequestData struct {
	Host    string
	Path    string
	Query   string
	Headers http.Header
	Body    []byte
}
