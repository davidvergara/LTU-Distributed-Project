package dht

import (
	"fmt"
	"dht"
)

func main() {
	id0 := "00"
	id1 := "01"
	id2 := "02"
	id3 := "03"
	id4 := "04"
	id5 := "05"
	id6 := "06"
	id7 := "07"

	node0b := dht.MakeDHTNode(&id0, "localhost", "1111")
	node1b := dht.MakeDHTNode(&id1, "localhost", "1112")
	node2b := dht.MakeDHTNode(&id2, "localhost", "1113")
	node3b := dht.MakeDHTNode(&id3, "localhost", "1114")
	node4b := dht.MakeDHTNode(&id4, "localhost", "1115")
	node5b := dht.MakeDHTNode(&id5, "localhost", "1116")
	node6b := dht.MakeDHTNode(&id6, "localhost", "1117")
	node7b := dht.MakeDHTNode(&id7, "localhost", "1118")

	node1b.AddToRing(node2b)
	node1b.AddToRing(node0b)
	node1b.AddToRing(node3b)
	node1b.AddToRing(node4b)
	node4b.AddToRing(node5b)
	node3b.AddToRing(node6b)
	node3b.AddToRing(node7b)

	fmt.Println("-> ring structure")
	node1b.PrintRing()
}
