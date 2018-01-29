package main

import (
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/stretchr/testify/assert"
)

func Test_WrongStateFinishRegistrationPhase(t *testing.T) {
	stub := shim.NewMockStub("Test_WrongState", new(SmartContract))

	checkFailingInvoke(t, stub, [][]byte{[]byte("finishRegistrationPhase")})
}

func Test_FinishRegistrationPhase(t *testing.T) {
	stub := shim.NewMockStub("test_beginVote", new(SmartContract))

	stub.MockTransactionStart("t123")
	PutState(stub, "state", SIGNUP)
	PutState(stub, "totalRegistered", 3)
	stub.MockTransactionEnd("t123")

	checkInvoke(t, stub, [][]byte{[]byte("finishRegistrationPhase")})

	var state StateEnum
	GetState(stub, "state", &state)
	assert.Equal(t, VOTE, state)
}
