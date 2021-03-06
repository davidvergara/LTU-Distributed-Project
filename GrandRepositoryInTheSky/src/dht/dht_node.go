//Authors: Alejandro Márquez Ferrer / David Vergara Manrique
//Date: 28/10/2015
//
//Description: This file contains the code related to each node. It is here where
//the dht node is defined, and all operations related to them are in this file.

package dht

import (
	"fmt"
	"bytes"
	"math/big"
	"encoding/hex"
	"sync"
	"time"
	"strconv"
)

/* Consts */
const SPACESIZE = 3

const LOOKUPEXPIRATION = time.Second * 5
const HEARTBEATPERIOD = time.Second * 5
const HEARTBEATEXPIRATION = time.Second * 4
const GETDATAEXPIRATION = time.Second * 4
const REPLICATEPERIOD = time.Second * 10
const UNREPLICATEPERIOD = time.Second * 20
const UPDATEFINGERPERIOD = time.Second * 10

/* Contact struct */
type Contact struct {
	ip   string
	port string
}

/* Node struct */
type DHTNode struct {
	nodeId   		string
	Successor  		*NetworkNode
	Predecessor		*NetworkNode
	PredOfPred		*NetworkNode	
	contact			Contact
	fingers			[]*Finger
	Data			DataSet
	
	/* Channels and ids for the different answers */
	NumLookup int
	LookupRequest 	map[int]chan *NetworkNode
	NumHeartBeat int
	HeartBeatRequest map[int]chan *NetworkNode
	NumGetData int
	GetDataRequest map[int]chan DataSet
	NumSetData int
	SetDataRequest map[int]chan bool
	NumPutData int
	PutDataRequest map[int]chan bool
	NumDeleteData int
	DeleteDataRequest map[int]chan bool
	
	/* Mutex part */
	mutexNumLookup  sync.Mutex
	mutexNumHeartBeat sync.Mutex
	mutexPredeccessor  sync.Mutex
	mutexSuccessor  sync.Mutex
	mutexPredOfPred	sync.Mutex
	mutexNumGetData sync.Mutex
	mutexSetData sync.Mutex
	mutexPutData sync.Mutex
	mutexDeleteData sync.Mutex
}

/* Finger struct */
type Finger struct {
	fingerId 		int
	nodeIdent		*NetworkNode
}


//Creates a DHTNode with the atributes passed as parameters, and initializes
//some global variables.
func MakeDHTNode(nodeId *string, ip string, port string) *DHTNode {
	dhtNode := new(DHTNode)
	dhtNode.contact.ip = ip
	dhtNode.contact.port = port

	if nodeId == nil {
		genNodeId := generateNodeId()
		dhtNode.nodeId = genNodeId
	} else {
		dhtNode.nodeId = *nodeId
	}

	dhtNode.Successor = nil
	dhtNode.Predecessor = nil
	dhtNode.PredOfPred = nil
	dhtNode.fingers = make([]*Finger, SPACESIZE)
	
	dhtNode.NumLookup=0
	dhtNode.NumHeartBeat=0
	dhtNode.NumGetData=0
	dhtNode.NumSetData=0
	dhtNode.NumPutData=0
	dhtNode.NumDeleteData=0
	//dhtNode.LookupRequest= make map[int]chan *NetworkNode

	//var LookupRequest 	map[int]chan *NetworkNode
	
	dhtNode.Data = MakeDataSet()

	dhtNode.LookupRequest = make(map[int]chan *NetworkNode)
	dhtNode.HeartBeatRequest = make(map[int]chan *NetworkNode)
	dhtNode.GetDataRequest = make(map[int]chan DataSet)
	dhtNode.SetDataRequest = make(map[int]chan bool)
	dhtNode.PutDataRequest = make(map[int]chan bool)
	dhtNode.DeleteDataRequest = make(map[int]chan bool)
	
	return dhtNode
}

/* GETTERS of the node */
func (dhtNode *DHTNode) GetNodeId () string {
	return dhtNode.nodeId
}

func (dhtNode *DHTNode) GetSuccessor () *NetworkNode {
	return dhtNode.Successor
}

func (dhtNode *DHTNode) GetPredecessor () *NetworkNode {
	return dhtNode.Predecessor
}

func (dhtNode *DHTNode) GetContact () Contact {
	return dhtNode.contact
}

func (dhtNode *DHTNode) GetIp () string {
	return dhtNode.contact.ip
}

func (dhtNode *DHTNode) GetPort () string {
	return dhtNode.contact.port
}

//Local node looks for the place to insert the new node passed as parameter,
//and tells the node responsible of that place to insert it.
func (dhtNode *DHTNode) AddToRing(newDHTNode *NetworkNode) {
	nodeResponsible := dhtNode.Lookup(newDHTNode.NodeId,dhtNode.ToNetworkNode(),"")
	if nodeResponsible == nil {
		
		/* Can't insert */
		fmt.Println("Error adding the node "+newDHTNode.NodeId)
	}
	dhtNode.SendInsertNodeBeforeMe(nodeResponsible,newDHTNode)
}

//Local node inserts the new node before it if that place is still valid. If not,
//calls AddToRing again to reallocate the node.
func (dhtNode *DHTNode) InsertNodeBeforeMe(newNode *NetworkNode) {
	dhtNode.mutexPredeccessor.Lock()
	valueNode,_ :=hex.DecodeString(dhtNode.nodeId)
	valueNodeNew,_ := hex.DecodeString(newNode.NodeId)
	if bytes.Compare(valueNode, valueNodeNew)==0 {
		
		/* Can't be 2 nodes with same ID */
		fmt.Println("hola jeje")
		fmt.Println("Error, tried to add a nodeId that is already in the ring: "+newNode.NodeId)
		
	} else {
		
		if(dhtNode.Predecessor == nil){
			
			/* First node in the ring */
			dhtNode.SendSetSuccessor(newNode, dhtNode.ToNetworkNode())
			dhtNode.SendSetPredecessor(newNode, dhtNode.ToNetworkNode())
			dhtNode.Successor = newNode
			dhtNode.Predecessor = newNode
			dhtNode.updateFingerTables()
			
			/* Send data corresponding to new node */
			dhtNode.UpdateAndSendData(newNode, true)
			
		} else{
			
			valueNodePredecessor,_ := hex.DecodeString(dhtNode.Predecessor.NodeId)
			if bytes.Compare(valueNodePredecessor, valueNodeNew) == 0{
				
				/* Can't be 2 nodes with same ID */
				fmt.Println("Error, tried to add a nodeId that is already in the ring: "+newNode.NodeId)
			} else if between(valueNodePredecessor, valueNode, valueNodeNew){
				
				/* Trying to insert the node in the right place */
				dhtNode.SendSetSuccessor(newNode, dhtNode.ToNetworkNode())
				dhtNode.SendSetPredecessor(newNode, dhtNode.Predecessor)
				dhtNode.SendSetSuccessor(dhtNode.Predecessor, newNode)
				dhtNode.Predecessor = newNode
				dhtNode.updateFingerTables()
				
				/* Send data corresponding to new node */
				dhtNode.UpdateAndSendData(newNode, false)
			} else{
				/* This is not the right place anymore -> we have to find it out again */
				dhtNode.AddToRing(newNode)
			}
		}
	}
	dhtNode.mutexPredeccessor.Unlock()
}

//Updates data of the actual node and sends data corresponding to new node to him
func (dhtNode *DHTNode) UpdateAndSendData(newNode *NetworkNode, onlyOneNode bool){
	
	/* Send data corresponding to new node */
	dataSetToBeSend :=MakeDataSet()
	dataToBeDeleted :=MakeDataSet()
	for k,v := range dhtNode.Data.DataStored{
		
		/* Send your replicas to the new node (as replica) */
		if !v.Original {
			dataSetToBeSend.StoreData(k,v.Value,false)
			
			/* Delete our replicas */
			dhtNode.Data.deleteData(k)
		
		}
	}
	
	newNodeValue,_ := hex.DecodeString(newNode.NodeId)
	actualNodeValue,_ := hex.DecodeString(dhtNode.nodeId)
	
	for k,v := range dhtNode.Data.DataStored{
		if v.Original {
			kValue,_ := hex.DecodeString(k)
			if k != dhtNode.nodeId{
				if k == newNode.NodeId {
					
					//Send the corresponding original data to the new node
					dataSetToBeSend.StoreData(k,v.Value,true)
					//Store replica
					dhtNode.Data.changeOriginalReplica(k)
					if !onlyOneNode {
						
						dataToBeDeleted.StoreData(k,v.Value,false)
					}						
				} else if !between(newNodeValue, actualNodeValue, kValue){
					dataSetToBeSend.StoreData(k,v.Value,true)
					dhtNode.Data.changeOriginalReplica(k)
					//Store replica
					if !onlyOneNode {
						
						dataToBeDeleted.StoreData(k,v.Value,false)
					}		
				}
			}
		}
	}
	if !onlyOneNode {
		
		//Sends the successor the datas he has to delete
		dhtNode.SendDeleteData(dhtNode.Successor,dataToBeDeleted)
	}
	
	//Send the new node its data
	dhtNode.SendSetData(newNode,dataSetToBeSend)
}

//Update the finger table of all the nodes in the ring
func (dhtNode *DHTNode) updateFingerTables(){
	dhtNode.calcFingerTable()
	if dhtNode.Successor != nil {
		
		//More than one node in the ring
		dhtNode.SendUpdateFingerTablesAux(dhtNode.ToNetworkNode(), dhtNode.Successor)
	}
}

//Continues updating the finger table of all the nodes in the ring
func (dhtNode *DHTNode) updateFingerTablesAux(original *NetworkNode){
	if dhtNode.GetNodeId() != original.NodeId {
		
		/* Not reached the beginning of the ring */
		dhtNode.SendUpdateFingerTablesAux(original, dhtNode.Successor)
	}
}

//Update the finger table of the Local node
func (dhtNode *DHTNode) calcFingerTable (){
	for k:=1; k<=SPACESIZE; k++ {
		n,_ :=hex.DecodeString(dhtNode.nodeId)
		idFinger,_ :=calcFinger(n, k, SPACESIZE)
		nodeFinger:= dhtNode.Lookup(idFinger,dhtNode.ToNetworkNode(),"")
		if nodeFinger != nil {
			dhtNode.fingers[k-1] = new(Finger)
			dhtNode.fingers[k-1].fingerId=k
			dhtNode.fingers[k-1].nodeIdent=nodeFinger
		} else {
			dhtNode.fingers[k-1] = new(Finger)
			dhtNode.fingers[k-1].fingerId=k
			dhtNode.fingers[k-1].nodeIdent=dhtNode.ToNetworkNode()
		}
	}
}

//Puts newPredecessor as the predecessor of the local node
func (dhtNode *DHTNode) SetPredecessor(newPredecessor *NetworkNode){
	dhtNode.mutexPredeccessor.Lock()
	dhtNode.Predecessor=newPredecessor
	dhtNode.mutexPredeccessor.Unlock()
}

//Puts local node's successor as local node's predecessor's successor of successor
//Puts newSucessor as the successor of the local node
func (dhtNode *DHTNode) SetSuccessor(newSuccessor *NetworkNode){
	dhtNode.mutexSuccessor.Lock()
	dhtNode.Successor = newSuccessor
	dhtNode.mutexSuccessor.Unlock()
}

//Puts newPredOfPred as the predecessor of the predecessor of the local node
func (dhtNode *DHTNode) SetPredOfPred(newPredOfPred *NetworkNode){
	dhtNode.mutexPredOfPred.Lock()
	dhtNode.PredOfPred = newPredOfPred
	dhtNode.mutexPredOfPred.Unlock()
}

//Returns the node responsible for the key passed as parameter
//SourceNode: Node that invoqued lookup originally
//idLookup: Id of the request (may be empty)
func (dhtNode *DHTNode) Lookup(key string, sourceNode *NetworkNode, idLookup string) *NetworkNode {
	
	if dhtNode.nodeId == key || dhtNode.Successor == nil {
		
		/* I am responsible for the key || I am the only node in the ring */
		if (dhtNode.nodeId == sourceNode.NodeId){
			
			/* key == nodeID, I answer (me) my query*/
			return dhtNode.ToNetworkNode()
		} else{
			
			/* I send the answer (me) to the node that asked originally for it */
			dhtNode.SendLookupAnswer(dhtNode.ToNetworkNode() , sourceNode, idLookup)
			return dhtNode.ToNetworkNode()
		}
	} else{
		
		keyBytes,_ := hex.DecodeString(key)
		nodeIDBytes,_ := hex.DecodeString(dhtNode.nodeId)
		sucessorIDBytes,_ := hex.DecodeString(dhtNode.Successor.NodeId)
		
		if between(nodeIDBytes, sucessorIDBytes, keyBytes){

			/* key between nodeID and its successor */
			if (dhtNode.nodeId == sourceNode.NodeId){
				
				/* key == nodeID, I answer (my successor) my query*/
				return dhtNode.Successor
			} else{
				
				/* I send the answer (my successor) to the node that asked originally for it */
				dhtNode.SendLookupAnswer(dhtNode.Successor , sourceNode, idLookup)
				return dhtNode.Successor
			}
		} else{
			
			/* key not between nodeID and its successor */
			/* return the closest finger to the key */
			dhtMinNode:= dhtNode.calcNodeMinDist(key)
			channel := dhtNode.SendLookup(key, dhtMinNode, sourceNode, idLookup)
			
			/* Waiting the answer in the channel*/
			if (dhtNode.nodeId == sourceNode.NodeId) {
				select {
					case answer := <- channel:
						return answer
					case <-time.After(LOOKUPEXPIRATION):
						fmt.Println("Waiting time for lookup answer expirated")
						return nil
				}
				return nil
			} else {
				
				/* The node that called the function is just and intermediary */
				return nil
			}
		}
	}
}

//Return the closest finger of dhtNode to the key
func (dhtNode *DHTNode) calcNodeMinDist(key string) *NetworkNode {
	dhtNodeMin := dhtNode.Successor
	/* Key to HEX */
	keyBytes,_ := hex.DecodeString(key)
	/* dhtNodeMin to HEX */
	nodeIdBytes,_ := hex.DecodeString(dhtNodeMin.NodeId)
	minDist := distance(nodeIdBytes, keyBytes,SPACESIZE)
	
	/* Iterates over all the values of dhtNode.fingers */
	for i,v := range dhtNode.fingers {
		if v!= nil{
			if v.nodeIdent!= nil {
				
				/* finger[i] != nil */
				fingerBytes,_ := hex.DecodeString(dhtNode.fingers[i].nodeIdent.NodeId)
				
				if bytes.Compare(fingerBytes,nodeIdBytes) != 0 {
					if between(fingerBytes, nodeIdBytes,keyBytes){
						/* FingerID to HEX */
						distance:= distance(fingerBytes, keyBytes,SPACESIZE)
						if minDist.Cmp(distance) == 1{
							
							/* Updating closest finger */
							minDist=distance
							dhtNodeMin = dhtNode.fingers[i].nodeIdent
						}
					}
				}
			}
		}
	}
	return dhtNodeMin
}

//Prints myself and sends printRingAux message to my successor
func (dhtNode *DHTNode) PrintRing(){
	ring := fmt.Sprintln("Printing ring...")
	
	data := "Data: " + "\n"
	for k,v := range dhtNode.Data.DataStored{
		data = data + k + " " + v.Value + " " + strconv.FormatBool(v.Original) + "\n"
	}
	
	ring = ring + fmt.Sprintln("Node " + dhtNode.GetNodeId() + " - " + dhtNode.GetPort()) + fmt.Sprintln(data)
	
	/* Uncomment to print fingers */
//	ring = ring + dhtNode.PrintFingerTable()
	if dhtNode.GetSuccessor() != nil {
		
		/* More than one node in the ring */
		
		dhtNode.SendPrintRingAux(dhtNode.ToNetworkNode(), dhtNode.GetSuccessor(), ring)
	} else{
		fmt.Println(ring)
	}
}

//If I was not the first node printing the ring, prints myself and sends
//printRingAux to my successor
func (dhtNode *DHTNode) PrintRingAux(original *NetworkNode, ring string){
	if dhtNode.GetNodeId() != original.NodeId {
			data := "Data: " + "\n"
			for k,v := range dhtNode.Data.DataStored{
				data = data + k + " " + v.Value + " " + strconv.FormatBool(v.Original) + "\n"
			}
		/* Not printed all the ring */
		ring = ring + fmt.Sprintln("Node " + dhtNode.GetNodeId() + " - " + dhtNode.GetPort()) + fmt.Sprintln(data)
		
		/* Uncomment to print fingers */
//		ring = ring + dhtNode.PrintFingerTable()

		dhtNode.SendPrintRingAux(original, dhtNode.GetSuccessor(),ring)
	} else{
		fmt.Println(ring)
	}
}

//Returns true if this node is responsible for the key passed as parameter.
//Otherwise, returns false
func (dhtNode *DHTNode) responsible(key string) bool {
	nodeResponsible:= dhtNode.Lookup(key,dhtNode.ToNetworkNode(),"")
	return nodeResponsible.NodeId == dhtNode.nodeId
}

//Prints the finger k of this node (SPACESIZE = m)
func (dhtNode *DHTNode) PrintFinger(k int, m int){
	//	fmt.Println("calculating result = (n+2^(k-1)) mod (2^m)")

	// convert the n to a bigint
	nBigInt := big.Int{}
	n,_ := hex.DecodeString(dhtNode.nodeId)
	nBigInt.SetBytes(n)

	fmt.Printf("n            %s\n",dhtNode.nodeId)

	fmt.Printf("k            %d\n", k)

	fmt.Printf("m            %d\n", m)

	// get the right addend, i.e. 2^(k-1)
	two := big.NewInt(2)
	addend := big.Int{}
	addend.Exp(two, big.NewInt(int64(k-1)), nil)

	fmt.Printf("2^(k-1)      %s\n", addend.String())

	// calculate sum
	sum := big.Int{}
	sum.Add(&nBigInt, &addend)

	fmt.Printf("(n+2^(k-1))  %s\n", sum.String())

	// calculate 2^m
	ceil := big.Int{}
	ceil.Exp(two, big.NewInt(int64(m)), nil)

	fmt.Printf("2^m          %s\n", ceil.String())

	// apply the mod
	result := big.Int{}
	result.Mod(&sum, &ceil)
	
	resultBytes := result.Bytes()
	if len(resultBytes) == 0 {
		resultBytes = []byte{0}
	}
	resultHex := fmt.Sprintf("%x", resultBytes)

	fmt.Printf("result       %s\n", result.String())
	fmt.Printf("successor    %s\n", dhtNode.Lookup(resultHex,dhtNode.ToNetworkNode(),"").NodeId)
}

//Prints all the fingers of this node
func (dhtNode *DHTNode) PrintFingerTable() string{
	result := "\n"
	for i,v := range dhtNode.fingers {
		if v != nil {
			result = result + fmt.Sprintf("   -Finger %d -> %s\n",v.fingerId,v.nodeIdent.NodeId)
		} else{
			result = result + fmt.Sprintf("   -Finger %d -> null\n",i+1)
		}
		
	}
	return result
}

//Starts the heartbeat routine
func (dhtNode *DHTNode) StartHeartBeats(){
	for {
		time.Sleep(HEARTBEATPERIOD)
		if dhtNode.Predecessor != nil{
			
			/* Channel to receive our PredPred */
			channel := dhtNode.SendHeartBeat(dhtNode.Predecessor)
			select {
				case answer := <- channel:
				{	
					/* Our predecessor is alive */
					dhtNode.SetPredOfPred(answer)
				}	
				case <-time.After(HEARTBEATEXPIRATION):
				{
					
					/* Our predecessor is dead */
					dhtNode.DeadPredecessor()
				}					
			}
		}
	}
}

//When our predecessor is dead, sets your PredOfPred as your predecessor. Also 
//Change the replicas of that node to original
func (dhtNode *DHTNode) DeadPredecessor(){
	dhtNode.mutexPredeccessor.Lock()
	fmt.Println("Node " + dhtNode.Predecessor.NodeId + " died")
	if dhtNode.PredOfPred == nil || dhtNode.PredOfPred.NodeId == dhtNode.nodeId {
		
		/* This is the only node remaining in the ring */
		dhtNode.Successor = nil
		dhtNode.Predecessor = nil
		dhtNode.PredOfPred = nil
		
		/* Changing my replicas to data */
		for k,v := range dhtNode.Data.DataStored{
			if !v.Original {
				dhtNode.Data.changeReplicaOriginal(k)
			}
		}
	} else{
		
		/* There are more nodes in the ring */
		dhtNode.Predecessor = dhtNode.PredOfPred
		dhtNode.SendSetSuccessor(dhtNode.Predecessor,dhtNode.ToNetworkNode())
		dhtNode.PredOfPred = nil
		
		/* Changing my replicas to data from the node
		   and send them as replicas to the successor */
		dataSetToBeSend :=MakeDataSet()
		for k,v := range dhtNode.Data.DataStored{
			if !v.Original {
				dhtNode.Data.changeReplicaOriginal(k)
				dataSetToBeSend.StoreData(k,v.Value,false)
			}
		}
		dhtNode.SendSetData(dhtNode.Successor,dataSetToBeSend)
		
		/* Ask the predecessor for its original data */
		c:=dhtNode.SendGetData("original",dhtNode.Predecessor)
		select {
				case answer := <- c:
				{	
					/* Data received, storing... */
					for k,v := range answer.DataStored{
						dhtNode.Data.StoreData(k,v.Value,false)
					}
				}	
				case <-time.After(GETDATAEXPIRATION):
				{
					
					/* Expiration time */
					fmt.Println("ERROR: Predecessor does not answer...")
				}					
		}
	}
	dhtNode.updateFingerTables()
	dhtNode.mutexPredeccessor.Unlock()
}

//Starts the routine in charge of replicate data
func (dhtNode *DHTNode) StartReplicateRoutine(){
	for {
		time.Sleep(REPLICATEPERIOD)
		if dhtNode.Successor != nil {
			dataSetToBeSend :=MakeDataSet()
			for k,v := range dhtNode.Data.DataStored{
				if v.Original {
					
					/* Replicates your original data to the successor */
					dataSetToBeSend.StoreData(k,v.Value,false)
				}
			}
			dhtNode.SendSetData(dhtNode.Successor,dataSetToBeSend)
		}
	}
}

//Starts the routine in charge of delete data replicated more than once
func (dhtNode *DHTNode) StartUnreplicateRoutine(){
	for {
		time.Sleep(UNREPLICATEPERIOD)
		if dhtNode.Predecessor != nil && dhtNode.PredOfPred != nil {
			for k,v := range dhtNode.Data.DataStored{
				if !v.Original {
					valuePred,_ :=hex.DecodeString(dhtNode.Predecessor.NodeId)
					valuePredOfPred,_ := hex.DecodeString(dhtNode.PredOfPred.NodeId)
					valueK,_ := hex.DecodeString(k)
					
					if bytes.Compare(valueK,valuePredOfPred) == 0{
						
						/* Data replicated more than once */
						dhtNode.Data.deleteData(k)
					} else if bytes.Compare(valueK,valuePred) != 0{
						if !between(valuePredOfPred, valuePred, valueK) {
							
							/* Data replicated more than once */
							dhtNode.Data.deleteData(k)
						}
					}
				}
			}
		}
	}
}

//Starts the routine in charge of update the local node's finger table
func (dhtNode *DHTNode) StartUpdateFingersRoutine(){
	for {
		time.Sleep(UPDATEFINGERPERIOD)
		dhtNode.calcFingerTable()
	}
}

//Add a pair [key,value] to a node
//The key and value are inserted as parameter
//Return true if all was correct and the node where was inserted the pair
//Return false if something failed
func (dhtNode *DHTNode) HttpPost(key string, value string) (bool, *NetworkNode){
	nodeResponsible := dhtNode.Lookup(key, dhtNode.ToNetworkNode(), "")
	dataSetToBeSend :=MakeDataSet()
	dataSetToBeSend.StoreData(key,value,true)
	
	channel := dhtNode.SendSetDataWithAnswer(nodeResponsible,dataSetToBeSend)
	
	/* Waiting the answer in the channel*/
	select {
		case answer := <- channel:
			return answer, nodeResponsible
		case <-time.After(LOOKUPEXPIRATION):
			fmt.Println("Waiting time for SendSetDataWithAnswer answer expirated")
			return false, nodeResponsible
	}
}

//Find a key insterted as parameter in the ring
//Return true if the key exists, the Data of the key and
//the node where is it saved
//If the key does not exist return false
func (dhtNode *DHTNode) HttpGet(key string) (bool, string, *NetworkNode){
	nodeResponsible := dhtNode.Lookup(key, dhtNode.ToNetworkNode(), "")
	
	channel:= dhtNode.SendGetData("all",nodeResponsible)
		
	/* Waiting the answer in the channel*/
	select {
		case dataSet := <- channel:
			data, success :=dataSet.getData(key)
			return success, data.Value, nodeResponsible
		case <-time.After(LOOKUPEXPIRATION):
			fmt.Println("Waiting time for SendGetData answer expirated")
			return false, "", nodeResponsible
	}
}

//Update the value of the data of a key
//The key and the value are inserted as parameter
//Return true if all was correct and the node where was made the update
//Return false if there was an error
func (dhtNode *DHTNode) HttpPut(key string, value string) (bool, *NetworkNode){
	nodeResponsible := dhtNode.Lookup(key, dhtNode.ToNetworkNode(), "")
	dataSetToBeSend :=MakeDataSet()
	dataSetToBeSend.StoreData(key,value,true)
	channel := dhtNode.SendPutDataWithAnswer(nodeResponsible, dataSetToBeSend)
	
	/* Waiting the answer in the channel*/
	select {
		case answer := <- channel:
			return answer, nodeResponsible
		case <-time.After(LOOKUPEXPIRATION):
			fmt.Println("Waiting time for SendPutDataWithAnswer answer expirated")
			return false, nodeResponsible
	}
}

//Delete a pain [key,value] of the ring
//The key is inserted as parameter
//Return true and the node where was deleted the pair if all was
//correct 
//Return false if there was an error
func (dhtNode *DHTNode) HttpDelete(key string) (bool, *NetworkNode){
	nodeResponsible := dhtNode.Lookup(key, dhtNode.ToNetworkNode(), "")
	dataSetToBeSend :=MakeDataSet()
	dataSetToBeSend.StoreData(key,"",true)
	channel := dhtNode.SendDeleteDataWithAnswer(nodeResponsible,dataSetToBeSend)
	
	/* Waiting the answer in the channel*/
	select {
		case answer := <- channel:
			return answer, nodeResponsible
		case <-time.After(LOOKUPEXPIRATION):
			fmt.Println("Waiting time for SendDeleteDataWithAnswer answer expirated")
			return false, nodeResponsible
	}
}

func (dhtNode *DHTNode) testCalcFingers(m int, bits int) {
	/* idBytes, _ := hex.DecodeString(dhtNode.nodeId)
	fingerHex, _ := calcFinger(idBytes, m, bits)
	fingerSuccessor := dhtNode.lookup(fingerHex)
	fingerSuccessorBytes, _ := hex.DecodeString(fingerSuccessor.nodeId)
	fmt.Println("successor    " + fingerSuccessor.nodeId)

	dist := distance(idBytes, fingerSuccessorBytes, bits)
	fmt.Println("distance     " + dist.String()) */
}
