package dht

import (
	"fmt"
	"bytes"
	"math/big"
	"encoding/hex"
	"sync"
)

/* Consts */
const SPACESIZE = 160
//var NumLookup = 0
//var LookupRequest 	map[int]chan *NetworkNode
//var mutexNumLookup = &sync.Mutex{}
//var mutexPredeccessor = &sync.Mutex{}
//var mutexSuccessor = &sync.Mutex{} 

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
	contact			Contact
	fingers			[]*Finger
	NumLookup int
	LookupRequest 	map[int]chan *NetworkNode
	mutexNumLookup  sync.Mutex
	mutexPredeccessor  sync.Mutex
	mutexSuccessor  sync.Mutex
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
	dhtNode.fingers = make([]*Finger, SPACESIZE)
	
	dhtNode.NumLookup=0
	//dhtNode.LookupRequest= make map[int]chan *NetworkNode

	//var LookupRequest 	map[int]chan *NetworkNode

	dhtNode.LookupRequest = make(map[int]chan *NetworkNode)
	
	
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
		fmt.Println("Error, tried to add a nodeId that is already in the ring: "+newNode.NodeId)
		
	} else {
		if(dhtNode.Predecessor == nil){
			
			/* First node in the ring */
			dhtNode.SendSetSuccessor(newNode, dhtNode.ToNetworkNode())
			dhtNode.SendSetPredecessor(newNode, dhtNode.ToNetworkNode())
			dhtNode.Predecessor = newNode
			dhtNode.Successor = newNode
			dhtNode.updateFingerTables()
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
			} else{
				/* This is not the right place anymore -> we have to find it out again */
				dhtNode.AddToRing(newNode)
			}
		}
	}
	dhtNode.mutexPredeccessor.Unlock()
}

//Update the finger table of all the nodes in the ring
func (dhtNode *DHTNode) updateFingerTables(){
	dhtNode.calcFingerTable()
	if dhtNode.GetSuccessor() != nil {
		
		/* More than one node in the ring */
		dhtNode.SendUpdateFingerTablesAux(dhtNode.ToNetworkNode(), dhtNode.GetSuccessor())
	}
}

//Continues updating the finger table of all the nodes in the ring
func (dhtNode *DHTNode) updateFingerTablesAux(original *NetworkNode){
	if dhtNode.GetNodeId() != original.NodeId {
		
		/* Not reached the beginning of the ring */
		dhtNode.calcFingerTable()
		dhtNode.SendUpdateFingerTablesAux(original, dhtNode.GetSuccessor())
	}
}

//Update the finger table of the Local node
func (dhtNode *DHTNode) calcFingerTable (){
	for k:=1; k<=SPACESIZE; k++ {
		n,_ :=hex.DecodeString(dhtNode.nodeId)
		idFinger,_ :=calcFinger(n, k, SPACESIZE)
		nodeFinger:= dhtNode.Lookup(idFinger,dhtNode.ToNetworkNode(),"")
		dhtNode.fingers[k-1] = new(Finger)
		dhtNode.fingers[k-1].fingerId=k
		dhtNode.fingers[k-1].nodeIdent=nodeFinger
	}
}

//Puts newPredecessor as the predecessor of the local node
func (dhtNode *DHTNode) SetPredecessor(newPredecessor *NetworkNode){
	dhtNode.mutexPredeccessor.Lock()
	dhtNode.Predecessor=newPredecessor
	dhtNode.mutexPredeccessor.Unlock()
}

//Puts newSucessor as the successor of the local node
func (dhtNode *DHTNode) SetSuccessor(newSuccessor *NetworkNode){
//	fmt.Println("Nodo " + dhtNode.nodeId + " poniendo sucesor " + newSuccessor.NodeId)
	dhtNode.mutexSuccessor.Lock()
	dhtNode.Successor = newSuccessor
	dhtNode.mutexSuccessor.Unlock()
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
			return nil
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
				return nil
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
		
		if v!= nil {
			
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
	return dhtNodeMin
}

//Prints myself and sends printRingAux message to my successor
func (dhtNode *DHTNode) PrintRing(){
	fmt.Println("Node " + dhtNode.GetNodeId())
	//dhtNode.PrintFingerTable()
	if dhtNode.GetSuccessor() != nil {
		
		/* More than one node in the ring */
		dhtNode.SendPrintRingAux(dhtNode.ToNetworkNode(), dhtNode.GetSuccessor())
	}
}

//If I was not the first node printing the ring, prints myself and sends
//printRingAux to my successor
func (dhtNode *DHTNode) PrintRingAux(original *NetworkNode){
	if dhtNode.GetNodeId() != original.NodeId {
		
		/* Not printed all the ring */
		fmt.Println("Node " + dhtNode.GetNodeId())
		
		/* TO TEST PRINTING FINGER TABLE WITH TEST22() */
//		if dhtNode.GetPort() == "1180" || dhtNode.GetPort() == "1160"{
//			dhtNode.PrintFingerTable()
//		}

		dhtNode.SendPrintRingAux(original, dhtNode.GetSuccessor())
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
func (dhtNode *DHTNode) PrintFingerTable() {
	for i,v := range dhtNode.fingers {
		if v != nil {
			fmt.Printf("   -Finger %d -> %s\n",v.fingerId,v.nodeIdent.NodeId)
		} else{
			fmt.Printf("   -Finger %d -> null\n",i+1)
		}
		
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
