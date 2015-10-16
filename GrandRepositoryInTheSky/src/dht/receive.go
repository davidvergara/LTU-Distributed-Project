
package dht

import (
	"net"
	"fmt"
	"runtime"
	"encoding/json"
	"strconv"
)

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
		case message.Type == "ADDTORING":
		{
			receive.receiveAddToRing(message)
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
			receive.receivePrintRing(message)
		}
		case message.Type == "PRINTRINGAUX":
		{
			receive.receivePrintRingAux(message)
		}
		case message.Type == "UPDATEFINGERTABLES":
		{
			receive.receiveUpdateFingerTables(message)
		}
		case message.Type == "UPDATEFINGERTABLESAUX":
		{
			receive.receiveUpdateFingerTablesAux(message)
		}
		case message.Type == "INSERTNODEBEFOREME":
		{
			receive.receiveInsertNodeBeforeMe(message)
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

func (receive *DHTNode) receivePrintRing(message Msg){
	receive.PrintRing()
}

func (receive *DHTNode) receivePrintRingAux(message Msg){
	origin := message.Source
	receive.PrintRingAux(origin)
}

func (receive *DHTNode) receiveAddToRing(message Msg){
	newNode := message.Source
	receive.AddToRing(newNode)
}

func (receive *DHTNode) receiveUpdateFingerTables(message Msg){
	receive.updateFingerTables()
}

func (receive *DHTNode) receiveUpdateFingerTablesAux(message Msg){
	origin := message.Source
	receive.updateFingerTablesAux(origin)
}

func (receive *DHTNode) receiveInsertNodeBeforeMe(message Msg){
	nodeToInsert := message.Source
	receive.InsertNodeBeforeMe(nodeToInsert)
}