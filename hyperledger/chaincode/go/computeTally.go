package main

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

func (s *SmartContract) computeTally(stub shim.ChaincodeStubInterface) sc.Response {

	if !s.inState(stub, VOTE) {
		return shim.Error("Wrong state")
	}

	curve := crypto.S256()

	var totalRegistered int
	GetState(stub, "totalregistered", &totalRegistered)

	var voters []Voter
	GetState(stub, "voters", &voters)

	// is a point on the curve, I guess
	tempResult := []*big.Int{big.NewInt(0), big.NewInt(0)}

	// Sum all votes
	for _, voter := range voters {
		voterAddress := voter.address
		fmt.Printf("key[%s] value[%s]\n", voterAddress, voter)

		var votecast map[string]bool
		GetState(stub, "votecast", &votecast)

		value, found := votecast[voterAddress]
		if found && !value {
			return shim.Error("Voter " + voterAddress + " has not voted")
		}

		vote := voter.vote
		logger.Info(vote)

		tempResult[0], tempResult[1] = curve.Add(tempResult[0], tempResult[1], vote[0], vote[1])

		if tempResult[0] == big.NewInt(0) {
			finalTally := []int{0, totalRegistered}
			finalTallyBytes, _ := json.Marshal(finalTally)
			return shim.Success(finalTallyBytes)
		} else {
			tempG := []*big.Int{curve.Params().Gx, curve.Params().Gy}
			for i := 1; i <= totalRegistered; i++ {
				if tempResult[0] == tempG[0] {
					finalTally := []int{i, totalRegistered}
					finalTallyBytes, _ := json.Marshal(finalTally)
					return shim.Success(finalTallyBytes)
				}
				tempG[0], tempG[1] = curve.Add(tempG[0], tempG[1], curve.Params().Gx, curve.Params().Gy)
			}
		}

		// Something bad happened, we should never get here
		finalTally := []int{0, 0}
		finalTallyBytes, _ := json.Marshal(finalTally)
		return shim.Success(finalTallyBytes)
	}

	// All votes have been accounted for...
	// Get tally and change state to 'Finished'
	s.transitionToState(stub, FINISHED)

	// todo
	return shim.Error("Not implemented yet")

}
