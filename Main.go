package main

import (
	"fmt"
	"net"
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

		_ = value

		writer := NewWriter(connection)
		writer.Write(Value{typ: "string", str: "OK"})
	}
}
