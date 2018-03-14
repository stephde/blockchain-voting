package main

import (
	"encoding/json"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func Ignore_ComputeTally(t *testing.T) {
	stub := shim.NewMockStub("test_computeTally", new(SmartContract))

	stub.MockTransactionStart("t123")
	PutState(stub, "state", VOTE)
	PutState(stub, "totalRegistered", 3)

	//voters := []Voter{}
	var i int
	for i = 0; i < 3; i++ {

	}

	stub.MockTransactionEnd("t123")

	result := Result{2, map[int]int{0: 1, 1: 1}}
	bytes, _ := json.Marshal(result)

	checkQuery(t, stub, "computeTally", string(bytes))
}
