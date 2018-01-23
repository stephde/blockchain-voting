package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

func (s *SmartContract) submitVote(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if !s.inState(stub, VOTE) {
		return shim.Error("Wrong state")
	}

	// get sender address
	creator, err := stub.GetCreator()
	userID := args[0]

	logger.Info("Creator is ", creator)
	logger.Info("Err is ", err)
	logger.Info("UserId is ", userID)

	// Make sure the sender can vote and hasn't already voted
	var registered map[string]bool
	var votecast map[string]bool
	GetState(stub, "registered", &registered)
	GetState(stub, "votecast", &votecast)

	value1, found1 := registered[userID]
	value2, found2 := votecast[userID]

	logger.Info("User is registered? ", found1)
	logger.Info("User has voted? ", value2)

	if found1 && found2 && value1 && !value2 {
		// User is registered and did not cast vote yet
		logger.Info("User is allowed to vote")
	} else {
		return shim.Error("User is not allowed to vote")
	}

	return shim.Error("Not implemented yet")
}
