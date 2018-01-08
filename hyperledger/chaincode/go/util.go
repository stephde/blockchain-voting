package main

import (
	"encoding/json"
	"errors"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func GetState(stub shim.ChaincodeStubInterface, name string, v interface{}) error {
	bytes, err := stub.GetState(name)
	if err != nil {
		return errors.New("Failed to get state")
	}
	if bytes == nil {
		return errors.New("Entity not found")
	}

	return json.Unmarshal(bytes, &v)
}

func PutState(stub shim.ChaincodeStubInterface, name string, v interface{}) error {
	bytes, err := json.Marshal(v)

	if err != nil {
		return errors.New("Failed to put state")
	}
	if bytes == nil {
		return errors.New("Entity not found")
	}

	return stub.PutState(name, bytes)
}
