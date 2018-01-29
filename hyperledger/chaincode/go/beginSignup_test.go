package main

import (
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/stretchr/testify/assert"
)

func Test_BeginSignup(t *testing.T) {

	scc := new(SmartContract)
	stub := shim.NewMockStub("test_beginSignup", scc)

	stub.MockTransactionStart("t123")
	PutState(stub, "state", SETUP)
	stub.MockTransactionEnd("t123")

	checkInvoke(t, stub, [][]byte{[]byte("beginSignUp"), []byte("This is a question!")})

	var state StateEnum
	GetState(stub, "state", &state)
	assert.Equal(t, SIGNUP, state)

	var question string
	GetState(stub, "question", &question)
	assert.Equal(t, "This is a question!", question)
}
