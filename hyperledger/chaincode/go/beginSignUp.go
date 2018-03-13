package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

/*
 * Starts a vote by setting the question and Initializing the voting stage.
 */
func (s *SmartContract) beginSignUp(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if !s.inState(stub, SETUP) {
		return shim.Error("Wrong state, expected SETUP")
	}

	if len(args) != 1 {
		return shim.Error("Expected one argument: question")
	}

	question := args[0]
	PutState(stub, "question", question)
	s.transitionToState(stub, SIGNUP)

	return shim.Success(nil)
}
