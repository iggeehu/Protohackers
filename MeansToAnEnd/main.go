package main

//Each client tracks a different asset
//Each connection is a separate session, each session a different asset

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"sort"
)

type Price struct{
	Time uint32
	Amount uint32
}

type ByTime []Price

func (p ByTime) Len() int           { return len(p) }
func (p ByTime) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p ByTime) Less(i, j int) bool { return p[i].Time < p[j].Time }


//main.go
func main(){
	PORT:=":10000"
	l, err := net.Listen("tcp4", PORT)
	if err != nil {	
		log.Fatal(err)
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("Listening...")
		go handle(c)
	}
}

func handle(conn net.Conn){
	fmt.Printf("Serving %s\n", conn.RemoteAddr().String())
	record:=make([]Price, 0)
	defer conn.Close()
	packet:=make([]byte, 9)
	for{
		_, err:=io.ReadFull(conn, packet)
		if err!=nil{
			if err==io.EOF{
				fmt.Println("EOF")
			}
			println("error reading from connection")
			break
		}
		ptype, first, second, err:=parse(packet)
		switch ptype{
		case "I":
				curr:=Price{first, second}
				record=append(record, curr)
		case "Q":
				deliverable := make([]byte, 4)

				sort.Sort(ByTime(record))
				bi := sort.Search(len(record), func(i int) bool { return record[i].Amount >= first })
				ei:= sort.Search(len(record), func(i int) bool { return record[i].Amount > second })
				//If there are no samples within the requested period, or if mintime comes after maxtime, the value returned must be 0.
				if bi==len(record) || ei==len(record){
					binary.LittleEndian.PutUint32(deliverable, 0)
					conn.Write(deliverable)
					break
				}
				var sum uint32=0
				for i:=bi; i<=ei-1; i++{
						sum+=record[i].Amount
				}
				avg:=sum/(uint32)(ei-bi)
				binary.LittleEndian.PutUint32(deliverable, avg)
				_, err = conn.Write(deliverable)
		
			}
	}
}



func parse(msg []byte) (ptype string, first uint32, second uint32, err error){
	if len(msg)!=9{
		return "error",0,0,errors.New("invalid message length")
	}
	fmt.Println("message is ", msg[0]) 
	switch msg[0] {
	case 0x71:
		ptype="Q"
	case 0x49:
		ptype="I"
	default:
		return "",0,0,errors.New("invalid message type")
	}
	first=binary.BigEndian.Uint32(msg[1:5])
	second=binary.BigEndian.Uint32(msg[5:9])
	fmt.Println("first is ", first, "second is ", second, "type is ", ptype)
	return ptype, first, second, nil	
	
}


