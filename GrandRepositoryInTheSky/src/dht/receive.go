
package dht

import (
	"net"

	"fmt"
    "io/ioutil"

	"runtime"
	"encoding/json"
	
)

//type nodeReceiver struct{
//	dhtNode				*dht.DHTNode
//	sendDataObject		*sendData.NodeSender
//}

func (receive *DHTNode) StartListenServer(){
	
	addr, err := net.ResolveUDPAddr("udp",(":"+receive.GetPort()))
	if err != nil {
		panic(err)
	}
	
	conn, err := net.ListenUDP("udp",addr)
	if err != nil {
		panic(err)
	}
	go func() {
		for {
			buffer :=make([]byte,1024) 
			readed, err := conn.Read(buffer)
			if err != nil {
				panic(err)
			}
			d1 := []byte("hello\ngo\n")
    es := ioutil.WriteFile("/Temp/sdssssss", d1, 0644)
     fmt.Println(es)
			message := buffer[0:readed]

			receive.decryptMessage(message)
			runtime.Gosched()
		}
	}()
}

func (receive *DHTNode) decryptMessage (bytesReceived []byte){
	d1 := []byte("hello\ngo\n")
    es := ioutil.WriteFile("/Temp/sddsd", d1, 0644)
     fmt.Println(es)
	var message Msg
	err := json.Unmarshal(bytesReceived, &message)
	if err != nil {
		panic (err)
	}
	switch 
	{
	case message.Type == "LOOKUP":
	{
		 d1 := []byte("hello\ngo\n")
    err := ioutil.WriteFile("/Temp/dat1", d1, 0644)
     fmt.Println(err)
		fmt.Println("recibido")
	}
	default: 
	{
		fmt.Println("nada")
	}
	}
	

}
