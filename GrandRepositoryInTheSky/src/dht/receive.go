
package dht

import (
	"net"
	"fmt"
	"runtime"
	"encoding/json"
	"strconv"
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

			go receive.decryptMessage(message)
			runtime.Gosched()
		}
	}()
}

func (receive *DHTNode) decryptMessage (bytesReceived []byte){
	var message Msg
	err := json.Unmarshal(bytesReceived, &message)
	if err != nil {
		panic (err)
	}
	switch 
	{
		case message.Type == "LOOKUP":
		{
			receive.receiveLookup(message)
		}
		case message.Type == "LOOKUPANSWER":
		{
			receive.receiveLookupAnswer(message)
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
			receive.receiveSetPredecessor(message)
		}
		case message.Type == "SETSUCCESSOR":
		{
			receive.receiveSetSuccessor(message)
		}
		case message.Type == "PRINTRING":
		{
			//Llamar funcion PRINTRING
		}
		default: 
		{
			fmt.Println("Wrong message")
		}
	}
}
	
func (receive *DHTNode) receiveLookup (message Msg){
	receive.Lookup(message.Args["key"],message.Source, message.Args["lookUpId"]) 
}

func (receive *DHTNode) receiveLookupAnswer (message Msg){
	idLookup,_ := strconv.Atoi(message.Args["lookUpId"])
	LookupRequest[idLookup] <- message.Source
}

func (receive *DHTNode) receiveSetPredecessor (message Msg){
	newPredecessor := message.Source
	receive.SetPredecessor(newPredecessor)
}

func (receive *DHTNode) receiveSetSuccessor (message Msg){
	newSuccessor := message.Source
	receive.SetSuccessor(newSuccessor)
}

