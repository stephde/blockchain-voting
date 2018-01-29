package main

import (
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

func (s *SmartContract) submitVote(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if !s.inState(stub, VOTE) {
		return shim.Error("Wrong state")
	}

	if len(args) != 2 {
		return shim.Error("Expecting two arguments: UserID and Vote")
	}

	userID := args[0]
	vote, err := strconv.Atoi(args[1])
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

	logger.Debug("User is registered? ", found1)
	logger.Debug("User has voted? ", value2)

	if found1 && found2 && value1 && !value2 {
		// User is registered and did not cast vote yet
		var voters map[string]Voter
		GetState(stub, "voters", &voters)
		voter := voters[userID]
		voter.Vote = vote
		PutState(stub, "voters", voters)

		return shim.Success(nil)
	} else {
		return shim.Error("User is not allowed to vote")
	}
}
