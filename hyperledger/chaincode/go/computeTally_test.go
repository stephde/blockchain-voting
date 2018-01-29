package main

import (
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func Test_ComputeTally(t *testing.T) {

	scc := new(SmartContract)
	stub := shim.NewMockStub("test_computeTally", scc)

	stub.MockTransactionStart("t123")
	PutState(stub, "state", VOTE)
	stub.MockTransactionEnd("t123")

	checkInvoke(t, stub, [][]byte{[]byte("computeTally")})
}
