package main

import (
	"crypto/ecdsa"
	"log"
	"math/big"
	"strconv"
	"testing"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/stretchr/testify/assert"
)

/*
 * Ignored per default to reduce execution time of unit tests
 */
// func Test_End2End(t *testing.T) {
func Test_End2End(t *testing.T) {

	users := []int{10, 50, 100, 200, 300, 400, 500, 600, 700, 800, 900, 1000, 2000, 3000, 4000, 5000, 6000, 7000, 8000, 9000, 10000}
	//users := []int{10}

	for _, numUser := range users {
		runScenario(t, numUser)
	}
}

func runScenario(t *testing.T, numberOfUsers int) {
	scc := new(SmartContract)
	stub := shim.NewMockStub("test_invalid_question", scc)

	voters := []Voter{}
	registered := map[string]bool{}
	votecast := map[string]bool{}

	for i := 0; i < numberOfUsers; i++ {
		publicKeyECDSA, _ := generateKeyPair()
		var emptyVote []*big.Int
		emptyReconstructedKey := new(ecdsa.PublicKey)
		userID := "userID" + strconv.Itoa(i)
		voter := Voter{userID, publicKeyECDSA, emptyReconstructedKey, emptyVote}
		voters = append(voters, voter)
		registered[userID] = true
		votecast[userID] = false
	}

	logger.Info("Registered all voters in e2e")

	stub.MockTransactionStart("t123")
	PutState(stub, "registered", registered)
	PutState(stub, "votecast", votecast)

	voters = scc.reconstructKeys(numberOfUsers, voters)
	PutState(stub, "voters", voters)
	stub.MockTransactionEnd("t123")

	stub.MockTransactionStart("t1235535")
	var foo []Voter
	GetState(stub, "voters", &foo)
	stub.MockTransactionEnd("t1235535")

	logger.Info("Reconstructed all voting keys in e2e")

	start := time.Now()

	for _, voter := range voters {
		w := generateRandomSeed()
		r := generateRandomSeed()
		d := generateRandomSeed()
		x := generateRandomSeed()
		result1, result2 := create1outof2ZKPYesVote(voter, w, r, d, x)

		y := new(ecdsa.PublicKey)
		a1 := new(ecdsa.PublicKey)
		b1 := new(ecdsa.PublicKey)
		a2 := new(ecdsa.PublicKey)
		b2 := new(ecdsa.PublicKey)

		y.X, y.Y = result1[0], result1[1]
		a1.X, a1.Y = result1[2], result1[3]
		b1.X, b1.Y = result1[4], result1[5]
		a2.X, a2.Y = result1[6], result1[7]
		b2.X, b2.Y = result1[8], result1[9]

		stub.MockTransactionStart("t124")
		response := scc.submitVoteInternal(stub, voter.address, result2, y, a1, b1, a2, b2)
		assert.True(t, shim.OK == response.Status)
		stub.MockTransactionEnd("t124")
	}

	elapsed := time.Since(start)
	log.Printf("Voting took %s", elapsed)
}
