package main

import (
	"fmt"
	"net"
)

func main() {
	fmt.Println("Starting redis server...")
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection:", err)
		return
	}
	conn.Write([]byte("Ping to redis server!\r\n"))
	defer conn.Close()
}
