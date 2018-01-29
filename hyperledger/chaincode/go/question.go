package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

func (s *SmartContract) question(stub shim.ChaincodeStubInterface) sc.Response {
	var question string
	GetState(stub, "question", &question)
	return shim.Success([]byte(question))
}
