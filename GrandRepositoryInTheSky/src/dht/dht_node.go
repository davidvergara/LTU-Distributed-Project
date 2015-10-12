package dht

import (
	"fmt"
	"bytes"
	"math/big"
)


/* Consts */
const SPACESIZE = 3


type Contact struct {
	ip   string
	port string
}

type DHTNode struct {
	nodeId      string
	successor   *DHTNode
	predecessor *DHTNode
	contact     Contact
	fingers		[]*Finger
}

type Finger struct {
	fingerId 		int
	nodeIdent		*DHTNode
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

	return dhtNode
}

func (dhtNode *DHTNode) AddToRing(newDHTNode *DHTNode) {
	
	/* Just one node in the ring ->
	   newNode is successor and predecessor of that node */
	if dhtNode.successor == nil {
		dhtNode.successor = newDHTNode
		dhtNode.predecessor = newDHTNode
		newDHTNode.successor = dhtNode
		newDHTNode.predecessor = dhtNode
		dhtNode.updateFingerTables()
		
	} else {
		
		/* More than one node */
		valueNode :=[]byte(dhtNode.nodeId)
		valueNodeNew := []byte(newDHTNode.nodeId)
		valueNodeNext := []byte(dhtNode.successor.nodeId)
		
		/* Look if dhtNode is last node in the ring */
		if bytes.Compare(valueNode, valueNodeNext) == 1 {
			
			/* New node between last and first nodes */
			if bytes.Compare(valueNodeNew,valueNode) == 1 ||
				 bytes.Compare(valueNodeNew,valueNodeNext) == -1 {
				
				/* valueNodeNew > valueNode || valueNodeNew < valueNodeNext */
				oldSuccessorDhtNode := dhtNode.successor
				dhtNode.successor = newDHTNode
				newDHTNode.successor = oldSuccessorDhtNode
				newDHTNode.predecessor = dhtNode
				oldSuccessorDhtNode.predecessor = newDHTNode
				dhtNode.updateFingerTables()
			} else {
				
				/* New node is not after last node ->
			       recursion with first node */
				dhtNode.successor.AddToRing(newDHTNode)
			}
		} else {
			
			/* Trying to insert between 2 consecutive nodes */
			if bytes.Compare(valueNodeNew,valueNode) == 1 && 
				bytes.Compare(valueNodeNew,valueNodeNext) == -1{
				/* valueNodeNew > valueNode && valueNodeNew < valueNodeNext */
				
				/* New node id bigger than dhtNode id and smaller than next node id ->
			       inserction between those nodes */
					oldSuccessorDhtNode := dhtNode.successor
					dhtNode.successor = newDHTNode
					newDHTNode.successor = oldSuccessorDhtNode
					newDHTNode.predecessor = dhtNode
					oldSuccessorDhtNode.predecessor = newDHTNode
					dhtNode.updateFingerTables()
			} else {
				
				/* New node is not between those nodes ->
			       recursion with next node */
				dhtNode.successor.AddToRing(newDHTNode)
			}	
		}
	}
}

func (dhtNode *DHTNode) updateFingerTables() {
	dhtNode.calcFingerTable()
	
	/* There is more than one node */
	if dhtNode.successor != nil {
		dhtNode.successor.updateFingerTablesAux(dhtNode.nodeId)
	}
}

func(dhtNode * DHTNode) updateFingerTablesAux(nodeID string) {
	if dhtNode.nodeId != nodeID {
		
		/* This is not the first node */
		dhtNode.calcFingerTable()
		dhtNode.successor.updateFingerTablesAux(nodeID)
	}
}

func (dhtNode *DHTNode) calcFingerTable (){
	for k:=1; k<=SPACESIZE; k++ {
		idFinger,_ :=calcFinger([]byte(dhtNode.nodeId), k, SPACESIZE)
		nodeFinger:= dhtNode.acceleratedLookupUsingFingers(idFinger)
		dhtNode.fingers[k-1].fingerId=k
		dhtNode.fingers[k-1].nodeIdent=nodeFinger
	}
}

func (dhtNode *DHTNode) Lookup(key string) *DHTNode {
	
	if dhtNode.nodeId == key || dhtNode.successor == nil {
		/* key == nodeID */
		return dhtNode
	} else if between([]byte(dhtNode.nodeId), 
		[]byte(dhtNode.successor.nodeId),
		[]byte(key)){
	
		/* key between nodeID and its successor */
		return dhtNode.successor
	} else{
		/* key not between nodeID and its successor */
		return dhtNode.successor.Lookup(key)
	}
}

/* GETTERS of the node */
func (dhtNode *DHTNode) GetNodeId () string {
	return dhtNode.nodeId
}

func (dhtNode *DHTNode) GetSuccessor () *DHTNode {
	return dhtNode.successor
}

func (dhtNode *DHTNode) GetPredecessor () *DHTNode {
	return dhtNode.predecessor
}

func (dhtNode *DHTNode) GetContact () Contact {
	return dhtNode.contact
}


func (dhtNode *DHTNode) acceleratedLookupUsingFingers(key string) *DHTNode {
	// TODO
	
	
	if dhtNode.nodeId == key || dhtNode.successor == nil {
		/* key == nodeID */
		return dhtNode
		
	} else if between([]byte(dhtNode.nodeId), 
		[]byte(dhtNode.successor.nodeId),
		[]byte(key)){
	
		/* key between nodeID and its successor */
		return dhtNode.successor
	} else{
		/* key not between nodeID and its successor */
		
		/* return the closest finger to the key */
		dhtMinNode:= dhtNode.calcNodeMinDist(key)
		return dhtMinNode.acceleratedLookupUsingFingers(key)
	}
	
}

/**
 * Return the closest finger of dhtNode to the key
 */
func (dhtNode *DHTNode) calcNodeMinDist(key string) *DHTNode {
	dhtNodeMin := dhtNode.fingers[0].nodeIdent
	minDist := distance([]byte(dhtNodeMin.nodeId), []byte(key),SPACESIZE)
	for i:=0; i<len(dhtNode.fingers); i++ {
		distance:= distance([]byte(dhtNode.fingers[i].nodeIdent.nodeId), []byte(key),SPACESIZE)
		if minDist.Cmp(distance) == 1{
			minDist=distance
			dhtNodeMin = dhtNode.fingers[i].nodeIdent
		}
	}
	return dhtNodeMin
}

/* Return the responsible node for the key */
func (dhtNode *DHTNode) responsible(key string) bool {
	nodeResponsible:= dhtNode.acceleratedLookupUsingFingers(key)
	return nodeResponsible.nodeId == dhtNode.nodeId
}

func (dhtNode *DHTNode) PrintFinger(k int, m int){
	fmt.Println("calculating result = (n+2^(k-1)) mod (2^m)")

	// convert the n to a bigint
	nBigInt := big.Int{}
	nBigInt.SetBytes([]byte(dhtNode.nodeId))

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
	resultHex := fmt.Sprintf("%x", resultBytes)

	fmt.Printf("result       %s\n", result.String())
	fmt.Printf("successor	 %s\n", dhtNode.acceleratedLookupUsingFingers(resultHex))
}


func (dhtNode *DHTNode) PrintRing() {
	fmt.Println(dhtNode.nodeId)
	
	/* There is more than one node */
	if dhtNode.successor != nil {
		dhtNode.successor.printRingAux(dhtNode.nodeId)
	}
}

func(dhtNode * DHTNode) printRingAux(nodeID string) {
	if dhtNode.nodeId != nodeID {
		
		/* This is not the first node */
		fmt.Println(dhtNode.nodeId)
		dhtNode.successor.printRingAux(nodeID)
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
