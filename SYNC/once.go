package main

import (
	"fmt"
	"sync"
)

func main(){
	var wg sync.WaitGroup
	var once sync.Once
	defer wg.Wait()

	load := func() {
		fmt.Println("Initialization function run only once...")
	}

	for i:=0; i<10; i++{
		wg.Add(1)
		go func() {
			defer wg.Done()
			once.Do(load) //Magic here
		}()
	}
}
