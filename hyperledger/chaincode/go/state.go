package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

/*
 * The following variables and functions are used to store the state of the vote.
 */
type StateEnum int

const (
	SETUP StateEnum = iota
	SIGNUP
	VOTE
	FINISHED
)

var states = [...]string{
	"SETUP",
	"SIGNUP",
	"VOTE",
	"FINISHED",
}

func (state StateEnum) String() string {
	return states[state]
}

func (s *SmartContract) inState(stub shim.ChaincodeStubInterface, expectedState StateEnum) bool {
	var state StateEnum
	GetState(stub, "state", &state)

	logger.Info("State is " + state.String())
	return expectedState == state
}

func (s *SmartContract) transitionToState(stub shim.ChaincodeStubInterface, newState StateEnum) {
	PutState(stub, "state", newState)
}
