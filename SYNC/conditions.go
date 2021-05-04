package main

import (
	"fmt"
	"sync"
)

var sharedRsc = make(map[string]interface{})

func main(){
	var wg sync.WaitGroup
	defer wg.Wait()

	mu := sync.Mutex{}
	c := sync.NewCond(&mu)

	wg.Add(1)
	go func() {
		defer wg.Done()

		c.L.Lock()
		for len(sharedRsc) < 1{
			c.Wait()
		}

		fmt.Println(sharedRsc["data"])
		c.L.Unlock()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		c.L.Lock()
		for len(sharedRsc) < 2{
			c.Wait()
		}

		fmt.Println(sharedRsc["data2"])
		c.L.Unlock()
	}()

	c.L.Lock()
	sharedRsc["data"] = "foo brother!"
	sharedRsc["data2"] = "Bar brother!"
	c.Broadcast()
	c.L.Unlock()
}
