package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

func (s *SmartContract) initVote(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	// TODO: What are these two lines good for?
	// G[0] := Gx
	// G[1] := Gy

	PutState(stub, "state", SETUP)
	PutState(stub, "question", "No Question set")

	return shim.Success(nil)
}
