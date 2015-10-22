package dht

import (
	"net"
	"encoding/json"
	"fmt"
	"strconv"
)

//Set an UDP connection with the node passed as parameter
func SetConnection(dest *NetworkNode) *net.UDPConn {
//	fmt.Println("Hola")
//	if dest == nil {
//		fmt.Println("dest es null")
//	}
	addr, err := net.ResolveUDPAddr("udp", dest.Ip+":"+dest.Port)
	if err != nil {
		panic(err)
	}
	connection, err := net.DialUDP("udp", nil, addr)
	
	if err != nil {
		panic(err)
	}
	return connection
}

//Sends the message passed as parameter to the destination
func Send(dest *NetworkNode, message Msg) {
	
	if dest == nil {
		fmt.Println(message.Source.NodeId)
		fmt.Println(message.Type)
//		fmt.Println(message.Args)
	}
	conn := SetConnection(dest)

	go func() {
		buffer, err := json.Marshal(message)
		
		if err !=nil {
			panic(err)
		}
		
		conn.Write(buffer)
		conn.Close()
		
	}()
}

func (dhtNode *DHTNode) SendLookup(key string, dhtMinNode *NetworkNode,
	 sourceNode *NetworkNode, idLookup string)chan *NetworkNode{

	if (dhtNode.nodeId == sourceNode.NodeId) {
		/* We need a channel to save the answer */
		dhtNode.mutexNumLookup.Lock()
		numLookupString := strconv.Itoa(dhtNode.NumLookup)
		mess := Msg{Source: sourceNode,
			Dest: dhtMinNode,
 			Type: "LOOKUP", 
 			Args: map[string]string{
 					"key": string(key),
 					"lookUpId": numLookupString},
 			Data: DataSet{}}
			
		answerChannel:= make(chan *NetworkNode)
		dhtNode.LookupRequest[dhtNode.NumLookup]=answerChannel
		dhtNode.NumLookup++
		dhtNode.mutexNumLookup.Unlock()
		Send(dhtMinNode, mess)
		return answerChannel
	} else {
		/* The node is just an intermediary */
		mess := Msg{Source: sourceNode,
			Dest: dhtMinNode,
 			Type: "LOOKUP", 
 			Args: map[string]string{
 					"key": string(key),
 					"lookUpId": idLookup},
 			Data: DataSet{}}
		Send(dhtMinNode, mess)
		return nil
	}
}

//Sends to the destination (sourceNode) a LOOKUPANSER type message. 
//AnswerNode: node responsible for the key that a node was looking for
//SourceNode: node that made the first lookup request
//idLookup: id of the lookup request, to be stored in that channel
func (dhtNode *DHTNode) SendLookupAnswer(answerNode *NetworkNode, sourceNode *NetworkNode, idLookup string){

	mess := Msg{Source: answerNode,
			Dest: sourceNode,
	 		Type: "LOOKUPANSWER", 
	 		Args: map[string]string{
	 				"lookUpId":idLookup},
 			Data: DataSet{}}
		
	Send(sourceNode, mess)
}	

//Sends to the destination a SETPREDECESSOR message
func (dhtNode *DHTNode) SendSetPredecessor(dest *NetworkNode, newPredecessor *NetworkNode){
	mess := Msg{Source: newPredecessor,
				Dest: dest,
				Type: "SETPREDECESSOR",
				Args: nil,
 				Data: DataSet{}}
	
	Send(dest,mess)
}

//Sends to the destination a SETSUCCESSOR message
func (dhtNode *DHTNode) SendSetSuccessor(dest *NetworkNode, newSuccessor *NetworkNode){
	mess := Msg{Source: newSuccessor,
			Dest: dest,
			Type: "SETSUCCESSOR",
			Args: nil,
 			Data: DataSet{}}
	
	Send(dest, mess)
}

//Sends to the destination a PRINTRING message (starting printing ring)
func (dhtNode *DHTNode) SendPrintRing(dest *NetworkNode){
	mess := Msg{Source: dhtNode.ToNetworkNode(),
				Dest: dest,
				Type: "PRINTRING",
				Args: nil,
 				Data: DataSet{}}
	
	Send(dest,mess)
}

//Sends to the destination a PRINTRINGAUX message (continuing printing ring)
func (dhtNode *DHTNode) SendPrintRingAux(original *NetworkNode, dest *NetworkNode, ring string){
	mess := Msg{Source: original,
				Dest: dest,
				Type: "PRINTRINGAUX",
				Args:  map[string]string{
	 				"ring":ring},
 				Data: DataSet{}}
	
	Send(dest,mess)
}

//Sends to the destination a ADDTORING message.
//NewNode: node to be added
func (dhtNode *DHTNode) SendAddToRing(dest *NetworkNode, newNode *NetworkNode){
	mess := Msg{Source: newNode,
				Dest: dest,
				Type: "ADDTORING",
				Args: nil,
 				Data: DataSet{}}
	
	Send(dest,mess)
}

func SendAddToRingForeign(destIP string, destPort string, newNode *NetworkNode){
	auxNetwork := new(NetworkNode)
	auxNetwork.Ip = destIP
	auxNetwork.Port = destPort
	auxNetwork.NodeId = ""
	mess := Msg{Source: newNode,
				Dest: auxNetwork,
				Type: "ADDTORING",
				Args: nil,
 				Data: DataSet{}}
	
	Send(auxNetwork,mess)
}

func SendDataToRingForeign(destIP string, destPort string,data DataSet){
	auxNetwork := new(NetworkNode)
	auxNetwork.Ip = destIP
	auxNetwork.Port = destPort
	auxNetwork.NodeId = ""
	mess := Msg{Source: new(NetworkNode),
				Dest: auxNetwork,
				Type: "ADDDATA",
				Args: nil,
 				Data: data}
	
	Send(auxNetwork,mess)
}

//Sends to the destination a UPDATEFINGERTABLES message (starting updating fingers)
func (dhtNode *DHTNode) SendUpdateFingerTables(dest *NetworkNode){
	mess := Msg{Source: dhtNode.ToNetworkNode(),
				Dest: dest,
				Type: "UPDATEFINGERTABLES",
				Args: nil,
 				Data: DataSet{}}
	
	Send(dest,mess)
}

//Sends to the destination a UPDATEFINGERTABLESAUX message (continuing updating finger)
func (dhtNode *DHTNode) SendUpdateFingerTablesAux(original *NetworkNode, dest *NetworkNode){
	if dest != nil {
		mess := Msg{Source: original,
					Dest: dest,
					Type: "UPDATEFINGERTABLESAUX",
					Args: nil,
	 				Data: DataSet{}}
		
		Send(dest,mess)
	}
}

//Sends to the destination a SENDINSERTNODEBEFOREME message
//NodeToInsert: node to be inserted
//NodeResponsible: node that has to insert the node
func (dhtNode *DHTNode) SendInsertNodeBeforeMe (nodeResponsible *NetworkNode,nodeToInsert *NetworkNode){
	mess := Msg{Source: nodeToInsert,
				Dest: nodeResponsible,
				Type: "INSERTNODEBEFOREME",
				Args: nil,
 				Data: DataSet{}}
	Send(nodeResponsible,mess)
}

func (dhtNode *DHTNode) SendHeartBeat(dest *NetworkNode)chan *NetworkNode{

	//fmt.Println("+Nodo " + dhtNode.nodeId + " sending heartbeat to " + dest.NodeId)
	/* We need a channel to save the answer */
	dhtNode.mutexNumHeartBeat.Lock()
	numHeartBeat := strconv.Itoa(dhtNode.NumHeartBeat)
	mess := Msg{Source: dhtNode.ToNetworkNode(),
			Dest: dest,
 			Type: "HEARTBEAT", 
 			Args: map[string]string{"heartBeatId": numHeartBeat},
 			Data: DataSet{}}
		
	answerChannel:= make(chan *NetworkNode)
	dhtNode.HeartBeatRequest[dhtNode.NumHeartBeat]=answerChannel
	dhtNode.NumHeartBeat++
	dhtNode.mutexNumHeartBeat.Unlock()
	if dhtNode.GetPort() == "1201"{
		fmt.Println("Enviando latido a... " + dhtNode.Predecessor.Port)
	}
	Send(dest, mess)
	return answerChannel
}

func (dhtNode *DHTNode) SendHeartBeatAnswer(dest *NetworkNode, idHeartBeat string){

	//fmt.Println("-Nodo " + dhtNode.nodeId + " answering heartbeat to " + dest.NodeId + " with answer " + dhtNode.Predecessor.NodeId)
	mess := Msg{Source: dhtNode.Predecessor,
			Dest: dest,
	 		Type: "HEARTBEATANSWER", 
	 		Args: map[string]string{
	 				"heartBeatId":idHeartBeat},
 			Data: DataSet{}}
		
	Send(dest, mess)
}

func (dhtNode *DHTNode) SendSetData(dest *NetworkNode, datasetToSend DataSet){
	
		mess := Msg{Source: dhtNode.ToNetworkNode(),
			Dest: dest,
	 		Type: "SETDATA", 
	 		Args: nil,
 			Data: datasetToSend} 
		
		Send(dest,mess)
}

func (dhtNode *DHTNode) SendGetData(typeData string, dest *NetworkNode)chan DataSet{

	/* We need a channel to save the answer */
	dhtNode.mutexNumGetData.Lock()
	numGetData:= strconv.Itoa(dhtNode.NumGetData)
	mess := Msg{Source: dhtNode.ToNetworkNode(),
			Dest: dest,
 			Type: "GETDATA", 
 			Args: map[string]string{"getDataId": numGetData,
 									"requestType": typeData},
 			Data: DataSet{}}
		
	answerChannel:= make(chan DataSet)
	dhtNode.GetDataRequest[dhtNode.NumGetData] =  answerChannel
	dhtNode.NumGetData++
	dhtNode.mutexNumGetData.Unlock()
	Send(dest, mess)
	return answerChannel
}

func (dhtNode *DHTNode) SendGetDataAnswer(dest *NetworkNode, dataSet DataSet, dataId string){

	mess := Msg{Source: dhtNode.ToNetworkNode(),
			Dest: dest,
	 		Type: "GETDATAANSWER", 
	 		Args: map[string]string{
	 				"getDataId":dataId},
 			Data: dataSet}
		
	Send(dest, mess)
}	

func (dhtNode *DHTNode) SendDeleteData(dest *NetworkNode, datasetToSend DataSet){
	
		mess := Msg{Source: dhtNode.ToNetworkNode(),
			Dest: dest,
	 		Type: "DELETEDATA", 
	 		Args: nil,
 			Data: datasetToSend} 
		
		Send(dest,mess)
}
