package main

import (
	"io"
	"log"
	"net"
	"time"
)

func main(){
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil{
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil{
			continue
		}
		go handleConn(conn)
	}
}

func handleConn(c net.Conn){
	defer c.Close()
	for {
		_, err := io.WriteString(c, "Response from server\n")
		if err != nil{
			return
		}
		time.Sleep(time.Second)
	}
}
