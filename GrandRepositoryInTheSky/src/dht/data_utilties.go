package dht

import (
	"fmt"
)

//DataSet struct
type DataSet struct {
	DataStored map[string]Data
}

//Data struct
type Data struct {
	Value string
	Original bool
}


//Stores the Data{key,value,original} in the dataset
func (dataSet DataSet) StoreData(key string, value string, original bool) bool{
	_,is := dataSet.DataStored[key]
	
	if is {
		
		/* Data is already stored */
		return false
	} else{
		dataSet.DataStored[key] = Data{Value : value,
										Original : original}
		return true
	}
}

//Deletes the data with the key passed as parameter from the dataset
func (dataSet DataSet) deleteData(key string) bool{
	_,is := dataSet.DataStored[key]
	
	if is {
		
		/* Data is stored */
		delete(dataSet.DataStored,key)
		return true
	} else{
		
		/* Data is not stored */
		return false
	}
}

//Gets the data with the key passed as parameter from the dataset
func (dataSet DataSet) getData(key string) (Data, bool){
	data,is := dataSet.DataStored[key]
	
	if is {
		
		/* Data can be got */
		return data, true
	} else{
		
		/* Data is not stored */
		fmt.Println("Error: data with key \"" + key + "\"can't be found")
		return Data{"",false}, false
	}
}

//Change the data with the key passed as parameter from replica to original
func (dataSet DataSet) changeReplicaOriginal(key string){
	oldData,_ :=dataSet.getData(key)
	dataSet.deleteData(key) 
	dataSet.StoreData(key, oldData.Value, true)
}

//Change the data with the key passed as parameter from original to replica
func (dataSet DataSet) changeOriginalReplica(key string){
	oldData,_ :=dataSet.getData(key)
	dataSet.deleteData(key) 
	dataSet.StoreData(key, oldData.Value, false)
}

//Gets all the dataset
func (dataSet DataSet) getStoredData() map[string]Data{
	return dataSet.DataStored
}

func MakeDataSet() DataSet{
	return DataSet{make(map[string]Data)}
}