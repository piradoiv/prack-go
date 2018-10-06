package main

import (
	"net/http"
)

// Request is what we get from Prack Server
type Request struct {
	Identifier  string            `json:"identifier"`
	Environment map[string]string `json:"environment"`
}

// Response is what we send to Prack Server
type Response struct {
	Identifier string            `json:"identifier"`
	Code       int               `json:"code"`
	Headers    map[string]string `json:"headers"`
	Body       string            `json:"body"`
}

// BuildDefaultResponse builds a default Response
// with a 200 Status Code, text/html Content Type and
// Close Connection
func BuildDefaultResponse(req Request) Response {
	return Response{
		Identifier: req.Identifier,
		Code:       http.StatusOK,
		Headers: map[string]string{
			"Content-Type": "text/html",
			"Connection":   "close",
		},
	}
}
