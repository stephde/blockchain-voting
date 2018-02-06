package main

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

func (s *SmartContract) submitVoteInternal(stub shim.ChaincodeStubInterface,
	userID string,
	params [4]*big.Int,
	y *ecdsa.PublicKey,
	a1 *ecdsa.PublicKey,
	b1 *ecdsa.PublicKey,
	a2 *ecdsa.PublicKey,
	b2 *ecdsa.PublicKey) sc.Response {

	// logger.Info("Creator is ", creator)
	// logger.Info("Err is ", err)
	// logger.Info("UserId is ", userID)

	// Make sure the sender can vote and hasn't already voted
	var registered map[string]bool
	var votecast map[string]bool
	var voters []Voter
	GetState(stub, "registered", &registered)
	GetState(stub, "votecast", &votecast)
	GetState(stub, "voters", &voters)

	value1, found1 := registered[userID]
	value2, found2 := votecast[userID]

	if found1 && found2 && value1 && !value2 {
		for i := range voters {
			if voters[i].address == userID {
				v := voters[i]
				s.verify1outOf2ZKP(v, params, y, a1, b1, a2, b2)
				votecast[userID] = true
				v.vote = []*big.Int{y.X, y.Y}
				voters[i] = v

			}
		}

		PutState(stub, "voters", voters)
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
