package main

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

func (s *SmartContract) submitVoteInternal(stub shim.ChaincodeStubInterface,
	v Voter,
	params [4]*big.Int,
	y *ecdsa.PublicKey,
	a1 *ecdsa.PublicKey,
	b1 *ecdsa.PublicKey,
	a2 *ecdsa.PublicKey,
	b2 *ecdsa.PublicKey) sc.Response {
	// get sender address
	userID := v.address

	// logger.Info("Creator is ", creator)
	// logger.Info("Err is ", err)
	// logger.Info("UserId is ", userID)

	// Make sure the sender can vote and hasn't already voted
	var registered map[string]bool
	var votecast map[string]bool
	GetState(stub, "registered", &registered)
	GetState(stub, "votecast", &votecast)

	value1, found1 := registered[userID]
	value2, found2 := votecast[userID]

	// logger.Info("User is registered? ", found1)
	// logger.Info("User has voted? ", value2)

	if found1 && found2 && value1 && !value2 {
		// User is registered and did not cast vote yet
		// logger.Info("User is allowed to vote")

		// if s.verify1outOf2ZKP(v, params, y, a1, b1, a2, b2) {
		votecast[userID] = true
		v.vote = []*big.Int{y.X, y.Y}

		PutState(stub, "votecast", votecast)
		return shim.Success(nil)
		// } else {
		// return shim.Error("Verirfy1outOf2ZKP failed")
		// }

	} else {
		return shim.Error("User is not allowed to vote")
	}
}

func (s *SmartContract) submitVote(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if !s.inState(stub, VOTE) {
		return shim.Error("Wrong state")
	}

	return shim.Error("Not implemented yet")
	// return s.submitVoteInternal(stub, v, params, y, a1, b1, a2, b2)
}
