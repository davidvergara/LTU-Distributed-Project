
package dht

import (
	"net"
	"fmt"
	"runtime"
	"encoding/json"
	"strconv"
//	"time"
)

//Starts a go function that will be receiving messages
//for a node all the time
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
		
//		if receive.GetPort() == "1200" {
//			tick := time.Tick(40000 * time.Millisecond)
//			fmt.Println("Entering")
//
//			for {
//				fmt.Println("hey")
//				
//				buffer :=make([]byte,1024) 
//				readed, err := conn.Read(buffer)
//				if err != nil {
//					panic(err)
//				}
//	
//				message := buffer[0:readed]
//	
//				go receive.decryptMessage(message)
//				runtime.Gosched()
//				select{
//					case <- tick: 
//					fmt.Println("pam")
//					break;
//	//				case default 
//				}
//				fmt.Println("no nos hemos cargado el bucle")
//			} 
//		} else {
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
//		}
	}()
	
	go receive.StartHeartBeats()
	go receive.StartReplicateRoutine()
	go receive.StartUnreplicateRoutine()
}

//Decrypts a message received and calls the function that has 
//to work with that message
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
		case message.Type == "HEARTBEAT":
		{
			receive.receiveHeartBeat(message)
		}
		case message.Type == "HEARTBEATANSWER":
		{
			receive.receiveHeartBeatAnswer(message)
		}
		case message.Type == "SETDATA":
		{
			receive.receiveSetData(message)
		}
		case message.Type == "GETDATA":
		{
			receive.receiveGetData(message)
		}
		case message.Type == "GETDATAANSWER":
		{
			receive.receiveGetDataAnswer(message)
		}
		case message.Type == "DELETEDATA":
		{
			receive.receiveDeleteData(message)
		}
		case message.Type == "ADDDATA":
		{
			receive.receiveAddData(message)
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
	ring := message.Args["ring"]
	receive.PrintRingAux(origin,ring)
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

func (receive *DHTNode) receiveHeartBeat(message Msg){
	idHeartBeat := message.Args["heartBeatId"]
	receive.SendHeartBeatAnswer(message.Source,idHeartBeat)
}

func (receive *DHTNode) receiveHeartBeatAnswer(message Msg){
	idHeartBeat,_ := strconv.Atoi(message.Args["heartBeatId"])
	receive.HeartBeatRequest[idHeartBeat] <- message.Source
}

func (receive *DHTNode) receiveSetData(message Msg){
	dataToInsert := message.Data
	for k,v := range dataToInsert.DataStored{
		receive.Data.StoreData(k,v.Value,v.Original)
	}
}

func (receive *DHTNode) receiveGetData(message Msg){
	if message.Args["requestType"]=="original"{
		dataSetToBeSend :=MakeDataSet()
		for k,v := range  receive.Data.DataStored{
			if v.Original {
				dataSetToBeSend.StoreData(k,v.Value,false)
			}
		}
		receive.SendGetDataAnswer(message.Source,dataSetToBeSend,message.Args["getDataId"])
	}
}

func (receive *DHTNode) receiveGetDataAnswer (message Msg){
	idGetData,_ := strconv.Atoi(message.Args["getDataId"])
	receive.GetDataRequest[idGetData] <- message.Data
}

func (receive *DHTNode) receiveDeleteData(message Msg){
	dataToDelete := message.Data
	for k,_ := range dataToDelete.DataStored{
		receive.Data.deleteData(k)
	}
}

func (receive *DHTNode) receiveAddData(message Msg){
	for k,v := range message.Data.DataStored{
		dataSetToBeSend :=MakeDataSet()
		nodeResponsible:=receive.Lookup(k,receive.ToNetworkNode(),"")
		dataSetToBeSend.StoreData(k,v.Value,true)
		receive.SendSetData(nodeResponsible, dataSetToBeSend)
	}
}