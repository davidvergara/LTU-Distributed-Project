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
	Node string
}


func (dataSet DataSet) storeData(key string, value string, nodeId string) bool{
	_,is := dataSet.DataStored[key]
	
	if is {
		
		/* Data is already stored */
		return false
	} else{
		dataSet.DataStored[key] = Data{Value : value,
										Node : nodeId}
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

func (dataSet DataSet) getData(key string) Data{
	data,is := dataSet.DataStored[key]
	
	if is {
		
		/* Data can be got */
		return data
	} else{
		
		/* Data is not stored */
		fmt.Println("Error: data with key \"" + key + "\"can't be found")
		return Data{"","ERROR"}
	}
}

func (dataSet DataSet) getStoredData(key string) map[string]Data{
	return dataSet.DataStored
}

func MakeDataSet() DataSet{
	return DataSet{make(map[string]Data)}
}