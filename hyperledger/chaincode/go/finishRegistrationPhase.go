package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

func (s *SmartContract) finishRegistrationPhase(stub shim.ChaincodeStubInterface) sc.Response {
	if !s.inState(stub, SIGNUP) {
		return shim.Error("Wrong state")
	}

	var totalRegistered int
	GetState(stub, "totalRegistered", &totalRegistered)
	if totalRegistered < 3 {
		return shim.Error("Too few voters registered, need at least 3")
	}

	return shim.Error("Not implemented yet")
}
