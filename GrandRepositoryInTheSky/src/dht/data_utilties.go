package dht

import (
	"fmt"
)

//DataSet struct
type DataSet struct {
	DataStored map[string]Data
}

type Data struct {
	Value string
	Original bool
}


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

func (dataSet DataSet) changeReplicaOriginal(key string){
	oldData,_ :=dataSet.getData(key)
	dataSet.deleteData(key) 
	dataSet.StoreData(key, oldData.Value, true)
}

func (dataSet DataSet) changeOriginalReplica(key string){
	oldData,_ :=dataSet.getData(key)
	dataSet.deleteData(key) 
	dataSet.StoreData(key, oldData.Value, false)
}

func (dataSet DataSet) getStoredData(key string) map[string]Data{
	return dataSet.DataStored
}

func MakeDataSet() DataSet{
	return DataSet{make(map[string]Data)}
}