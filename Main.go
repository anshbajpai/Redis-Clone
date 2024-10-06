package main

import (
	"fmt"
	"net"
	"strings"
)

func main() {
	listener, error := net.Listen("tcp", ":6379")
	if error != nil {
		fmt.Println(error)
		return
	}

	connection, error := listener.Accept()
	if error != nil {
		fmt.Println(error)
	}

	defer connection.Close()

	for {
		resp := NewResp(connection)
		value, err := resp.Read()
		if err != nil {
			fmt.Println(err)
			return
		}

		if value.typ != "array" {
			fmt.Println("Invalid request, expected array")
			continue
		}

		if len(value.array) == 0 {
			fmt.Println("Invalid request, expected array length > 0")
			continue
		}

		command := strings.ToUpper(value.array[0].bulk)
		args := value.array[1:]

		writer := NewWriter(connection)

		handler, ok := Handlers[command]
		if !ok {
			fmt.Println("Invalid command: ", command)
			writer.Write(Value{typ: "string", str: ""})
			continue
		}

		result := handler(args)
		writer.Write(result)
	}
}
