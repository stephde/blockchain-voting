package main

import (
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/stretchr/testify/assert"
)

func Test_SetEligible(t *testing.T) {
	stub := shim.NewMockStub("test_setEligible", new(SmartContract))

	stub.MockTransactionStart("t123")
	PutState(stub, "state", SETUP)
	stub.MockTransactionEnd("t123")

	checkInvoke(t, stub, [][]byte{
		[]byte("setEligible"),
		[]byte("a"),
		[]byte("b"),
		[]byte("c"),
	})

	var eligible map[string]bool
	var totalEligible int
	GetState(stub, "eligible", &eligible)
	GetState(stub, "totalEligible", &totalEligible)

	assert.Equal(t, 3, totalEligible)

	assert.True(t, eligible["a"])
	assert.True(t, eligible["b"])
	assert.True(t, eligible["c"])
}
