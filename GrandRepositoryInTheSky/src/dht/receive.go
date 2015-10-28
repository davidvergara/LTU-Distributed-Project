//Authors: Alejandro Márquez Ferrer / David Vergara Manrique
//Date: 28/10/2015
//
//Description: This file contains the code related to receive and handle the
//different messages, as well as functions for each type of message that can
//be received.

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
	
	go receive.StartHeartBeats()
	go receive.StartReplicateRoutine()
	go receive.StartUnreplicateRoutine()
	go receive.StartUpdateFingersRoutine()
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
		case message.Type == "SETDATAHTTP":
		{
			receive.receiveSetDataHttp(message)
		}
		case message.Type == "SETDATAHTTPANSWER":
		{
			receive.receiveSetDataHttpAnswer(message)
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
		case message.Type == "DELETEDATASUC":
		{
			receive.receiveDeleteDataSuc(message)
		}
		case message.Type == "DELETEDATAHTTP":
		{
			receive.receiveDeleteDataHttp(message)
		}
		case message.Type == "DELETEDATAHTTPANSWER":
		{
			receive.receiveDeleteDataHttpAnswer(message)
		}
		case message.Type == "ADDDATA":
		{
			receive.receiveAddData(message)
		}
		case message.Type == "PUTDATAHTTP":
		{
			receive.receivePutDataHttp(message)
		}
		case message.Type == "PUTDATAHTTPANSWER":
		{
			receive.receivePutDataHttpAnswer(message)
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

//Receive message HEARTBEAT
//Calls SendHeartBeatAnswer to answer the heartbeat
func (receive *DHTNode) receiveHeartBeat(message Msg){
	idHeartBeat := message.Args["heartBeatId"]
	receive.SendHeartBeatAnswer(message.Source,idHeartBeat)
}

//Receive message "HEARTBEATANSWER"
//Stores the node answered in the corresponding channel
func (receive *DHTNode) receiveHeartBeatAnswer(message Msg){
	idHeartBeat,_ := strconv.Atoi(message.Args["heartBeatId"])
	receive.HeartBeatRequest[idHeartBeat] <- message.Source
}

//Receive message "SETDATA"
//Stores all the data received
func (receive *DHTNode) receiveSetData(message Msg){
	dataToInsert := message.Data
	for k,v := range dataToInsert.DataStored{
		receive.Data.StoreData(k,v.Value,v.Original)
	}
}

//Receive message "GETDATA"
//Calls the function SendGetDataAnswer to send the corresponding answer to the sender
func (receive *DHTNode) receiveGetData(message Msg){
	if message.Args["requestType"]=="original"{
		dataSetToBeSend :=MakeDataSet()
		for k,v := range  receive.Data.DataStored{
			if v.Original {
				dataSetToBeSend.StoreData(k,v.Value,false)
			}
		}
		receive.SendGetDataAnswer(message.Source,dataSetToBeSend,message.Args["getDataId"])
	} else if message.Args["requestType"] == "all"{
		dataSetToBeSend :=MakeDataSet()
		dataSetToBeSend.DataStored = receive.Data.getStoredData()
		receive.SendGetDataAnswer(message.Source,dataSetToBeSend,message.Args["getDataId"])
	}
}

//Receive message "GETDATAANSWER"
//Stores the answer in the corresponding channel
func (receive *DHTNode) receiveGetDataAnswer (message Msg){
	idGetData,_ := strconv.Atoi(message.Args["getDataId"])
	receive.GetDataRequest[idGetData] <- message.Data
}

//Receive message "DELETEDATA"
//Deletes all the data specified and tells its successor to do the same
func (receive *DHTNode) receiveDeleteData(message Msg){
	dataToDelete := message.Data
	receive.SendDeleteDataSuc(receive.Successor, dataToDelete)
	for k,_ := range dataToDelete.DataStored{
		receive.Data.deleteData(k)
	}
}

//Receive message "DELETEDATASUC"
//Deletes all the data specified
func (receive *DHTNode) receiveDeleteDataSuc(message Msg){
	dataToDelete := message.Data
	for k,_ := range dataToDelete.DataStored{
		receive.Data.deleteData(k)
	}
}

//Receive message "ADDDATA"
//Stores the all the data that was sent 
func (receive *DHTNode) receiveAddData(message Msg){
	for k,v := range message.Data.DataStored{
		dataSetToBeSend :=MakeDataSet()
		nodeResponsible:=receive.Lookup(k,receive.ToNetworkNode(),"")
		dataSetToBeSend.StoreData(k,v.Value,true)
		receive.SendSetData(nodeResponsible, dataSetToBeSend)
	}
}

//Receive message "SETDATAHTTP"
//Stores the data passed as message and sends an answer to the source node
func (receive *DHTNode) receiveSetDataHttp(message Msg){
	idSetData := message.Args["setDataId"]
	dataToInsert := message.Data
	var exito bool
	for k,v := range dataToInsert.DataStored{
		exito = receive.Data.StoreData(k,v.Value,v.Original)
	}
	receive.SendSetDataAnswer(message.Source,idSetData,exito)
}

//Receive message "SETDATAHTTPANSWER"
//Puts in the result of the operation setData in the corresponding channel
func (receive *DHTNode) receiveSetDataHttpAnswer(message Msg){
	idSetData,_ := strconv.Atoi(message.Args["setDataId"])
	receive.SetDataRequest[idSetData] <- message.Args["bool"] == "true"
}

//Receive message "PUTDATAHTTP"
//Updates the data passed as message and sends an answer to the source node
func (receive *DHTNode) receivePutDataHttp(message Msg){
	idPutData := message.Args["putDataId"]
	dataToPut := message.Data
	var exito bool
	for k,v := range dataToPut.DataStored{
		exito = receive.Data.updateData(k,v.Value)
		if exito {
			
			/* If data was updated successfully, then update the replica too */
			data,_ := receive.Data.getData(k)
			if data.Original{
				if receive.Successor != nil {
					dataToPut.changeOriginalReplica(k)
					receive.SendDeleteDataSuc(receive.Successor,dataToPut)
					receive.SendSetData(receive.Successor,dataToPut)
				}
			}
		}
	}
	receive.SendPutDataAnswer(message.Source,idPutData,exito)
}

//Receive message "PUTDATAHTTPANSWER"
//Puts in the result of the operation putData in the corresponding channel
func (receive *DHTNode) receivePutDataHttpAnswer(message Msg){
	idPutData,_ := strconv.Atoi(message.Args["putDataId"])
	receive.PutDataRequest[idPutData] <- message.Args["bool"] == "true"
}

//Receive message "DELETEDATAHTTP"
//Deletes the data passed in the message and sends an answer to the source node
func (receive *DHTNode) receiveDeleteDataHttp(message Msg){
	idDeleteData := message.Args["deleteDataId"]
	dataToDelete := message.Data
	var exito bool
	for k,_ := range dataToDelete.DataStored{
		exito = receive.Data.deleteData(k)
		if exito {
			if receive.Successor != nil {
				
				/* If data deleted successfully, delete the replica too */
				receive.SendDeleteDataSuc(receive.Successor,dataToDelete)
			}
		}
	}
	receive.SendDeleteDataAnswer(message.Source,idDeleteData,exito)
}

//Receive message "DELETEDATAHTTPANSWER"
//Puts in the result of the operation deleteData in the corresponding channel
func (receive *DHTNode) receiveDeleteDataHttpAnswer(message Msg){
	idDeleteData,_ := strconv.Atoi(message.Args["deleteDataId"])
	receive.DeleteDataRequest[idDeleteData] <- message.Args["bool"] == "true"
}