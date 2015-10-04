package dht

import (
	"strconv"
	"fmt"
)

type Contact struct {
	ip   string
	port string
}

type DHTNode struct {
	nodeId      string
	successor   *DHTNode
	predecessor *DHTNode
	contact     Contact
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
	} else {
		
		/* More than one node */
		valueNode,_ := strconv.Atoi(dhtNode.nodeId)
		valueNodeNew,_ := strconv.Atoi(newDHTNode.nodeId)
		valueNodeNext,_ := strconv.Atoi(dhtNode.successor.nodeId)
		
		/* Look if dhtNode is last node in the ring */
		if valueNode > valueNodeNext {
			
			/* New node between last and first nodes */
			if valueNodeNew > valueNode || valueNodeNew < valueNodeNext {
				oldSuccessorDhtNode := dhtNode.successor
				dhtNode.successor = newDHTNode
				newDHTNode.successor = oldSuccessorDhtNode
				newDHTNode.predecessor = dhtNode
				oldSuccessorDhtNode.predecessor = newDHTNode
			} else {
				
				/* New node is not after last node ->
			       recursion with first node */
				dhtNode.successor.AddToRing(newDHTNode)
			}
		} else {
			
			/* Trying to insert between 2 consecutive nodes */
			if valueNodeNew > valueNode && valueNodeNew < valueNodeNext {
				
				/* New node id bigger than dhtNode id and smaller than next node id ->
			       inserction between those nodes */
					oldSuccessorDhtNode := dhtNode.successor
					dhtNode.successor = newDHTNode
					newDHTNode.successor = oldSuccessorDhtNode
					newDHTNode.predecessor = dhtNode
					oldSuccessorDhtNode.predecessor = newDHTNode
			} else {
				
				/* New node is not between those nodes ->
			       recursion with next node */
				dhtNode.successor.AddToRing(newDHTNode)
			}	
		}
	}
}

func (dhtNode *DHTNode) Lookup(key string) *DHTNode {
	
	if dhtNode.nodeId == key {
		
		/* This is the node */
		return dhtNode
	} else {
		
		/* We look the successor */
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
	return dhtNode // XXX This is not correct obviously
}

func (dhtNode *DHTNode) responsible(key string) bool {
	// TODO
	return false
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
