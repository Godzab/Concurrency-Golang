package main

import (
	"fmt"
	"time"
)

func main(){
	// Direct call
	fun("This is a test.")

	// Goroutine with different variants of function call
	go fun("Goroutine-1")

	// Goroutine with anonymous function
	go func() {
		fun("Goroutine-2")
	}()

	//Goroutine with function value call
	fv := fun
	go fv("Goroutine-3")

	// Wait for go routines to end
	time.Sleep(100 * time.Millisecond)
	fmt.Println("Done!")
}

func fun(str string){
	for i := 0; i < 3; i++{
		fmt.Println(str)
		time.Sleep(1 * time.Millisecond)
	}
}