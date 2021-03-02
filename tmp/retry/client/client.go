package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"runtime"
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
	ctx := req.Context()

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
		drainBody(res.Body)
		wait := c.Backoff(attemptNum)

		select {
		case <-ctx.Done():
			if err := req.Body.Close(); err != nil {
				return nil, err
			}
			return nil, ctx.Err()
		case <-time.After(wait):
		}

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

func drainBody(body io.ReadCloser) {
	io.Copy(io.Discard, body)
	body.Close()
}

func PooledTransport() *http.Transport {
	var t *http.Transport
	dt := http.DefaultTransport.(*http.Transport)
	t = dt.Clone()
	t.MaxIdleConnsPerHost = runtime.NumCPU()
	//t.DisableKeepAlives = true
	return t
}

func PooledClient() *http.Client {
	return &http.Client{
		Transport: PooledTransport(),
	}
}

func main() {
	//client := PooledClient()
	client := http.DefaultClient
	for i := 0; i < 30; i++ {
		res, err := client.Get("http://127.0.0.1:8080")
		if err != nil {
			log.Fatal(err)
		}
		buf, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("res: %s\n", buf)
		if err := res.Body.Close(); err != nil {
			log.Fatal(err)
		}
		time.Sleep(1 * time.Second)
	}
	//client.CloseIdleConnections()
	time.Sleep(3 * time.Minute)
}
