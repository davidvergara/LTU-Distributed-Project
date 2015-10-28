//Authors: Alejandro Márquez Ferrer / David Vergara Manrique
//Date: 28/10/2015
//
//Description: This file contains the code used to make tests of the different parts
//of the project. To test them, you have to uncomment the one you want to execute,
//and comment the rest.

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
//	task32()
	task31()
//	task22()
//	task21()
}

func task31(){
	
	//Replication test (to be launched in the terminal)
	typeNode := os.Args[1]
	port := os.Args[2]
	
	if typeNode == "initial" {
		
		//Args[1] = type of request
		//Args[2] = port of the first node
		//Args[3] = id of the node
		id := os.Args[3]
		fmt.Println("Creando nodo inicial con puerto " + port)
		node0b := dht.MakeDHTNode(&id, "localhost", port)
		node0b.InitializeWebServer(port)
 		node0b.StartListenServer()
 		
 		for {
    		runtime.Gosched()
		}
	} else if typeNode == "connect" {
		
		//Args[1] = type of request
		//Args[2] = port of the node
		//Args[3] = port to connect to
		//Args[4] = id of the new node
		portToConnect := os.Args[3]
		id := os.Args[4]
		fmt.Println("Connecting node " + port + " to node " + portToConnect)
		node0b := dht.MakeDHTNode(&id, "localhost", port)
		node0b.InitializeWebServer(port)
		node0b.StartListenServer()
		dht.SendAddToRingForeign("localhost",portToConnect,node0b.ToNetworkNode())
		for {
    		runtime.Gosched()
		}
	} else if typeNode == "addData" {
		
		//Args[1] = type of request
		//Args[2] = port of the node to connect
		//Args[3] = data to send
		fmt.Println("Sending data to node " + port)
		
		dataSetToBeSend :=dht.MakeDataSet()
		dataSetToBeSend.StoreData(os.Args[3],"0",true)
		dht.SendDataToRingForeign("localhost", port,dataSetToBeSend)
		time.Sleep(1000 * time.Millisecond)
	} else if typeNode == "deleteData" {
		
		//Args[1] = type of request
		//Args[2] = port of the node to connect
		//Args[3] = data to delete
		fmt.Println("Sending delete data request to node " + port)
		dataSetToBeSend :=dht.MakeDataSet()
		dataSetToBeSend.StoreData(os.Args[3],"0",true)
		dht.SendDeleteDataForeign("localhost", port,dataSetToBeSend)
		time.Sleep(1000 * time.Millisecond)
	}else {
		
		//Args[1] = type of request
		//Args[2] = port of the node to send print ring request
		fmt.Println("Sending Print Ring request to node " + port)
		dht.SendPrintRingForeign("localhost", port)
		time.Sleep(1000 * time.Millisecond)
	}
}

func task32(){
	
	//Replication test (to be launched in the terminal 160bits)
	typeNode := os.Args[1]
	port := os.Args[2]
	
	if typeNode == "initial" {
		
		//Args[1] = type of request
		//Args[2] = port of the first node
		node0b := dht.MakeDHTNode(nil, "localhost", port)
		fmt.Println("Creando nodo inicial " + node0b.GetNodeId() + ":" + port)
		node0b.InitializeWebServer(port)
 		node0b.StartListenServer()
 		
 		for {
    		runtime.Gosched()
		}
	} else if typeNode == "connect" {
		
		//Args[1] = type of request
		//Args[2] = port of the node
		//Args[3] = port to connect to
		portToConnect := os.Args[3]
		node0b := dht.MakeDHTNode(nil, "localhost", port)
		fmt.Println("Connecting node "+ node0b.GetNodeId() + ":" + port + " to port " + portToConnect)
		node0b.StartListenServer()
		dht.SendAddToRingForeign("localhost",portToConnect,node0b.ToNetworkNode())
		for {
    		runtime.Gosched()
		}
	} else if typeNode == "addData" {
		
		//Args[1] = type of request
		//Args[2] = port of the node to connect
		//Args[3] = data to send
		fmt.Println("Sending data to node " + port)
		
		dataSetToBeSend :=dht.MakeDataSet()
		dataSetToBeSend.StoreData(os.Args[3],"0",true)
		dht.SendDataToRingForeign("localhost", port,dataSetToBeSend)
		time.Sleep(1000 * time.Millisecond)
	} else if typeNode == "deleteData" {
		
		//Args[1] = type of request
		//Args[2] = port of the node to connect
		//Args[3] = data to delete
		fmt.Println("Sending delete data request to node " + port)
		dataSetToBeSend :=dht.MakeDataSet()
		dataSetToBeSend.StoreData(os.Args[3],"0",true)
		dht.SendDeleteDataForeign("localhost", port,dataSetToBeSend)
		time.Sleep(1000 * time.Millisecond)
	}else {
		
		//Args[1] = type of request
		//Args[2] = port of the node to send print ring request
		fmt.Println("Sending Print Ring request to node " + port)
		dht.SendPrintRingForeign("localhost", port)
		time.Sleep(1000 * time.Millisecond)
	}
}
