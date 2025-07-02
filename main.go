package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	port := "6379"
	fmt.Println("Starting redis server...")
	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	fmt.Printf("Listening on port %s.\n", port)
	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection:", err)
		return
	}

	defer conn.Close()

	for {
		resp := NewRESP(conn)
		value, err := resp.Read()
		if err != nil {
			if err == io.EOF {
				fmt.Println("Connection closed by client.")
				break
			}
			fmt.Println("Error reading from connection:", err.Error())
			os.Exit(1)
		}
		fmt.Println(value)
		conn.Write([]byte("+OK\r\n"))
	}
}
