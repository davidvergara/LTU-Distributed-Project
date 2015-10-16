package dht

import (
	"fmt"
	"bytes"
	"math/big"
	"encoding/hex"
)

/* Consts */
const SPACESIZE = 3


type Contact struct {
	ip   string
	port string
}

type DHTNode struct {
	nodeId   		string
	successor  		*NetworkNode
	predecessor		*NetworkNode
	contact			Contact
	fingers			[]*Finger
//	communication 	
}

type Finger struct {
	fingerId 		int
	nodeIdent		*NetworkNode
}

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

	dhtNode.successor = nil
	dhtNode.predecessor = nil
	dhtNode.fingers = make([]*Finger, SPACESIZE)
	
	//dhtNode.communication =  

	return dhtNode
}

func (dhtNode *DHTNode) AddToRing(newDHTNode *NetworkNode) {
	nodeResponsible := dhtNode.Lookup(newDHTNode.NodeId,dhtNode.ToNetworkNode(),"")
	if nodeResponsible == nil {
		fmt.Println("Error adding the node "+newDHTNode.NodeId)
	}
	dhtNode.SendInsertNodeBeforeMe(nodeResponsible,newDHTNode)
}

func (dhtNode *DHTNode) InsertNodeBeforeMe(newNode *NetworkNode) {
	mutexPredeccessor.Lock()
	valueNode,_ :=hex.DecodeString(dhtNode.nodeId)
	valueNodeNew,_ := hex.DecodeString(newNode.NodeId)
	if bytes.Compare(valueNode, valueNodeNew)==0 {
		/* The key is the same, error */
		fmt.Println("Error, tryed to add a nodeId that is in the ring "+newNode.NodeId)
		
	} else {
		valueNodePredecessor,_ := hex.DecodeString(dhtNode.predecessor.NodeId)
	
		if between(valueNodePredecessor, valueNode, valueNodeNew){
			/* Trying to insert the node in the right place */
			dhtNode.SendSetSuccessor(newNode, dhtNode.ToNetworkNode())
			dhtNode.SendSetPredecessor(newNode, dhtNode.predecessor)
			dhtNode.SendSetSuccessor(dhtNode.predecessor, newNode)
			dhtNode.SetPredecessor(newNode)
		} else{
			/* We have to look for the right place */
			dhtNode.AddToRing(newNode)
		}
	}
	mutexPredeccessor.Unlock()
}

func (dhtNode *DHTNode) updateFingerTables(){
	dhtNode.calcFingerTable()
	if dhtNode.GetSuccessor() != nil {
		
		/* More than one node in the ring */
		dhtNode.SendUpdateFingerTablesAux(dhtNode.ToNetworkNode(), dhtNode.GetSuccessor())
	}
}

func (dhtNode *DHTNode) updateFingerTablesAux(original *NetworkNode){
	if dhtNode.GetNodeId() != original.NodeId {
		dhtNode.calcFingerTable()
		dhtNode.SendUpdateFingerTablesAux(original, dhtNode.GetSuccessor())
	}
}

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

func (dhtNode *DHTNode) SetPredecessor(newPredecessor *NetworkNode){
	mutexPredeccessor.Lock()
	dhtNode.predecessor=newPredecessor
	mutexPredeccessor.Unlock()
}

func (dhtNode *DHTNode) SetSuccessor(newSuccessor *NetworkNode){
	mutexSuccessor.Lock()
	dhtNode.successor = newSuccessor
	mutexSuccessor.Unlock()
}

/* GETTERS of the node */
func (dhtNode *DHTNode) GetNodeId () string {
	return dhtNode.nodeId
}

func (dhtNode *DHTNode) GetSuccessor () *NetworkNode {
	return dhtNode.successor
}

func (dhtNode *DHTNode) GetPredecessor () *NetworkNode {
	return dhtNode.predecessor
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


func (dhtNode *DHTNode) Lookup(key string, sourceNode *NetworkNode, idLookup string) *NetworkNode {
	
	if dhtNode.nodeId == key || dhtNode.successor == nil {
		
		if (dhtNode.nodeId == sourceNode.NodeId){
			/* key == nodeID, I answer my query*/
			return dhtNode.ToNetworkNode()
		} else{
			dhtNode.SendLookupAnswer(dhtNode.ToNetworkNode() , sourceNode, idLookup)
			return nil
		}
	} else{
		
		keyBytes,_ := hex.DecodeString(key)
		nodeIDBytes,_ := hex.DecodeString(dhtNode.nodeId)
		sucessorIDBytes,_ := hex.DecodeString(dhtNode.successor.NodeId)
		
		if between(nodeIDBytes, sucessorIDBytes, keyBytes){

			/* key between nodeID and its successor */
			if (dhtNode.nodeId == sourceNode.NodeId){
				/* key == nodeID, I answer my query*/
				return dhtNode.successor
			} else{
				dhtNode.SendLookupAnswer(dhtNode.successor , sourceNode, idLookup)
				return nil
			}
		} else{
			/* key not between nodeID and its successor */
			
			/* return the closest finger to the key */
			dhtMinNode:= dhtNode.calcNodeMinDist(key)
			
			channel := dhtNode.SendLookup(key, dhtMinNode, sourceNode, idLookup)
			
			/* Waiting the answer */
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

/**
 * Return the closest finger of dhtNode to the key
 */
func (dhtNode *DHTNode) calcNodeMinDist(key string) *NetworkNode {
	dhtNodeMin := dhtNode.successor
	/* Key to HEX */
	keyBytes,_ := hex.DecodeString(key)
	/* dhtNodeMin to HEX */
	nodeIdBytes,_ := hex.DecodeString(dhtNodeMin.NodeId)
	minDist := distance(nodeIdBytes, keyBytes,SPACESIZE)
	for i,v := range dhtNode.fingers {
		
		if v!= nil {
			fingerBytes,_ := hex.DecodeString(dhtNode.fingers[i].nodeIdent.NodeId)
			
			if bytes.Compare(fingerBytes,nodeIdBytes) != 0 {
				if between(fingerBytes, nodeIdBytes,keyBytes){
					/* FingerID to HEX */
					distance:= distance(fingerBytes, keyBytes,SPACESIZE)
					if minDist.Cmp(distance) == 1{
						minDist=distance
						dhtNodeMin = dhtNode.fingers[i].nodeIdent
						//fmt.Println("Nodo " + dhtNode.nodeId + " Distancia minima -> " + dhtNodeMin.nodeId + " key " + key)
					}
				}
			}
			
		}
	}
	return dhtNodeMin
}


func (dhtNode *DHTNode) PrintRing(){
	fmt.Println("Node " + dhtNode.GetNodeId())
	if dhtNode.GetSuccessor() != nil {
		
		/* More than one node in the ring */
		dhtNode.SendPrintRingAux(dhtNode.ToNetworkNode(), dhtNode.GetSuccessor())
	}
}

func (dhtNode *DHTNode) PrintRingAux(original *NetworkNode){
	if dhtNode.GetNodeId() != original.NodeId {
		fmt.Println("Node " + dhtNode.GetNodeId())
		/* Not printed all the ring */
		dhtNode.SendPrintRingAux(original, dhtNode.GetSuccessor())
	}
}

func (dhtNode *DHTNode) responsible(key string) bool {
	nodeResponsible:= dhtNode.Lookup(key,dhtNode.ToNetworkNode(),"")
	return nodeResponsible.NodeId == dhtNode.nodeId
}

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




func (dhtNode *DHTNode) ToNetworkNode() *NetworkNode {
	networkNode := new(NetworkNode)
	networkNode.NodeId = dhtNode.GetNodeId()
	networkNode.Ip = dhtNode.GetIp()
	networkNode.Port = dhtNode.GetPort()
	return networkNode
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
