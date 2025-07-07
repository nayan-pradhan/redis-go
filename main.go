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

		fmt.Println("command: ", command)
		fmt.Println("args: ", args)

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
		result := handler(args)
		writer.Write(result)
	}
}
