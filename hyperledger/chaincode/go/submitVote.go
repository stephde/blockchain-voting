package main

import (
	"fmt"

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

	txid := stub.GetTxID()
	userID := args[0]
	vote := args[1]

	// Make sure the sender can vote and hasn't already voted
	var registered map[string]struct{}
	GetState(stub, "registered", &registered)

	votecastCompositeIndex := "varName~userID~txID"
	name := "vote"

	deltaResultsIterator, deltaErr := stub.GetStateByPartialCompositeKey(votecastCompositeIndex, []string{"votecast", userID})
	if deltaErr != nil {
		return shim.Error(fmt.Sprintf("Could not retrieve value for %s: %s", name, deltaErr.Error()))
	}
	defer deltaResultsIterator.Close()

	_, found1 := registered[userID]
	hasVoted := deltaResultsIterator.HasNext()

	if !found1 || hasVoted {
		return shim.Error(userID + " is not allowed to vote")
	}

	// TODO: userID could be voting key and vote could be ZKP encrypted
	compositeIndexName := "varName~userID~vote~txID"
	compositeKey, compositeErr := stub.CreateCompositeKey(compositeIndexName, []string{name, userID, vote, txid})
	if compositeErr != nil {
		return shim.Error(fmt.Sprintf("Could not create a composite key for %s: %s", name, compositeErr.Error()))
	}

	// User is registered and did not cast vote yet
	compositePutErr := stub.PutState(compositeKey, []byte{0x00})
	if compositePutErr != nil {
		return shim.Error(fmt.Sprintf("Could not put operation for %s in the ledger: %s", name, compositePutErr.Error()))
	}

	// Saving votecast
	votecastCompositeKey, votecastCompositeErr := stub.CreateCompositeKey(votecastCompositeIndex, []string{"votecast", userID, txid})
	if votecastCompositeErr != nil {
		return shim.Error(fmt.Sprintf("Could not create a composite key for %s: %s", "votecast", votecastCompositeErr.Error()))
	}
	compositePutErr = stub.PutState(votecastCompositeKey, []byte{0x00})
	if compositePutErr != nil {
		return shim.Error(fmt.Sprintf("Could not put operation for %s in the ledger: %s", "votecast", compositePutErr.Error()))
	}

	return shim.Success(nil)
}
