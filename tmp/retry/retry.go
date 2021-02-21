package main

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type CheckRetry func(*http.Response, error) bool
type Backoff func(attemptNum int) time.Duration

type Client struct {
	RetryMax   int
	HTTPClient *http.Client
	CheckRetry CheckRetry
	Backoff    Backoff
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	buf, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	if err := req.Body.Close(); err != nil {
		return nil, err
	}
	bodyReader := func() io.ReadCloser {
		return ioutil.NopCloser(bytes.NewReader(buf))
	}
	var attemptNum int
	for {
		attemptNum++
		req.Body = bodyReader() // TODO: net/http.Transportのrewindの実装を確認する
		res, err := c.HTTPClient.Do(req)
		shouldRetry := c.CheckRetry(res, err)
		if !shouldRetry {
			return res, err
		}
		if c.RetryMax < attemptNum {
			return nil, errors.New("retry max exceeded")
		}
		wait := c.Backoff(attemptNum)
		time.Sleep(wait)
	}
}

func main() {
	client := Client{
		RetryMax:   3,
		HTTPClient: http.DefaultClient,
		CheckRetry: func(res *http.Response, err error) bool {
			if err == nil {
				return false
			}
			return res.StatusCode >= http.StatusInternalServerError
		},
		Backoff: func(attemptNum int) time.Duration {
			return time.Second * time.Duration(attemptNum)
		},
	}
	req, err := http.NewRequest("GET", "http://www.google.com/robots.txt", nil)
	if err != nil {
		log.Fatal(err)
	}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	if err := res.Body.Close(); err != nil {
		log.Fatal(err)
	}
}
