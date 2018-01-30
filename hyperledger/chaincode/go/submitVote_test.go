package main

import (
	"strconv"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/stretchr/testify/assert"
)

func Test_SubmitVote(t *testing.T) {
	stub := shim.NewMockStub("test_submitVote", new(SmartContract))

	userID := "1"
	stub.MockTransactionStart("t123")
	PutState(stub, "state", VOTE)

	registered := map[string]bool{userID: true}
	PutState(stub, "registered", registered)

	checkInvoke(t, stub, [][]byte{[]byte("submitVote"), []byte(userID), []byte(strconv.Itoa(1))})

	name := "vote"
	deltaResultsIterator, _ := stub.GetStateByPartialCompositeKey("varName~userID~vote~txID", []string{name, userID})
	assert.True(t, deltaResultsIterator.HasNext())
}

func Test_InvalidUserID(t *testing.T) {
	stub := shim.NewMockStub("test_submitVote", new(SmartContract))

	stub.MockTransactionStart("t123")
	PutState(stub, "state", VOTE)
	stub.MockTransactionEnd("t123")

	checkFailingInvoke(t, stub, [][]byte{[]byte("submitVote"), []byte("invalidUser"), []byte(strconv.Itoa(1))})
}

func Test_DuplicateVote(t *testing.T) {
	stub := shim.NewMockStub("test_submitVote", new(SmartContract))

	userID := "1"
	stub.MockTransactionStart("t123")
	PutState(stub, "state", VOTE)

	registered := map[string]bool{userID: true}
	PutState(stub, "registered", registered)

	votecastCompositeKey, _ := stub.CreateCompositeKey("varName~userID~txID", []string{"votecast", userID, "1"})
	PutState(stub, votecastCompositeKey, []byte{0x00})
	stub.MockTransactionEnd("t123")

	checkFailingInvoke(t, stub, [][]byte{[]byte("submitVote"), []byte(userID), []byte(strconv.Itoa(1))})
}
