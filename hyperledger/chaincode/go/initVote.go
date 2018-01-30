package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

func (s *SmartContract) initVote(stub shim.ChaincodeStubInterface) sc.Response {
	PutState(stub, "state", SETUP)
	PutState(stub, "eligible", map[string]bool{})
	PutState(stub, "registered", map[string]bool{})
	PutState(stub, "totalRegistered", 0)
	return shim.Success(nil)
}
