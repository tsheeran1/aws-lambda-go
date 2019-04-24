package main

import "encoding/json"

//FunctionList is a struct to contain the json function list returned by aws list-functions
type FunctionList struct {
	Functions []Function
}

// Function is a struct to contain an individual function
type Function struct {
	FunctionName string
}

//NewFunctionList is a function to retrieve function list from aws
func NewFunctionList() (FunctionList, error) {
	data, err := run("aws", "lambda", "list-functions")
	if err != nil {
		return FunctionList{}, err
	}

	var res FunctionList
	err = json.Unmarshal(data, &res)
	return res, err
}

//HasFunction is a function to determine if a function name (fname) exists in aws
func (fl FunctionList) HasFunction(fname string) bool {
	for _, v := range fl.Functions {
		if v.FunctionName == fname {
			return true
		}
	}
	return false

}
