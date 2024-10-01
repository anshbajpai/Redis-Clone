package main

import (
	"fmt"
	"net"
	"os"
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
		if error != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println(value)

		connection.Write([]byte("+PONG OK\r\n"))
	}
}
