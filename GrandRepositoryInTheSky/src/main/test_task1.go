package main

import (
//	"fmt"
	"dht"
	"time"
//	"io/ioutil"
)
//
//func task1 () {
//	id0 := "00"
//	id1 := "01"
////	id2 := "02"
//	id3 := "03"
//	id4 := "03"
//	id5 := "05"
////	id6 := "06"
//	id7 := "07"
//
////	node0b := dht.MakeDHTNode(nil, "localhost", "1111")
////	node1b := dht.MakeDHTNode(nil, "localhost", "1112")
////	node2b := dht.MakeDHTNode(nil, "localhost", "1113")
////	node3b := dht.MakeDHTNode(nil, "localhost", "1114")
////	node4b := dht.MakeDHTNode(nil, "localhost", "1115")
//////	node5b := dht.MakeDHTNode(nil, "localhost", "1116")
////	node6b := dht.MakeDHTNode(nil, "localhost", "1117")
////	node7b := dht.MakeDHTNode(nil, "localhost", "1118")
////
//	node0b := dht.MakeDHTNode(&id0, "localhost", "1111")
//	node1b := dht.MakeDHTNode(&id1, "localhost", "1112")
////	node2b := dht.MakeDHTNode(&id2, "localhost", "1113")
//	node3b := dht.MakeDHTNode(&id3, "localhost", "1114")
//	node4b := dht.MakeDHTNode(&id4, "localhost", "1115")
//	node5b := dht.MakeDHTNode(&id5, "localhost", "1116")
////	node6b := dht.MakeDHTNode(&id6, "localhost", "1117")
//	node7b := dht.MakeDHTNode(&id7, "localhost", "1118")
//
////	node1b.AddToRing(node2b)
//	node1b.AddToRing(node0b)
//	node1b.AddToRing(node3b)
//	node1b.AddToRing(node4b)
//	node4b.AddToRing(node5b)
////	node3b.AddToRing(node6b)
//	node3b.AddToRing(node7b)
//
//	fmt.Println("-> ring structure")
//	node1b.PrintRing()
//	
//	fmt.Println()
//	
//	//nodeSearched := node1b.Lookup("08")
////	fmt.Print("Node searched (lookup) -> ")
//	//fmt.Println(nodeSearched.GetNodeId())
//}
//
//func task2() {
//	node0b := dht.MakeDHTNode(nil, "localhost", "1111")
//	
//	for i:=0; i<100; i++{
//		node1b := dht.MakeDHTNode(nil, "localhost", "1111")
//		node0b.AddToRing(node1b)
//		fmt.Printf("Node %d added \n",i)
//	}
//	
//	fmt.Println("-> ring structure")
//	node0b.PrintRing()
//}
//
func main() {
 	
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
 		id2 := "07"
 		nodo2 := dht.MakeDHTNode(&id2,"localhost","1114")
 		nodo2.StartListenServer()
 		nodoAConectar := dht.NetworkNode{
 			NodeId: "01",
 			Ip: "localhost",
 			Port: "1111"}
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
 	time.Sleep(3000 * time.Millisecond)
 	nodo1.PrintRing()
 	time.Sleep(30000000000 * time.Millisecond)

}
