package main

import (
	"fmt"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/stretchr/testify/assert"
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

func Test_ComputeTally(t *testing.T) {

	scc := new(SmartContract)
	stub := shim.NewMockStub("test_computeTally", scc)

	stub.MockTransactionStart("t123")
	PutState(stub, "state", VOTE)
	stub.MockTransactionEnd("t123")

	checkInvoke(t, stub, [][]byte{[]byte("computeTally")})
}

func Test_Register(t *testing.T) {

	scc := new(SmartContract)
	stub := shim.NewMockStub("test_register", scc)

	stub.MockTransactionStart("t123")
	PutState(stub, "state", SIGNUP)
	eligible := map[string]bool{"userId": true}
	PutState(stub, "eligible", eligible)
	registered := map[string]bool{"userId": false}
	PutState(stub, "registered", registered)
	stub.MockTransactionEnd("t124")

	checkInvoke(t, stub, [][]byte{
		[]byte("register"),
		[]byte("[1, 2]"),
		[]byte("[1, 2, 3]"),
		[]byte("1"),
	})

	var totalRegistered int
	GetState(stub, "totalRegistered", &totalRegistered)
	assert.Equal(t, 1, totalRegistered)

	GetState(stub, "registered", &registered)
	assert.True(t, registered["userId"])

	// TODO: verify that voter was stored
}

func Test_SetEligible(t *testing.T) {

	scc := new(SmartContract)
	stub := shim.NewMockStub("test_setEligible", scc)

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

func Test_InitVote(t *testing.T) {

	scc := new(SmartContract)
	stub := shim.NewMockStub("test_beginSignup", scc)

	checkInvoke(t, stub, [][]byte{[]byte("initVote")})

	var state StateEnum
	GetState(stub, "state", &state)
	assert.Equal(t, SETUP, state)
}

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

func Test_SubmitVote(t *testing.T) {

	scc := new(SmartContract)
	stub := shim.NewMockStub("test_submitVote", scc)

	stub.MockTransactionStart("t123")
	PutState(stub, "state", VOTE)
	stub.MockTransactionEnd("t123")

	userId := "1"
	stub.MockTransactionStart("t124")
	registered := map[string]bool{userId: true}
	PutState(stub, "registered", registered)
	stub.MockTransactionEnd("t124")

	stub.MockTransactionStart("t125")
	votecast := map[string]bool{userId: false}
	PutState(stub, "votecast", votecast)
	stub.MockTransactionEnd("t125")

	checkInvoke(t, stub, [][]byte{[]byte("submitVote"), []byte(userId)})
}
