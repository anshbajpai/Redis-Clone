package main

import (
	"fmt"
	"net"
)

func main() {
	listener, error := net.Listen("tcp", "6379")
	if error != nil {
		fmt.Println(error)
		return
	}

	connection, error := listener.Accept()
	if error != nil {
		fmt.Println(error)
	}

	defer connection.Close()
}
