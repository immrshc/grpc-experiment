package main

import (
	"bytes"
	"errors"
	"io"
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

func NewClient() *Client {
	return &Client{
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
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	var attemptNum int
	for {
		attemptNum++
		req, err := rewindBody(req)
		if err != nil {
			return nil, err
		}
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

func rewindBody(req *http.Request) (*http.Request, error) {
	if req.Body == nil || req.Body == http.NoBody {
		return req, nil
	}
	if req.GetBody == nil {
		buf, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		req.GetBody = func() (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewReader(buf)), nil
		}
	}
	if err := req.Body.Close(); err != nil {
		return nil, err
	}
	body, err := req.GetBody()
	if err != nil {
		return nil, err
	}
	newReq := *req
	newReq.Body = body
	return &newReq, nil
}

func main() {
	client := NewClient()
	req, err := http.NewRequest("POST", "http://www.google.com/robots.txt", bytes.NewBuffer([]byte("\"{\"x\":1,\"y\":2}\"")))
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
