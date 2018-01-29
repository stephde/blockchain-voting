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

func Test_InvalidFunctionName(t *testing.T) {
	scc := new(SmartContract)
	stub := shim.NewMockStub("test_invalid_question", scc)

	res := stub.MockInvoke("1", [][]byte{[]byte("someInvalidFunction")})
	if res.Status != shim.ERROR {
		fmt.Println("Query test_question succeeded", string(res.Message))
		t.FailNow()
	}
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
