package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

func (s *SmartContract) beginSignUp(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if !s.inState(stub, SETUP) {
		return shim.Error("Wrong state")
	}

	question := args[0]
	PutState(stub, "question", question)
	s.transitionToState(stub, SIGNUP)

	return shim.Success(nil)
}
