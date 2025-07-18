package main

import (
	"fmt"
	"io"
	"net"
	"strings"
)

func main() {
	port := "6379"
	fmt.Println("Starting redis server...")
	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}

	aof, err := NewAof("Database.aof")
	if err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println("Aof active.")
	}
	defer aof.Close()

	aof.Read(func(value Value) {
		command := strings.ToUpper(value.array[0].bulk)
		args := value.array[1:]
		handler, ok := Handlers[command]
		if !ok {
			fmt.Println("Invalid command in aof read ", command)
			return
		}
		handler(args)
	})

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
			return
		}

		if value.typ != "array" {
			fmt.Println("Invalid request, expected array!")
			continue
		}

		if len(value.array) == 0 {
			fmt.Println(("Invalid empty request - expected array lenght > 0"))
		}

		fmt.Println("Received value: ", value)
		command := strings.ToUpper(value.array[0].bulk)
		args := value.array[1:]

		writer := NewWriter(conn)

		handler, ok := Handlers[command]
		if !ok {
			fmt.Println("Invalid command: ", command)
			writer.Write(Value{
				typ: "error",
				str: fmt.Sprintf("ERR unknown command '%s'", command),
			})
			continue
		}
		if command == "SET" || command == "HSET" || command == "DEL" || command == "HDEL" {
			aof.Write(value)
		}
		result := handler(args)
		writer.Write(result)
	}
}
