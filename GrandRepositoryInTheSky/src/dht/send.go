package dht

import (
	"net"
	"encoding/json"
)

func SetConnection(dest *NetworkNode) *net.UDPConn {
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


func Send(dest *NetworkNode, message Msg) {
	conn := SetConnection(dest)

	go func() {
		buffer, err := json.Marshal(message)
		
		if err !=nil {
			panic(err)
		}
		conn.Write(buffer)
		
	}()
}

func (dhtNode *DHTNode) SendLookup(key string, dhtMinNode *NetworkNode,
	 sourceNode *NetworkNode, idLookup string)chan *NetworkNode{

	if (dhtNode.nodeId == sourceNode.NodeId) {
		/* We need a channel to save the answer */
		mutexNumLookup.Lock()
		mess := Msg{Source: sourceNode,
			Dest: dhtMinNode,
 			Type: "LOOKUP", 
 			Args: map[string]string{
 					"key": string(key),
 					"lookUpId": string(NumLookup)}}
			
		answerChannel:= make(chan *NetworkNode)
		LookupRequest[NumLookup]=answerChannel
		NumLookup++
		mutexNumLookup.Unlock()
		Send(dhtMinNode, mess)
		return answerChannel
	} else {
		/* The node is just an intermediary */
		mess := Msg{Source: sourceNode,
			Dest: dhtMinNode,
 			Type: "LOOKUP", 
 			Args: map[string]string{
 					"key": string(key),
 					"lookUpId": idLookup}}
		Send(dhtMinNode, mess)
		return nil
	}
}

func (dhtNode *DHTNode) SendLookupAnswer(answerNode *NetworkNode, sourceNode *NetworkNode, idLookup string){

	mess := Msg{Source: answerNode,
			Dest: sourceNode,
	 		Type: "LOOKUPANSWER", 
	 		Args: map[string]string{
	 				"lookUpId":idLookup}}
		
	Send(sourceNode, mess)
}	

func (dhtNode *DHTNode) SendSetPredecessor(dest *NetworkNode, newPredecessor *NetworkNode){
	mess := Msg{Source: newPredecessor,
				Dest: dest,
				Type: "SETPREDECESSOR",
				Args: nil}
	
	Send(dest,mess)
}

func (dhtNode *DHTNode) SendSetSuccessor(dest *NetworkNode, newSuccessor *NetworkNode){
	mess := Msg{Source: newSuccessor,
			Dest: dest,
			Type: "SETSUCCESSOR",
			Args: nil}
	
	Send(dest, mess)
}

func (dhtNode *DHTNode) SendPrintRing(dest *NetworkNode){
	mess := Msg{Source: dhtNode.ToNetworkNode(),
				Dest: dest,
				Type: "PRINTRING",
				Args: nil}
	
	Send(dest,mess)
}

func (dhtNode *DHTNode) SendPrintRingAux(original *NetworkNode, dest *NetworkNode){
	mess := Msg{Source: original,
				Dest: dest,
				Type: "PRINTRINGAUX",
				Args: nil}
	
	Send(dest,mess)
}

func (dhtNode *DHTNode) SendUpdateFingerTables(dest *NetworkNode){
	mess := Msg{Source: dhtNode.ToNetworkNode(),
				Dest: dest,
				Type: "UPDATEFINGERTABLES",
				Args: nil}
	
	Send(dest,mess)
}

func (dhtNode *DHTNode) SendUpdateFingerTablesAux(original *NetworkNode, dest *NetworkNode){
	mess := Msg{Source: original,
				Dest: dest,
				Type: "UPDATEFINGERTABLESAUX",
				Args: nil}
	
	Send(dest,mess)
}


func (dhtNode *DHTNode) SendInsertNodeBeforeMe (nodeResponsible *NetworkNode,nodeToInsert *NetworkNode){
	mess := Msg{Source: nodeToInsert,
				Dest: nodeResponsible,
				Type: "INSERTNODEBEFOREME",
				Args: nil}
	Send(nodeResponsible,mess)
}