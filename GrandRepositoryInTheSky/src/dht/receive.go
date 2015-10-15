
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
			fmt.Println("ssdsdsdsdsdsdsdsdsd")
			if err != nil {
				
				panic(err)
			}

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
		//Llamar funcion LOOKUP
	}
	case message.Type == "LOOKUPRESPONDE":
	{
		//Llamar funcion SETSUCCESSOR
	}
	case message.Type == "UPDATEFINGERS":
	{
		//Llamar funcion UPDATEFINGERS
	}
	case message.Type == "ADDRING":
	{
		//Llamar funcion ADDRING
	}
	case message.Type == "SETPREDECESSOR":
	{
		//Llamar funcion SETPREDECESSOR
	}
	case message.Type == "SETSUCCESSOR":
	{
		//Llamar funcion SETSUCCESSOR
	}
	default: 
	{
		fmt.Println("Wrong message")
	}
	}
	

}
