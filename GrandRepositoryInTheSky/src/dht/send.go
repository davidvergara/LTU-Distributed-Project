package dht

import (
	"net"
	"encoding/json"
	  "io/ioutil"
	  "fmt"
)

//type NodeSender struct{
//	dhtNode					*dht.DHTNode
//}

func SetConnection(dest *DHTNode) *net.UDPConn {
	addr, err := net.ResolveUDPAddr("udp", dest.GetIp()+":"+dest.GetPort())
	if err != nil {
		panic(err)
	}
	connection, err := net.DialUDP("udp", nil, addr)
	
	if err != nil {
		panic(err)
	}
	return connection
}


func Send(dest *DHTNode, message Msg) {
	conn := SetConnection(dest)

	
	go func() {
			d1 := []byte("hello\ngo\n")
    es := ioutil.WriteFile("/Temp/sddsd", d1, 0644)
     fmt.Println(es)
		buffer, err := json.Marshal(message)
		
		if err !=nil {
			panic(err)
		}
		conn.Write(buffer)
		
	}()

}