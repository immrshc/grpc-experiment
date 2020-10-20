package main

import (
	"context"
	"fmt"
	"sync"
)

//func main() {
//	var wg sync.WaitGroup
//	wg.Add(1)
//
//	ctx, cancel := context.WithCancel(context.Background())
//
//	go func(ctx context.Context) {
//		select {
//		case <-ctx.Done():
//			fmt.Println("----done----")
//			wg.Done()
//			return
//		}
//	}(ctx)
//
//	cancel()
//	wg.Wait()
//}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	c1, can1 := context.WithCancel(context.Background())

	go func(ctx context.Context) {
		c2, _ := context.WithCancel(ctx)

		go func(ctx context.Context) {
			c3, _ := context.WithCancel(ctx)
			select {
			case <-c3.Done():
				fmt.Println("----done----")
				wg.Done()
				return
			}
		}(c2)
	}(c1)

	can1()
	wg.Wait()
}

//func main() {
//	ch1, ch2 := make(chan struct{}), make(chan struct{})
//
//	go func() {
//		ch1 <- struct{}{}
//	}()
//	go func() {
//		time.Sleep(3 * time.Second)
//		ch2 <- struct{}{}
//	}()
//
//	ch1 = nil
//	select {
//	case <-ch1:
//		fmt.Println("ch1 received");
//	case <-ch2:
//		fmt.Println("ch2 received")
//	}
//}
