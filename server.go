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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: " + os.Args[0] + " [host]:[port]")
		fmt.Println("Example: " + os.Args[0] + " localhost:4242")
		return
	}
	url := fmt.Sprintf("http://%s%s", os.Args[1], apiEndpoint)

	for {
		loop(url)
	}
}

func loop(url string) {
	req, err := getNextRequest(url)
	if err != nil {
		if err == errPrackIsDown {
			time.Sleep(5 * time.Second)
		}
		return
	}

	res := BuildDefaultResponse(req)
	res.Body = "Hello, " + req.Identifier + "!"
	http.Post(url, "application/json", buildResponseJSON(res))
}

func getNextRequest(url string) (Request, error) {
	req := &Request{}
	res, err := http.Get(url)
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

func buildResponseJSON(res Response) io.Reader {
	encodedBody := base64.StdEncoding.EncodeToString([]byte(res.Body))
	res.Body = encodedBody
	marshalledResponse, _ := json.Marshal(res)
	return strings.NewReader(string(marshalledResponse))
}
