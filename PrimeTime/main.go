package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"reflect"
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
	tmp:=make([]byte, 1024)
	defer conn.Close()
	
	for{
		n, err := conn.Read(tmp)
		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err)
			}
			println("END OF FILE")
			break
		}
		packet = append(packet, tmp[:n]...)
	}

	if !json.Valid(packet) {
		conn.Write([]byte("Malformed JSON"))	
		fmt.Println("packet is not JSON")
		conn.Close()
	} else{
		fmt.Println("packet is JSON")
		//decode JSON object in packet
		obj:=make(map[string]interface{})
		err:=json.Unmarshal(packet, &obj)
		methodVal, methodOk:=obj["Method"]
		numberVal, numberOk:=obj["Prime"]

		if err!=nil{
		  fmt.Println(err)
		}


		if !methodOk || !numberOk || methodVal!="isPrime"||
		reflect.TypeOf(numberVal).String()!= "float64" {
			conn.Write([]byte("Malformed JSON"))
			conn.Close()
		}else{
		response:=make(map[string]interface{})
		response["Prime"]=isPrime(int(numberVal.(float64)))
		response["Method"]="isPrime"
		jsonResponse, _:=json.Marshal(response)
		conn.Write(jsonResponse)
		}
		}	
	}

func isPrime(num int) bool {
	for i := 2; i < num; i++ {
		if num%i == 0 {
			return false
		}
	}
	return true}

