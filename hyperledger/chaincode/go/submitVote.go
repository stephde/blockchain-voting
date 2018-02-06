package main

import (
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

func (s *SmartContract) submitVote(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if !s.inState(stub, VOTE) {
		return shim.Error("Wrong state, expected VOTE")
	}

	if len(args) != 2 {
		return shim.Error("Expecting two arguments: UserID and Vote")
	}

	userID := args[0]
	_, err := strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("Vote was non-readable int")
	}

	// Make sure the sender can vote and hasn't already voted
	var registered map[string]bool
	var votecast map[string]bool
	GetState(stub, "registered", &registered)
	GetState(stub, "votecast", &votecast)

	value1, found1 := registered[userID]
	value2, found2 := votecast[userID]

	if !found1 || !found2 || !value1 || value2 {
		return shim.Error(userID + " is not allowed to vote")
	}

	// User is registered and did not cast vote yet
	// var voters map[string]Voter
	// GetState(stub, "voters", &voters)
	// voter := voters[userID]
	// voter.Vote = vote
	// voters[userID] = voter
	// PutState(stub, "voters", voters)

	votecast[userID] = true
	PutState(stub, "votecast", votecast)

	return shim.Success(nil)
}
