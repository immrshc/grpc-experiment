package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
)

func Sum(nums []int) int {
	var s int
	for _, n := range nums {
		s += n
	}
	return s
}

func SumConcurrently(nums []int, cncrtNum int) int {
	totalNum := len(nums)
	numPerGrtne := totalNum / cncrtNum

	var s int64
	var wg sync.WaitGroup
	for i := 0; i < cncrtNum; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			start := i * numPerGrtne
			end := start + numPerGrtne
			if i == cncrtNum-1 {
				end = totalNum
			}
			w := nums[start:end]
			atomic.AddInt64(&s, int64(Sum(w)))
		}(i)
	}
	wg.Wait()
	return int(s)
}

func request() {
	res, err := http.Get("http://www.google.com/robots.txt")
	if err != nil {
		log.Fatal(err)
	}
	if _, err := ioutil.ReadAll(res.Body); err != nil {
		log.Fatal(err)
	}
	if err = res.Body.Close(); err != nil {
		log.Fatal(err)
	}
}

func Do(num int) {
	for i := 0; i < num; i++ {
		request()
	}
}

func DoConcurrently(cncrtNum int) {
	var wg sync.WaitGroup
	for i := 0; i < cncrtNum; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			request()
		}()
	}
	wg.Wait()
}

func main() {
	//Do(20)
	//DoConcurrently(20)
	//DoConcurrently(50)
	DoConcurrently(100)
	//DoConcurrently(200)
}
