package main

import (
	"fmt"
	"io"

	"log"

	"net"
)

func main() {

	PORT := ":10000"
	l, err := net.Listen("tcp4", PORT)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	fmt.Println("Server is listening.")
	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleConnection(c)
	}
}

func handleConnection(conn net.Conn) {
	fmt.Printf("Serving %s\n", conn.RemoteAddr().String())
	packet := make([]byte, 0)
	tmp := make([]byte, 1)
	defer conn.Close()
	for {
		_, err := conn.Read(tmp)
		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err)
			}
			println("END OF FILE")
			break
		}
		packet = append(packet, tmp...)
	}
	
	fmt.Printf("message is %s", string(packet))

}

