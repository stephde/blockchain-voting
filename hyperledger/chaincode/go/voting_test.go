package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func checkInit(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInit("1", args)
	if res.Status != shim.OK {
		fmt.Println("Init failed", string(res.Message))
		t.FailNow()
	}
}

func checkState(t *testing.T, stub *shim.MockStub, name string, value string) {
	bytes := stub.State[name]
	if bytes == nil {
		fmt.Println("State", name, "failed to get value")
		t.FailNow()
	}
	if string(bytes) != value {
		fmt.Println("State value", name, "was not", value, "as expected, but was", string(bytes))
		t.FailNow()
	}
}

func checkQuery(t *testing.T, stub *shim.MockStub, function string, name string, value string) {
	res := stub.MockInvoke("1", [][]byte{[]byte(function), []byte(name)})
	if res.Status != shim.OK {
		fmt.Println("Query", function, "with", name, "failed", string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		fmt.Println("Query", name, "failed to get value")
		t.FailNow()
	}
	if string(res.Payload) != value {
		fmt.Println("Query value", function, "was not", value, "as expected, but was", string(res.Payload))
		t.FailNow()
	}
}

func checkInvoke(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInvoke("1", args)
	if res.Status != shim.OK {
		fmt.Println("Invoke", args, "failed", string(res.Message))
		t.FailNow()
	}
}

// func Test_InitVote(t *testing.T) {
// 	scc := new(SmartContract)
// 	stub := shim.NewMockStub("test_initVote", scc)
//
// 	checkInvoke(t, stub, [][]byte{[]byte("initVote"), []byte("red"), []byte("green")})
// 	checkState(t, stub, "red", "{0}")
// 	checkState(t, stub, "green", "{0}")
// 	checkQuery(t, stub, "queryOptions", "", "[\"green\",\"red\"]")
// }

// func Test_Vote(t *testing.T) {
// 	scc := new(SmartContract)
// 	stub := shim.NewMockStub("test_vote", scc)
//
// 	checkInvoke(t, stub, [][]byte{[]byte("initVote"), []byte("red"), []byte("green")})
// 	checkInvoke(t, stub, [][]byte{[]byte("vote"), []byte("red")})
// 	checkState(t, stub, "red", "{1}")
// 	checkState(t, stub, "green", "{}")
// 	checkQuery(t, stub, "queryVotes", "", "[{key:\"green\",value:{0},{key:\"red\",value:{1}]")
// }

func Test_SubmitVote(t *testing.T) {
	scc := new(SmartContract)
	stub := shim.NewMockStub("test_submitVote", scc)

	stub.MockTransactionStart("t123")
	stateAsBytes, _ := json.Marshal(VOTE)
	stub.PutState("state", stateAsBytes)
	stub.MockTransactionEnd("t123")

	checkInvoke(t, stub, [][]byte{[]byte("submitVote"), []byte("")})
}
