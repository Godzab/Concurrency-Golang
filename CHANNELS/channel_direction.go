package main

import "fmt"

func genMsg(ch1 chan <- string){
	ch1 <- "message"
}

func relayMsg(ch1 <- chan string, ch2 chan <- string){
	m := <- ch1
	ch2 <- m + " - With Additional Information"
}

func main(){
	ch1 := make(chan string)
	ch2 := make(chan string)

	go genMsg(ch1)
	go relayMsg(ch1, ch2)

	v := <- ch2
	fmt.Println(v)
}