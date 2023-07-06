package main
//Each client tracks a different asset
//Each connection is a separate session, each session a different asset


import (
"net", 
"fmt",
"encoding/hex",
"encoding/binary",

)

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
		go handle(c)
	}
}

func handle(conn net.Conn){
	prices:=make(map[int32]int32)
	defer conn.Close()
	packet:=make([]byte, 9)
	for{
		bytesRead, err:=io.ReadFull(conn, packet)
		if err!=nil{
			if err==io.EOF{
				fmt.Println("EOF")
			}
			println("error reading from connection")
			break
		}
		type, first, second, err:=parse(packet)
		switch type{
		case "I":
				prices[first]=insertPrice(second)
		case "Q":
				keys=rangeOfKeys(first, second, prices)
				sum:=0
				for _, k := range keys{
					sum+=prices[k]
				}
				conn.Write(sum/len(keys))

			}

}

func rangeOfKeys(start int32, end int32, prices map[int32]int32) int32[]{
	keys := make([]int, len(prices))
	i := 0
	for k := range prices {
		keys[i] = k
		i++
	}
	sort.Ints(keys)
	for i, k := range keys {
		if(k>=first){
			key=keys[i:]
		}
		if(k>second){
			key=keys[:i]
		}
	}
	return keys
}

func parse(msg []bytes) (type byte, first int32, second int 32, err error){
	if len(msg)!=9{
		return 0,0,0,errors.New("invalid message length")
	}
	fmt.Println("message is ", msg[0]) 
	switch msg[0] {
	case 0x71:
		type="Q"
	case 0x49:
		type="I"
	default:
		return 0,0,0,errors.New("invalid message type")

	first=binary.BigEndian.Uint32(msg[1:5])
	second=binary.BigEndian.Uint32(msg[5:9])
	return type, first, second, nil	
	}
}




//handleConnection


//writeInsert
