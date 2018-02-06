package main

import (
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/stretchr/testify/assert"
)

func Test_End2End(t *testing.T) {

//	users := []int{10, 50, 100, 200, 300, 400, 500, 600, 700, 800, 900, 1000, 2000, 3000, 4000, 5000, 6000, 7000, 8000, 9000, 10000}
	users := []int{7000, 8000, 9000, 10000}

	for _, numUser := range users {
		runScenario(t, numUser)
	}
}

func runScenario(t *testing.T, numberOfUsers int) {
	scc := new(SmartContract)
	stub := shim.NewMockStub("test_invalid_question", scc)

	var voters []string

	stub.MockTransactionStart("t123")
	response := scc.initVote(stub)
	assert.True(t, shim.OK == response.Status)
	stub.MockTransactionEnd("t123")

	for i := 0; i < numberOfUsers; i++ {
		userID := "userID" + strconv.Itoa(i)
		voters = append(voters, userID)
	}

	stub.MockTransactionStart("t124")
	response = scc.setEligible(stub, voters)
	assert.True(t, shim.OK == response.Status)
	stub.MockTransactionEnd("t124")

	stub.MockTransactionStart("t127")
	response = scc.beginSignUp(stub, []string{""})
	assert.True(t, shim.OK == response.Status)
	stub.MockTransactionEnd("t128")

	for _, voter := range voters {
		stub.MockTransactionStart("t126")
		response = scc.register(stub, []string{voter})
		assert.True(t, shim.OK == response.Status)
		stub.MockTransactionEnd("t126")
	}

	logger.Info("Registered voters: " + strconv.Itoa(numberOfUsers))

	stub.MockTransactionStart("t125")
	response = scc.finishRegistrationPhase(stub)
	assert.True(t, shim.OK == response.Status)
	stub.MockTransactionEnd("t125")

	stub.MockTransactionStart("t124")
	start := time.Now()

	for _, voter := range voters {

		response := scc.submitVote(stub, []string{voter, "1"})
		assert.True(t, shim.OK == response.Status)

	}

	elapsed := time.Since(start)
	stub.MockTransactionEnd("t124")
	log.Printf("Voting took %s", elapsed)
}
