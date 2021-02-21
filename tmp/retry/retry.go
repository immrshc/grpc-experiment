package main

import (
	"errors"
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
	var attemptNum int
	for {
		attemptNum++
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
