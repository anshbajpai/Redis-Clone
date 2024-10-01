package main

import (
	"fmt"
	"io"
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
		buffer := make([]byte, 1024)

		_, error = connection.Read(buffer)
		if error != nil {
			if error == io.EOF {
				break
			}
			fmt.Println("Reading error from client: ", error.Error())
			os.Exit(1)
		}

		connection.Write([]byte("+PONG OK\r\n"))
	}
}
