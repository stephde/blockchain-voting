package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
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

// from secp256k1.
func generateKeyPair() (pubkey, privkey []byte) {
	key, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	pubkey = elliptic.Marshal(crypto.S256(), key.X, key.Y)
	return pubkey, math.PaddedBigBytes(key.D, 32)
}

func Test_verifyZKP(t *testing.T) {
	publicKey, privateKey := generateKeyPair()
	userID := ""

	xGx := new(big.Int).SetBytes(publicKey[0:31])
	xGy := new(big.Int).SetBytes(publicKey[32:63])
	xG := []big.Int{*xGx, *xGy}
	r := big.NewInt(3)
	vG := []big.Int{3, 4, 5}

	scc := new(SmartContract)
	assert.True(t, scc.verifyZKP(userID, xG, *r, vG))
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
		[]byte("userId"),
		[]byte("[1, 2]"),
		[]byte("[1, 2, 3]"),
		[]byte("1"),
	})

	var totalRegistered int
	GetState(stub, "totalRegistered", &totalRegistered)
	assert.Equal(t, 1, totalRegistered)

	GetState(stub, "registered", &registered)
	assert.True(t, registered["userId"])

	// TODO: verify that voter was storede
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
