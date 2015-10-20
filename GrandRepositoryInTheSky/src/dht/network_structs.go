package dht

import (
)

//Network node struct
type NetworkNode struct {
	NodeId string
	Ip string
	Port string
}

//Message struct
type Msg struct {
	Source *NetworkNode
	Dest *NetworkNode
	Type string
	Args map[string]string
	Data DataSet
}

//Converts a DHTNode in a NetworkNode, that can be sent
func (dhtNode *DHTNode) ToNetworkNode() *NetworkNode {
	networkNode := new(NetworkNode)
	networkNode.NodeId = dhtNode.GetNodeId()
	networkNode.Ip = dhtNode.GetIp()
	networkNode.Port = dhtNode.GetPort()
	return networkNode
}