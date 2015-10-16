package dht

import (
	"sync"
)

/* Network variables */
var LocalIp 		string
var LocalPort 		string
var NumLookup = 0
var LookupRequest 	map[int]chan *NetworkNode
var mutexNumLookup = &sync.Mutex{}
var mutexPredeccessor = &sync.Mutex{}
var mutexSuccessor = &sync.Mutex{} 

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



