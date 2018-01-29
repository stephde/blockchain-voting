package main

import (
	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

func (s *SmartContract) computeTally(stub shim.ChaincodeStubInterface) sc.Response {

	if !s.inState(stub, VOTE) {
		return shim.Error("Wrong state, expected VOTE")
	}

	var totalRegistered int
	GetState(stub, "totalRegistered", &totalRegistered)

	var voters map[string]Voter
	GetState(stub, "voters", &voters)

	// Initialize all results with 0
	tempResult := map[int]int{}

	// Sum all votes
	for voterAddress, voter := range voters {
		var votecast map[string]bool
		GetState(stub, "votecast", &votecast)

		value, found := votecast[voterAddress]
		if found && !value {
			return shim.Error("Voter " + voterAddress + " has not voted")
		}

		tempResult[voter.Vote]++
	}

	finalTally := Result{totalRegistered, tempResult}

	// All votes have been accounted for...
	// Get tally and change state to 'Finished'
	s.transitionToState(stub, FINISHED)

	finalTallyBytes, _ := json.Marshal(finalTally)
	return shim.Success(finalTallyBytes)
}
