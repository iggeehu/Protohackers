package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

const (
	HOST = "localhost"
	PORT = "10000"
	TYPE = "tcp"
)

func main() {
	tcpServer, err := net.ResolveTCPAddr(TYPE, HOST+":"+PORT)
	if err != nil {
		println("ResolveTCPAddr failed:", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP(TYPE, nil, tcpServer)
	if err != nil {
		println("Dial failed:", err.Error())
		os.Exit(1)
	}

	defer conn.Close()

	packet := []byte(`{"method":"isPrime","number":1997}`)
	_, err = conn.Write(packet)
	if err != nil {
		println("Write data failed:", err.Error())
		os.Exit(1)
	}
	//new line of code
	conn.CloseWrite()
	fmt.Println("Write to server = ", string(packet))

	received := make([]byte, 0)
	for {
		println("Reading data...")
		temp := make([]byte, 4096)
		_, err = conn.Read(temp)

		if err != nil {
			if err == io.EOF {
				println("End Of File")
			} else {
				println("Read data failed:", err.Error())
			}
			break
		}

		received = append(received, temp...)
	}

	println("Received message:", string(received))

}
