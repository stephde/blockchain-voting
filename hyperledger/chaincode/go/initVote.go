package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

func (s *SmartContract) initVote(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	PutState(stub, "state", SETUP)
	PutState(stub, "question", "No Question set")

	return shim.Success(nil)
}
