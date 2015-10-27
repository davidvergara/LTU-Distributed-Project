
package dht

import (
	"net/http"
	"fmt"
	//"strconv"
)

func (node *DHTNode)InitializeWebServer(port string){
	go func() {

		http.HandleFunc("/chordnode/", InitializeHtml)

		http.HandleFunc("/chordnode/post/", node.Post)

//		http.HandleFunc("/chordnode/post/", node.Get)
//		
//		http.HandleFunc("/chordnode/post/", node.Put)
//		
//		http.HandleFunc("/chordnode/post/", node.Delete)

		http.ListenAndServe("192.168.1.4:"+port, nil)
		fmt.Println("Web page up")
		

	}()
}

func InitializeHtml(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<body><h3>1.- Upload a new key-value pair</h3>" + 
				"<form action=\"/chordnode/post/\" method=\"post\" >" + 
				  "Key:   <input type=\"text\" name=\"key_value\"><br>" + 
				  "Value: <input type=\"text\" name=\"value_value\"><br>" + 
				  "<input type=\"submit\" value=\"Submit\">" + 
				"</form>" + 
				
				"<h3>2.- Returns the value for a specific key</h3>" + 
				"<form action=\"/chordnode/get/\" method=\"post\" >" + 
				  "Key:    <input type=\"text\" name=\"key_value\"><br>" + 
				  "<input type=\"submit\" value=\"Submit\">" + 
				"</form>" + 
				
				"<h3>3.- Update the value for a specific key</h3>" + 
				"<form action=\"/chordnode/put/\" method=\"post\" >" + 
				  "Key:   <input type=\"text\" name=\"key_value\"><br>" + 
				  "Value: <input type=\"text\" name=\"value_value\"><br>" + 
				  "<input type=\"submit\" value=\"Submit\">" + 
				"</form>" + 
				
				"<h3>4.- Delete a key-value pair with key</h3>" + 
				"<form action=\"/chordnode/delete/\" method=\"post\" >" + 
				  "Key:   <input type=\"text\" name=\"key_value\"><br>" + 
				  "Value: <input type=\"text\" name=\"value_value\"><br>" + 
				  "<input type=\"submit\" value=\"Submit\">" + 
				"</form>" + 
				
				"<p>Click on the submit button, and the input will be sent to the Chord network</p></body>")
}

func (node *DHTNode)Post(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("key_value")
	value := r.FormValue("value_value")
	fmt.Fprintf(w, "<body><h3>1.- Upload a new key-value pair</h3>")
fmt.Fprintf(w, "Key to save: %s </br>", key)
fmt.Fprintf(w, "Value to save: %s </br>",value)
	
	if node.HttpPost(key, value){
		fmt.Fprintf(w, "Key, value uploaded property")
	}else{
		fmt.Fprintf(w, "The upload failed")
	}
	
	fmt.Fprintf(w, "</br>")
	fmt.Fprintf(w, "<p><a href=\"/chordnode/\">Back to main page</a></p>")
	
}