package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)

	go func(c context.Context) {
		select {
		case <-c.Done():
			fmt.Println(c.Err())
			wg.Done()
			return
		}
	}(ctx)

	wg.Wait()
}
