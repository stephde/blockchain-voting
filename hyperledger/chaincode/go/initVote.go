package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

/*
 * Initializes a new vote.
 * Deletes state that belongs to a previous vote.
 */
func (s *SmartContract) initVote(stub shim.ChaincodeStubInterface) sc.Response {
	PutState(stub, "state", SETUP)
	PutState(stub, "eligible", map[string]bool{})
	PutState(stub, "registered", map[string]bool{})
	PutState(stub, "totalRegistered", 0)

	s.deleteCompositeKey(stub, "varName~userID~vote~txID", "vote")
	s.deleteCompositeKey(stub, "varName~userID~txID", "register")
	s.deleteCompositeKey(stub, "varName~userID~txID", "votecast")

	return shim.Success(nil)
}

/*
 * Helper function to delete a composite key from Hyperledger
 */
func (s *SmartContract) deleteCompositeKey(stub shim.ChaincodeStubInterface, compositeKey string, name string) sc.Response {
	deltaResultsIterator, deltaErr := stub.GetStateByPartialCompositeKey(compositeKey, []string{name})

	if deltaErr != nil {
		return shim.Error(fmt.Sprintf("Could not retrieve delta rows for %s: %s", name, deltaErr.Error()))
	}
	defer deltaResultsIterator.Close()

	// Ensure the variable exists
	if !deltaResultsIterator.HasNext() {
		return shim.Success([]byte(fmt.Sprintf("Variable %s already deleted", name)))
	}

	// Iterate through result set and delete all indices
	var i int
	for i = 0; deltaResultsIterator.HasNext(); i++ {
		responseRange, nextErr := deltaResultsIterator.Next()
		if nextErr != nil {
			return shim.Error(fmt.Sprintf("Could not retrieve next delta row: %s", nextErr.Error()))
		}

		deltaRowDelErr := stub.DelState(responseRange.Key)
		if deltaRowDelErr != nil {
			return shim.Error(fmt.Sprintf("Could not delete delta row: %s", deltaRowDelErr.Error()))
		}
	}

	return shim.Success([]byte(fmt.Sprintf("Deleted %s, %d rows removed", name, i)))
}
