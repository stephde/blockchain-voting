package main

import (
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/stretchr/testify/assert"
)

// Disabled per default to reduce execution time of unit tests
//func Test_End2End(t *testing.T) {
func Ignore_End2End(t *testing.T) {
	users := []int{10, 50, 100, 200, 300, 400, 500, 600, 700, 800, 900, 1000, 2000, 3000, 4000, 5000, 6000}

	for _, numUser := range users {
		runScenario(t, numUser)
	}
}

func runScenario(t *testing.T, numberOfUsers int) {
	scc := new(SmartContract)
	stub := shim.NewMockStub("test_invalid_question", scc)

	var voters []string
	start := time.Now()

	stub.MockTransactionStart("t123")
	response := scc.initVote(stub)
	assert.True(t, shim.OK == response.Status)
	stub.MockTransactionEnd("t123")

	log.Printf("Init Vote took %s", time.Since(start))
	start = time.Now()

	for i := 0; i < numberOfUsers; i++ {
		userID := "userID" + strconv.Itoa(i)
		voters = append(voters, userID)
	}

	log.Printf("Voter list took %s", time.Since(start))
	start = time.Now()

	stub.MockTransactionStart("t124")
	response = scc.setEligible(stub, voters)
	assert.True(t, shim.OK == response.Status)
	stub.MockTransactionEnd("t124")

	log.Printf("Set eligible took %s", time.Since(start))
	start = time.Now()

	stub.MockTransactionStart("t127")
	response = scc.beginSignUp(stub, []string{""})
	assert.True(t, shim.OK == response.Status)
	stub.MockTransactionEnd("t128")

	log.Printf("Begin Signup took %s", time.Since(start))
	start = time.Now()

	stub.MockTransactionStart("t126")
	for _, voter := range voters {
		response = scc.register(stub, []string{voter})
	}
	stub.MockTransactionEnd("t126")

	log.Printf("Register took %s", time.Since(start))
	start = time.Now()

	logger.Info("Registered voters: " + strconv.Itoa(numberOfUsers))

	stub.MockTransactionStart("t125")
	response = scc.finishRegistrationPhase(stub)
	assert.True(t, shim.OK == response.Status)
	stub.MockTransactionEnd("t125")

	log.Printf("FinishRegistration took %s", time.Since(start))
	start = time.Now()

	stub.MockTransactionStart("t124")

	for _, voter := range voters {

		response := scc.submitVote(stub, []string{voter, "1"})
		assert.True(t, shim.OK == response.Status)

	}

	stub.MockTransactionEnd("t124")
	log.Printf("Voting took %s", time.Since(start))
}
