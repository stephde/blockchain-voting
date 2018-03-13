package main

import (
	"encoding/json"
	"fmt"
	"strconv"
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

func checkQuery(t *testing.T, stub *shim.MockStub, function string, value string) {
	res := stub.MockInvoke("1", [][]byte{[]byte(function)})
	if res.Status != shim.OK {
		fmt.Println("Query", function, "failed", string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		fmt.Println("Query", function, "failed to get value")
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

func checkFailingInvoke(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInvoke("1", args)
	if res.Status != shim.ERROR {
		fmt.Println("Invoke", args, "succeeded", string(res.Message))
		t.FailNow()
	}
}

func Test_InvalidFunctionName(t *testing.T) {
	stub := shim.NewMockStub("test_invalid_question", new(SmartContract))
	checkFailingInvoke(t, stub, [][]byte{[]byte("someInvalidFunction")})
}

func Test_End2End(t *testing.T) {
	userID1 := "a"
	userID2 := "b"
	userID3 := "c"

	stub := shim.NewMockStub("test_invalid_question", new(SmartContract))
	// init vote
	checkInvoke(t, stub, [][]byte{[]byte("initVote")})

	// Set eligible users
	checkInvoke(t, stub, [][]byte{[]byte("setEligible"),
		[]byte(userID1),
		[]byte(userID2),
		[]byte(userID3)})

	// Begin SignUp
	checkInvoke(t, stub, [][]byte{[]byte("beginSignUp"), []byte("Do you like Blockchain?")})

	// Register users
	checkInvoke(t, stub, [][]byte{[]byte("register"), []byte(userID1)})
	checkInvoke(t, stub, [][]byte{[]byte("register"), []byte(userID2)})
	checkInvoke(t, stub, [][]byte{[]byte("register"), []byte(userID3)})

	// Begin Vote
	checkInvoke(t, stub, [][]byte{[]byte("finishRegistrationPhase")})

	// Vote
	checkInvoke(t, stub, [][]byte{[]byte("submitVote"), []byte(userID1), []byte(strconv.Itoa(0))})
	checkInvoke(t, stub, [][]byte{[]byte("submitVote"), []byte(userID2), []byte(strconv.Itoa(1))})
	checkInvoke(t, stub, [][]byte{[]byte("submitVote"), []byte(userID3), []byte(strconv.Itoa(1))})

	// Compute computeTally
	expectedResult, _ := json.Marshal(Result{3, map[int]int{0: 1, 1: 2}})
	checkQuery(t, stub, "computeTally", string(expectedResult))
}
