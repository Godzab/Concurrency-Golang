package main

import (
	"fmt"
	"sync"
)

func main(){
 	/*for n := range squares(squares(generator(2, 3))){
 		fmt.Println(n)
	}*/
	handleFanOperations()
}

func handleFanOperations(){
	in := generator(2, 3)

	ch1 := squares(in)
	ch2 := squares(in)

	for n := range merge(ch1, ch2){
		fmt.Println(n)
	}
}

func merge(cs ...<-chan int) <-chan int{
	out := make(chan int)
	var wg sync.WaitGroup

	output := func(c <-chan int) {
		for n := range c{
			out <- n
		}
		wg.Done()
	}

	wg.Add(len(cs))
	for _, c := range cs{
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}


func generator(nums ...int)<-chan int{
	out := make(chan int)
	go func() {
		defer close(out)
		for _, num := range nums{
			out <- num
		}
	}()
	return out
}

func squares(in <-chan int)<-chan int{
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in{
			out <- n * n
		}
	}()
	return out
}
