package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const apiEndpoint = "http://localhost:4242/api/v1/request"

var errPrackIsDown = errors.New("Prack server seems to be down")
var errNoRequestsPending = errors.New("There are no requests pending")

type request struct {
	Identifier  string            `json:"identifier"`
	Environment map[string]string `json:"environment"`
}

type response struct {
	Identifier string            `json:"identifier"`
	Code       int               `json:"code"`
	Headers    map[string]string `json:"headers"`
	Body       string            `json:"body"`
}

func main() {
	for {
		req, err := getNextRequest()
		if err != nil {
			if err == errPrackIsDown {
				time.Sleep(5 * time.Second)
			}
			continue
		}

		response := buildResponse(
			req.Identifier,
			200,
			map[string]string{
				"Content-Type": "text/html",
				"Connection":   "close",
			},
			"Hello, "+req.Identifier+"!",
		)

		http.Post(apiEndpoint, "application/json", strings.NewReader(response))
	}
}

func getNextRequest() (request, error) {
	res, err := http.Get(apiEndpoint)
	req := &request{}

	if err != nil {
		return *req, errPrackIsDown
	}

	if res.StatusCode != 200 {
		return *req, errNoRequestsPending
	}

	contents, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()
	json.Unmarshal(contents, &req)
	return *req, nil
}

func buildResponse(identifier string, code int, headers map[string]string, body string) string {
	res := &response{}
	res.Identifier = identifier
	res.Code = code
	res.Headers = headers
	res.Body = base64.StdEncoding.EncodeToString([]byte(body))

	marshalledResponse, _ := json.Marshal(res)
	return string(marshalledResponse)
}
