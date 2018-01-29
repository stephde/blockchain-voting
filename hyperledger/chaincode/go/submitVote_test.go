package main

import (
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func Test_SubmitVote(t *testing.T) {

	scc := new(SmartContract)
	stub := shim.NewMockStub("test_submitVote", scc)

	stub.MockTransactionStart("t123")
	PutState(stub, "state", VOTE)
	stub.MockTransactionEnd("t123")

	userId := "1"
	stub.MockTransactionStart("t124")
	registered := map[string]bool{userId: true}
	PutState(stub, "registered", registered)
	stub.MockTransactionEnd("t124")

	stub.MockTransactionStart("t125")
	votecast := map[string]bool{userId: false}
	PutState(stub, "votecast", votecast)
	stub.MockTransactionEnd("t125")

	checkInvoke(t, stub, [][]byte{[]byte("submitVote"), []byte(userId)})
}
