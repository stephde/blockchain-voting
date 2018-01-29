package main

import (
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/stretchr/testify/assert"
)

func Test_InitVote(t *testing.T) {

	scc := new(SmartContract)
	stub := shim.NewMockStub("test_beginSignup", scc)

	checkInvoke(t, stub, [][]byte{[]byte("initVote")})

	var state StateEnum
	GetState(stub, "state", &state)
	assert.Equal(t, SETUP, state)
}
