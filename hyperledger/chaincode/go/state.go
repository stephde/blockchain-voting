package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

/*
 * The following variables and functions are used to store the state of the vote.
 */
type StateEnum int

const (
	UNSET StateEnum = iota
	SETUP
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

/*
 * Verifies that vote is in a specific state
 */
func (s *SmartContract) inState(stub shim.ChaincodeStubInterface, expectedState StateEnum) bool {
	var state StateEnum
	GetState(stub, "state", &state)
	success := expectedState == state
	if !success {
		logger.Error("Expected state " + expectedState.String() + ", but was " + state.String())
	}

	return success
}

/*
 * At the end of a voting phase we might need to transition to the next state.
 */
func (s *SmartContract) transitionToState(stub shim.ChaincodeStubInterface, newState StateEnum) {
	PutState(stub, "state", newState)
}
