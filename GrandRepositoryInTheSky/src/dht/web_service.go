//Authors: Alejandro Márquez Ferrer / David Vergara Manrique
//Date: 28/10/2015
//
//Description: This file contains the code related to the web service which
//function is put/get/update/delete data in the chord ring.

package dht

import (
	"net/http"
	"fmt"
	//"strconv"
)

//The function intitializes the functions that are going
//to manage the post, get, put and delete petitions
//from the web page
func (node *DHTNode)InitializeWebServer(port string){
	go func() {

		http.HandleFunc("/chordnode/", InitializeHtml)

		http.HandleFunc("/chordnode/post/", node.Post)

		http.HandleFunc("/chordnode/get/", node.Get)
		
		http.HandleFunc("/chordnode/put/", node.Put)
		
		http.HandleFunc("/chordnode/delete/", node.Delete)

		http.ListenAndServe(":"+port, nil)
		fmt.Println("Web page up")
		

	}()
}

//Initialize the html code that will be shown in the 
//web page
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
				  "<input type=\"submit\" value=\"Submit\">" + 
				"</form>" + 
				
				"<p>Click on the submit button, and the input will be sent to the Chord network</p></body>")
}

//Post function
//Add a key and value to a node
//The key and the value are inserted by the user in the web page
//The function looks where they have to be insert and inserts them
//Finally shows a message saying how was the insertation
func (node *DHTNode)Post(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("key_value")
	value := r.FormValue("value_value")
	fmt.Fprintf(w, "<body><h3>1.- Upload a new key-value pair</h3>")
	fmt.Fprintf(w, "Key to save: %s </br>", key)
	fmt.Fprintf(w, "Value to save: %s </br></br>",value)
	fmt.Fprintf(w, "Trying to save the pair [key, value]...</br></br>")
	result, nodeResponsible :=node.HttpPost(key, value)
	if result{
		fmt.Fprintf(w, "<b>[Key, value] uploaded properly in node %s</b>", nodeResponsible.NodeId)
	}else{
		fmt.Fprintf(w, "<b>The upload failed</b>")
	}
	
	fmt.Fprintf(w, "</br>")
	fmt.Fprintf(w, "<p><a href=\"/chordnode/\">Back to main page</a></p>")
}

//Get function
//Gives the data that contains a key
//If the key does not exist shows an error
//The user insert the key in the web page
func (node *DHTNode)Get(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("key_value")
	fmt.Fprintf(w, "<body><h3>2.- Returns the value for a specific key</h3>")
	fmt.Fprintf(w, "Key to find: %s </br></br>", key)
	fmt.Fprintf(w, "Trying to find the key...</br></br>")
	
	result, dataValue, nodeResponsible:= node.HttpGet(key)
	
	if result{
		fmt.Fprintf(w, "<b>Data is saved in node %s</b></br>", nodeResponsible.NodeId)
		fmt.Fprintf(w, "<b>Value of the Data: %s</b>", dataValue)
	}else{
		fmt.Fprintf(w, "<b>The Get failed</b>")
	}
	
	fmt.Fprintf(w, "</br>")
	fmt.Fprintf(w, "<p><a href=\"/chordnode/\">Back to main page</a></p>")
}

//Put function
//Update the value of a key
//The key and the value are inserted by the user in the web page
//The function looks where is the key saved and update the data there
//If the key does not exist shows an error
//If all were correct shows a message of confirmation
func (node *DHTNode)Put(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("key_value")
	value := r.FormValue("value_value")
	fmt.Fprintf(w, "<body><h3>3.- Update the value for a specific key</h3>")
	fmt.Fprintf(w, "Key to update: %s </br>", key)
	fmt.Fprintf(w, "Value to update: %s </br></br>",value)
	fmt.Fprintf(w, "Trying to update the pair [key, value]...</br></br>")
	result, nodeResponsible :=node.HttpPut(key, value)
	if result{
		fmt.Fprintf(w, "<b>[Key, value] updated properly in node %s</b>", nodeResponsible.NodeId)
	}else{
		fmt.Fprintf(w, "<b>The update failed</b>")
	}
	
	fmt.Fprintf(w, "</br>")
	fmt.Fprintf(w, "<p><a href=\"/chordnode/\">Back to main page</a></p>")
}

//Delete function
//Delete the key and data of a key
//The key is inserted by the user in the web page
//The function looks where is the key saved and delete the key and the data
//If the key does not exist shows an error
//If all were correct shows a message of confirmation
func (node *DHTNode)Delete(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("key_value")
	fmt.Fprintf(w, "<body><h3>4.- Delete a key-value pair with key</h3>")
	fmt.Fprintf(w, "Key to delete: %s </br></br>", key)
	fmt.Fprintf(w, "Trying to delete the pair [key, value]...</br></br>")
	result, nodeResponsible :=node.HttpDelete(key)
	if result{
		fmt.Fprintf(w, "<b>[Key, value] deleted properly in node %s</b>", nodeResponsible.NodeId)
	}else{
		fmt.Fprintf(w, "<b>The delete failed</b>")
	}
	
	fmt.Fprintf(w, "</br>")
	fmt.Fprintf(w, "<p><a href=\"/chordnode/\">Back to main page</a></p>")
}