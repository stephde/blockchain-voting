package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

func (s *SmartContract) finishRegistrationPhase(stub shim.ChaincodeStubInterface) sc.Response {
	if !s.inState(stub, SIGNUP) {
		return shim.Error("Wrong state, expected SIGNUP")
	}

	var totalRegistered int
	GetState(stub, "totalRegistered", &totalRegistered)
	if totalRegistered < 3 {
		// Legacy from Anonymous Voting Protocol
		return shim.Error("Too few voters registered, need at least 3")
	}

	var registered map[string]bool
	GetState(stub, "registered", &registered)

	votecast := map[string]bool{}

	for voterUserID, registered := range registered {
		if registered {
			votecast[voterUserID] = false
		}
	}

	PutState(stub, "votecast", votecast)

	// Now we either enter the voting phase.
	s.transitionToState(stub, VOTE)

	return shim.Success([]byte("Success"))
}
