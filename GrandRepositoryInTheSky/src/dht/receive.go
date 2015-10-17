
package dht

import (
	"net"
	"fmt"
	"runtime"
	"encoding/json"
	"strconv"
)

//Starts a go function that will be receiving messages
//for a node all the time
func (receive *DHTNode) StartListenServer(){
	
	addr, err := net.ResolveUDPAddr("udp",(":"+receive.GetPort()))
	if err != nil {
		fmt.Println("tolai")
		panic(err)
	}
	
	conn, err := net.ListenUDP("udp",addr)
	if err != nil {
		fmt.Println("hijo")
		panic(err)
	}
	
	go func() {
		
		for {
			
			buffer :=make([]byte,1024) 
			readed, err := conn.Read(buffer)
			if err != nil {
				fmt.Println("de")
				panic(err)
			}

			message := buffer[0:readed]

			go receive.decryptMessage(message)
			runtime.Gosched()
		}
	}()
}

//Decrypts a message received and calls the function that has 
//to work with that message
func (receive *DHTNode) decryptMessage (bytesReceived []byte){
	var message Msg
	err := json.Unmarshal(bytesReceived, &message)
	if err != nil {
		fmt.Println("puta")
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
		case message.Type == "ADDTORING":
		{
			receive.receiveAddToRing(message)
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

//Received meesage 	"LOOKUP"
//A Lookup function will be run with the arguments of the message
func (receive *DHTNode) receiveLookup (message Msg){
	receive.Lookup(message.Args["key"],message.Source, message.Args["lookUpId"]) 
}

//Received meesage 	"LOOKUPANSWER"
//A answer to a Lookup has been received, saving that 
//answer in the corresponding channel
func (receive *DHTNode) receiveLookupAnswer (message Msg){
	idLookup,_ := strconv.Atoi(message.Args["lookUpId"])
	receive.LookupRequest[idLookup] <- message.Source
}

//Received message "SETPREDECESSOR"
//Uploading predecessor of the node...
func (receive *DHTNode) receiveSetPredecessor (message Msg){
	newPredecessor := message.Source
	receive.SetPredecessor(newPredecessor)
}

//Received message "SETSUCCESSOR"
//Uploading successor of the node...
func (receive *DHTNode) receiveSetSuccessor (message Msg){
	newSuccessor := message.Source
	receive.SetSuccessor(newSuccessor)
}

//Received message "PRINTRING"
//Calling PrintRing function to print the ring
func (receive *DHTNode) receivePrintRing(message Msg){
	receive.PrintRing()
}

//Received message "PRINTRINGAUX"
//Calling PrintRingAux function to print the ring
func (receive *DHTNode) receivePrintRingAux(message Msg){
	origin := message.Source
	receive.PrintRingAux(origin)
}

//Received message "ADDTORING"
//Calling AddToRing funtion to insert a node in the ring
func (receive *DHTNode) receiveAddToRing(message Msg){
	newNode := message.Source
	receive.AddToRing(newNode)
}

//Received message "UPDATEFINGERTABLES"
//Calling updateFingerTables to update the finger tables of the node
func (receive *DHTNode) receiveUpdateFingerTables(message Msg){
	receive.updateFingerTables()
}

//Received message "UPDATEFINGERTABLESAUX"
//Calling updateFingerTablesAux function to update the finger tables
//of the node
func (receive *DHTNode) receiveUpdateFingerTablesAux(message Msg){
	origin := message.Source
	receive.updateFingerTablesAux(origin)
}

//Received message "INSERTNODEBEFOREME"
//Calling InsertNodeBeforeMe function to insert a node before
//the local node
func (receive *DHTNode) receiveInsertNodeBeforeMe(message Msg){
	nodeToInsert := message.Source
	receive.InsertNodeBeforeMe(nodeToInsert)
}