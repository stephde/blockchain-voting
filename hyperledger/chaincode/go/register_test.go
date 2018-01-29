package main

import (
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/stretchr/testify/assert"
)

func Test_Register(t *testing.T) {

	scc := new(SmartContract)
	stub := shim.NewMockStub("test_register", scc)

	userID := "userId"

	stub.MockTransactionStart("t123")
	PutState(stub, "state", SIGNUP)
	eligible := map[string]bool{userID: true}
	PutState(stub, "eligible", eligible)
	registered := map[string]bool{userID: false}
	PutState(stub, "registered", registered)
	stub.MockTransactionEnd("t124")

	checkInvoke(t, stub, [][]byte{
		[]byte("register"),
		[]byte(userID),
	})

	var totalRegistered int
	GetState(stub, "totalRegistered", &totalRegistered)
	assert.Equal(t, 1, totalRegistered)

	GetState(stub, "registered", &registered)
	assert.True(t, registered[userID])

	var voters []Voter
	GetState(stub, "voters", &voters)
	assert.Equal(t, 1, len(voters))
	assert.Equal(t, userID, voters[0].UserId)
}
