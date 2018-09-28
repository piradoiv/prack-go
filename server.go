package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

const apiEndpoint = "/api/v1/request"

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
	if len(os.Args) != 2 {
		fmt.Println("Usage: " + os.Args[0] + " [host]:[port]")
		fmt.Println("Example: " + os.Args[0] + " localhost:4242")
		return
	}
	url := "http://" + os.Args[1] + apiEndpoint

	for {
		req, err := getNextRequest(url)
		if err != nil {
			if err == errPrackIsDown {
				time.Sleep(5 * time.Second)
			}
			continue
		}

		res := &response{
			Identifier: req.Identifier,
			Code:       200,
			Headers: map[string]string{
				"Content-Type": "text/html",
				"Connection":   "close",
			},
			Body: "Hello, " + req.Identifier + "!",
		}

		http.Post(url, "application/json", buildResponseJSON(res))
	}
}

func getNextRequest(url string) (request, error) {
	res, err := http.Get(url)
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

func buildResponseJSON(res *response) io.Reader {
	encodedBody := base64.StdEncoding.EncodeToString([]byte(res.Body))
	res.Body = encodedBody
	marshalledResponse, _ := json.Marshal(res)
	return strings.NewReader(string(marshalledResponse))
}
