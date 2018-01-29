package main

import (
	"strconv"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func Test_SubmitVote(t *testing.T) {
	stub := shim.NewMockStub("test_submitVote", new(SmartContract))

	stub.MockTransactionStart("t123")
	PutState(stub, "state", VOTE)
	stub.MockTransactionEnd("t123")

	userID := "1"
	stub.MockTransactionStart("t124")
	registered := map[string]bool{userID: true}
	PutState(stub, "registered", registered)
	stub.MockTransactionEnd("t124")

	stub.MockTransactionStart("t125")
	votecast := map[string]bool{userID: false}
	PutState(stub, "votecast", votecast)
	stub.MockTransactionEnd("t125")

	checkInvoke(t, stub, [][]byte{[]byte("submitVote"), []byte(userID), []byte(strconv.Itoa(1))})
}