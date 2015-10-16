package dht

import (

)

/* Network variables */
var LocalIp string
var LocalPort string

type NetworkNode struct {
	NodeId string
	Ip string
	Port string
}

type Msg struct {
	Source *NetworkNode
	Dest *NetworkNode
	Type string
//	Args map[string]string
}



