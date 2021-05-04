package main

import (
	"context"
	"fmt"
)

func main() {
	generator := func(ctx context.Context) <-chan int {
		ch := make(chan int)
		n := 1
		go func() {
			defer close(ch)
			for {
				select {
				case ch <- n:
				case <-ctx.Done():
					return
				}
				n++
			}
		}()
		return ch
	}

	ctx, cancel := context.WithCancel(context.Background())

	ch := generator(ctx)

	for n := range ch {
		fmt.Println(n)
		if n == 5 {
			cancel()
		}
	}
}
