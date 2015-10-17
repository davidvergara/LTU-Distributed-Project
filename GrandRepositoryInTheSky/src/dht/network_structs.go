package dht

import (
)

type NetworkNode struct {
	NodeId string
	Ip string
	Port string
}

type Msg struct {
	Source *NetworkNode
	Dest *NetworkNode
	Type string
	Args map[string]string
}

func (dhtNode *DHTNode) ToNetworkNode() *NetworkNode {
	networkNode := new(NetworkNode)
	networkNode.NodeId = dhtNode.GetNodeId()
	networkNode.Ip = dhtNode.GetIp()
	networkNode.Port = dhtNode.GetPort()
	return networkNode
}