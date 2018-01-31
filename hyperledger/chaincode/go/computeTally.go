package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

func (s *SmartContract) computeTally(stub shim.ChaincodeStubInterface) sc.Response {

	if !s.inState(stub, VOTE) {
		return shim.Error("Wrong state, expected VOTE")
	}

	var totalRegistered int
	GetState(stub, "totalRegistered", &totalRegistered)

	// Initialize all results with 0
	tempResult := map[int]int{}

	compositeIndexName := "varName~userID~vote~txID"
	name := "vote"
	deltaResultsIterator, deltaErr := stub.GetStateByPartialCompositeKey(compositeIndexName, []string{name})
	if deltaErr != nil {
		return shim.Error(fmt.Sprintf("Could not retrieve value for %s: %s", name, deltaErr.Error()))
	}
	defer deltaResultsIterator.Close()
	// Check the variable existed
	if !deltaResultsIterator.HasNext() {
		return shim.Error(fmt.Sprintf("No variable by the name %s exists", name))
	}

	// Sum all votes
	var i int
	for i = 0; deltaResultsIterator.HasNext(); i++ {
		// Get the next row
		responseRange, nextErr := deltaResultsIterator.Next()
		if nextErr != nil {
			return shim.Error(nextErr.Error())
		}

		// Split the composite key into its component parts
		_, keyParts, splitKeyErr := stub.SplitCompositeKey(responseRange.Key)
		if splitKeyErr != nil {
			return shim.Error(splitKeyErr.Error())
		}

		// Retrieve the delta value and operation
		vote, _ := strconv.Atoi(keyParts[2])
		tempResult[vote]++
	}

	if i < totalRegistered {
		logger.Error("Someone did not vote")
	} else if i > totalRegistered {
		return shim.Error("Someone voted multiple times")
	}

	finalTally := Result{i, tempResult}

	// All votes have been accounted for...
	// Get tally and change state to 'Finished'
	s.transitionToState(stub, FINISHED)

	finalTallyBytes, _ := json.Marshal(finalTally)
	return shim.Success(finalTallyBytes)
}
