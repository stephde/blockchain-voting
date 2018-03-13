package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

/*
 * This function casts a vote that was submitted to Hyperledger.
 * Expects userID and vote (not encrypted in this version)
 */
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

	/*
	 * Make sure the sender can vote and hasn't already voted.
	 * Check is done by querying two CompositeKeys
	 */

	registerName := "register"
	registerResultsIterator, registerErr := stub.GetStateByPartialCompositeKey("varName~userID~txID", []string{registerName, userID})
	if registerErr != nil {
		return shim.Error(fmt.Sprintf("Could not retrieve value for %s: %s", registerName, registerErr.Error()))
	}
	defer registerResultsIterator.Close()

	votecastCompositeIndex := "varName~userID~txID"
	votecastName := "votecast"
	deltaResultsIterator, deltaErr := stub.GetStateByPartialCompositeKey(votecastCompositeIndex, []string{votecastName, userID})
	if deltaErr != nil {
		return shim.Error(fmt.Sprintf("Could not retrieve value for %s: %s", votecastName, deltaErr.Error()))
	}
	defer deltaResultsIterator.Close()

	isRegisterd := registerResultsIterator.HasNext()
	hasVoted := deltaResultsIterator.HasNext()

	if !isRegisterd || hasVoted {
		if !isRegisterd {
			return shim.Error(userID + " is not allowed to vote - not registered")
		}
		return shim.Error(userID + " is not allowed to vote - already voted")
	}

	/*
	 * For the crypto version:
	 * userID would be personal voting key
	 * vote would be encrypted zero-knowledge proof
	 */
	compositeIndexName := "varName~userID~vote~txID"
	name := "vote"
	compositeKey, compositeErr := stub.CreateCompositeKey(compositeIndexName, []string{name, userID, vote, txid})
	if compositeErr != nil {
		return shim.Error(fmt.Sprintf("Could not create a composite key for %s: %s", name, compositeErr.Error()))
	}

	// User is registered and did not cast vote yet
	compositePutErr := stub.PutState(compositeKey, []byte{0x00})
	if compositePutErr != nil {
		return shim.Error(fmt.Sprintf("Could not put operation for %s in the ledger: %s", name, compositePutErr.Error()))
	}

	// Vote has been stored, now save votecast for this user
	votecastCompositeKey, votecastCompositeErr := stub.CreateCompositeKey(votecastCompositeIndex, []string{votecastName, userID, txid})
	if votecastCompositeErr != nil {
		return shim.Error(fmt.Sprintf("Could not create a composite key for %s: %s", votecastName, votecastCompositeErr.Error()))
	}
	compositePutErr = stub.PutState(votecastCompositeKey, []byte{0x00})
	if compositePutErr != nil {
		return shim.Error(fmt.Sprintf("Could not put operation for %s in the ledger: %s", votecastName, compositePutErr.Error()))
	}

	return shim.Success(nil)
}
