package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	var data int

	wg.Add(1)
	go func() {
		//This ensures that this method is called at the end of our function.
		defer wg.Done()
		data++
	}()

	wg.Wait()
	//Expect at least 1 to be printed
	fmt.Printf("the value is %v\n", data)
	fmt.Println("Done!")
}
