package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net"
	"reflect"
)


func main() {
	primeMap:= GenPrimes()
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
			continue
		}
		go handleConnection(c, primeMap)
	}
}

func handleConnection(conn net.Conn, primeMap map[int]bool) {
	fmt.Printf("Serving %s\n", conn.RemoteAddr().String())
	defer conn.Close()
	scn:=bufio.NewScanner(conn)
		

		for scn.Scan(){
				packet:=scn.Bytes()
				fmt.Println("packet is: ", string(packet))
				if(len(packet)!=0){
					connstat:=respondToPacket(packet, conn, primeMap)
					if(connstat==1){
						fmt.Println("conn closed")
						return
					}
				}
		
			}
		fmt.Println("Scanning ends")
		if err := scn.Err(); err != nil {
			log.Printf("Unexpected error: %s", err)
			return
		}
		return
		}
	


	


	func isPrime(num float64, mapOfPrimes map[int]bool) (prime bool) {
		if math.Trunc(num) != num {return false}
		if num <= 1 {return false}
		_,ok:=mapOfPrimes[int(num)]
		if ok{return true;}
		return false;
	 }

	type Response struct {
		Method string `json:"method"`
		Prime  bool   `json:"prime"`
	}

//0 means conn is not closed at the end, 1 means conn is closed
	func respondToPacket(packet []byte, conn net.Conn, primeMap map[int]bool) int {
	if !json.Valid(packet) {
		conn.Write([]byte("Malformed zsJSON"))	
		fmt.Println("packet is not JSON")
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



		if !methodOk || !numberOk || methodVal!="isPrime"||
		reflect.TypeOf(numberVal).Kind()!=reflect.Float64{
			
			fmt.Println("Packet is JSON but Malformed")
			conn.Write([]byte("Malformed JSON"))
			return 1
		}else{
			fmt.Println("Packet is well-formed")
			response:=Response{
				Method: "isPrime",
				Prime:	isPrime(numberVal.(float64), primeMap),
			}
	
			jsonResponse, _:=json.Marshal(response)
			fmt.Println(string(jsonResponse))
			conn.Write(jsonResponse)

			return 0
		}
	}}
