package main

import (
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func Test_Question(t *testing.T) {
	scc := new(SmartContract)
	stub := shim.NewMockStub("test_question", scc)

	question := "What is the question?"
	stub.MockTransactionStart("t123")
	PutState(stub, "question", question)
	stub.MockTransactionEnd("t123")

	checkQuery(t, stub, "question", question)
}
