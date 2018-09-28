package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const apiEndpoint = "http://localhost:4242/api/v1/request"

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
		res, err := http.Get(apiEndpoint)
		if err != nil {
			fmt.Println("Prack server seems to be down, sleeping")
			time.Sleep(5 * time.Second)
			continue
		}

		if res.StatusCode != 200 {
			continue
		}

		contents, _ := ioutil.ReadAll(res.Body)
		res.Body.Close()

		req := &request{}
		json.Unmarshal(contents, &req)

		response := &response{}
		response.Identifier = req.Identifier
		response.Code = 200
		response.Headers = map[string]string{
			"Content-Type": "text/html",
			"Connection":   "close",
		}
		response.Body = base64.StdEncoding.EncodeToString([]byte("Holi " + req.Identifier + "!\n"))

		marshalledResponse, _ := json.Marshal(response)
		http.Post(apiEndpoint, "application/json", strings.NewReader(string(marshalledResponse)))
	}
}
