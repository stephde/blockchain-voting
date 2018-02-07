package main

import (
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/stretchr/testify/assert"
)

func Test_Register(t *testing.T) {
	stub := shim.NewMockStub("test_register", new(SmartContract))

	userID := "userId"

	stub.MockTransactionStart("t123")
	PutState(stub, "state", SIGNUP)
	eligible := map[string]bool{userID: true}
	PutState(stub, "eligible", eligible)
	stub.MockTransactionEnd("t123")

	checkInvoke(t, stub, [][]byte{
		[]byte("register"),
		[]byte(userID),
	})

	name := "register"
	deltaResultsIterator, _ := stub.GetStateByPartialCompositeKey("varName~userID~txID", []string{name})
	assert.True(t, deltaResultsIterator.HasNext())
}
