package main

import (
	"fmt"
	"dht"
	"time"
	"strconv"
	"os"
	"runtime"
)

//Test function number 1 for the Objective 2
//It works with a SPACESIZE of 3 bits
func task21(){
	 	id1 := "01"
 	nodo1 := dht.MakeDHTNode(&id1,"localhost","1111")
 	nodo1.StartListenServer()
 	
 	go func() {
 		id2 := "02"
 		nodo2 := dht.MakeDHTNode(&id2,"localhost","1112")
 		nodo2.StartListenServer()
 		nodoAConectar := dht.NetworkNode{
 			NodeId: "01",
 			Ip: "localhost",
 			Port: "1111"}
 		nodo2.SendAddToRing(&nodoAConectar,nodo2.ToNetworkNode())
// 		nodo2.SendPrintRing(&nodoAConectar)
 	}()
 	
 	 	go func() {
 		id2 := "00"
 		nodo2 := dht.MakeDHTNode(&id2,"localhost","1113")
 		nodo2.StartListenServer()
 		nodoAConectar := dht.NetworkNode{
 			NodeId: "01",
 			Ip: "localhost",
 			Port: "1111"}
 		nodo2.SendAddToRing(&nodoAConectar,nodo2.ToNetworkNode())
// 		nodo2.SendPrintRing(&nodoAConectar)
 	}()
 	 	
		go func() {
		time.Sleep(1000 * time.Millisecond)
 		id2 := "07"
 		nodo2 := dht.MakeDHTNode(&id2,"localhost","1114")
 		nodo2.StartListenServer()
 		nodoAConectar := dht.NetworkNode{
 			NodeId: "02",
 			Ip: "localhost",
 			Port: "1112"}
 		nodo2.SendAddToRing(&nodoAConectar,nodo2.ToNetworkNode())
// 		nodo2.SendPrintRing(&nodoAConectar)
 	}()
 	 	 	
 	 	 	 	go func() {
 		id2 := "04"
 		nodo2 := dht.MakeDHTNode(&id2,"localhost","1115")
 		nodo2.StartListenServer()
 		nodoAConectar := dht.NetworkNode{
 			NodeId: "01",
 			Ip: "localhost",
 			Port: "1111"}
 		nodo2.SendAddToRing(&nodoAConectar,nodo2.ToNetworkNode())
// 		nodo2.SendPrintRing(&nodoAConectar)
 	}()
 	 	 	 	
  	 	go func() {
 		id2 := "05"
 		nodo2 := dht.MakeDHTNode(&id2,"localhost","1116")
 		nodo2.StartListenServer()
 		nodoAConectar := dht.NetworkNode{
 			NodeId: "01",
 			Ip: "localhost",
 			Port: "1111"}
 		nodo2.SendAddToRing(&nodoAConectar,nodo2.ToNetworkNode())
// 		nodo2.SendPrintRing(&nodoAConectar)
 	}()
 	 	 	 	 	
// 	 	 	 	 	 	go func() {
// 		id2 := "06"
// 		nodo2 := dht.MakeDHTNode(&id2,"localhost","1117")
// 		nodo2.StartListenServer()
// 		nodoAConectar := dht.NetworkNode{
// 			NodeId: "01",
// 			Ip: "localhost",
// 			Port: "1111"}
// 		nodo2.SendAddToRing(&nodoAConectar,nodo2.ToNetworkNode())
//// 		nodo2.SendPrintRing(&nodoAConectar)
// 	}()
 	nodo1.PrintRing()
 	time.Sleep(10000 * time.Millisecond)
 	nodo1.PrintRing()
 	time.Sleep(30000000000 * time.Millisecond)
}

//Test function number 2 for the Objective 2
//It works with a SPACESIZE of 160 bits
func task22(){
	node0b := dht.MakeDHTNode(nil, "localhost", "1111")
 	node0b.StartListenServer()
 	for i:=1112; i<1220; i++{
 		
 		time.Sleep(150 * time.Millisecond)
 		
	 	go func() { 
			fmt.Println("Sent routine "+strconv.Itoa(i))
	 		node1b := dht.MakeDHTNode(nil, "localhost", strconv.Itoa(i))
	 		node1b.StartListenServer()
	 		node0b.AddToRing(node1b.ToNetworkNode())
	 	}()
	 	
	 }
 	
	time.Sleep(5000 * time.Millisecond)
	node0b.PrintRing()
	time.Sleep(100000 * time.Millisecond)
	node0b.PrintRing()
	time.Sleep(20000 * time.Millisecond)
}


//Main function
func main() {
//	task31()
//	task22()
	task21()
}

func task31(){
	typeNode := os.Args[1]
	port := os.Args[2]
	if typeNode == "initial" {
		fmt.Println("Creando nodo inicial con puerto " + port)
		node0b := dht.MakeDHTNode(nil, "localhost", port)
 		node0b.StartListenServer()
 		for {
    		runtime.Gosched()
		}
	} else if typeNode == "connect" {
		portToConnect := os.Args[3]
		fmt.Println("Conectando nodo con puerto " + port + " al anillo " + portToConnect)
		node0b := dht.MakeDHTNode(nil, "localhost", port)
		node0b.StartListenServer()
		dht.SendAddToRingForeign("localhost",portToConnect,node0b.ToNetworkNode())
		for {
    		runtime.Gosched()
		}
	} else{
		portToConnect := os.Args[3]
		fmt.Println("Conectando nodo con puerto " + port + " al anillo " + portToConnect)
		node0b := dht.MakeDHTNode(nil, "localhost", port)
		node0b.StartListenServer()
		dht.SendAddToRingForeign("localhost",portToConnect,node0b.ToNetworkNode())
		time.Sleep(5000 * time.Millisecond)
		node0b.PrintRing()
		for {
			time.Sleep(10000 * time.Millisecond)
			node0b.PrintRing()
//    		runtime.Gosched()
		}
	}
}
