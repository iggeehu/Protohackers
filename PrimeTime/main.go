package main

import (
	"bufio"
	"encoding/json"
	"fmt"
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
	defer conn.Close()
	scn:=bufio.NewScanner(conn)
	for{
	
		packet:=make([]byte, 0)
		for scn.Scan(){
				packet=append(packet, scn.Bytes()...)
			}
		
		if(len(packet)!=0){
			connstat:=respondToPacket(packet, conn)
			if(connstat==1){
				fmt.Println("conn closed")
			}
		}
}}



func isPrime(num int) bool {
	for i := 2; i < num; i++ {
		if num%i == 0 {
			return false
		}
	}
	return true}

//0 means conn is not closed at the end, 1 means conn is closed
func respondToPacket(packet []byte, conn net.Conn) int {
	if !json.Valid(packet) {
		conn.Write([]byte("Malformed JSON"))	
		fmt.Println("packet is not JSON")
		conn.Close()
		return 1
	} else{
		fmt.Println("packet is JSON")
		//decode JSON object in packet
		obj:=make(map[string]interface{})
		err:=json.Unmarshal(packet, &obj)
		if(err!=nil){
			fmt.Println("unmarshal error: ", err)
		}
		methodVal, methodOk:=obj["method"]
		numberVal, numberOk:=obj["number"]
		fmt.Println("numberVal is", reflect.TypeOf(numberVal).Kind())

		if !methodOk || !numberOk || methodVal!="isPrime"||
		reflect.TypeOf(numberVal).Kind()!=reflect.Float64{
			
			fmt.Println("Packet is JSON but Malformed")
			conn.Write([]byte("Malformed JSON"))
			conn.Close()
			return 1
		}else{
			fmt.Println("Packet is JSON and well-formed")
			response:=make(map[string]interface{})
			response["Prime"]=isPrime(int(numberVal.(float64)))
			response["Method"]="isPrime"
			fmt.Println("response is", response)
			jsonResponse, _:=json.Marshal(response)
			conn.Write(jsonResponse)
			return 0
		}
	}}
