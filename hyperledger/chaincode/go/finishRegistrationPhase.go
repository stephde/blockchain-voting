package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

func (s *SmartContract) finishRegistrationPhase(stub shim.ChaincodeStubInterface) sc.Response {
	if !s.inState(stub, SIGNUP) {
		return shim.Error("Wrong state, expected SIGNUP")
	}

	name := "register"

	deltaResultsIterator, deltaErr := stub.GetStateByPartialCompositeKey("varName~userID~txID", []string{name})
	if deltaErr != nil {
		return shim.Error(fmt.Sprintf("Could not retrieve value for %s: %s", name, deltaErr.Error()))
	}
	defer deltaResultsIterator.Close()

	// Check the variable existed
	if !deltaResultsIterator.HasNext() {
		return shim.Error(fmt.Sprintf("No variable by the name %s exists", name))
	}

	registered := map[string]struct{}{}
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

		// Retrieve the userID
		userID := keyParts[1]
		registered[userID] = struct{}{}
	}

	PutState(stub, "totalRegistered", i)
	if i < 3 {
		// Legacy from Anonymous Voting Protocol
		return shim.Error("Too few voters registered, need at least 3")
	}

	PutState(stub, "registered", registered)

	// Now we either enter the voting phase.
	s.transitionToState(stub, VOTE)

	return shim.Success([]byte("Finished registration phase successfully"))
}
