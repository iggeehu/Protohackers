package main

import (
	"io"
	"net"
	"os"
)

const (
	HOST = "188.93.149.98"
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

	_, err = conn.Write([]byte("Ground Control To Major Tom"))
	if err != nil {
		println("Write data failed:", err.Error())
		os.Exit(1)
	}
	//new line of code
	conn.CloseWrite()

	received := make([]byte, 4096)
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
