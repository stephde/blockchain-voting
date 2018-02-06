package main

/*
func Test_SubmitVote(t *testing.T) {
	stub := shim.NewMockStub("test_submitVote", new(SmartContract))

	userID := "1"
	stub.MockTransactionStart("t123")
	PutState(stub, "state", VOTE)

	registered := map[string]bool{userID: true}
	PutState(stub, "registered", registered)

	votecast := map[string]bool{userID: false}
	PutState(stub, "votecast", votecast)

	PutState(stub, "voters", map[string]Voter{})
	stub.MockTransactionEnd("t123")

	checkInvoke(t, stub, [][]byte{[]byte("submitVote"), []byte(userID), []byte(strconv.Itoa(1))})

	GetState(stub, "votecast", &votecast)
	assert.True(t, votecast[userID])

	var voters map[string]Voter
	GetState(stub, "voters", &voters)
	assert.Equal(t, 1, voters[userID].Vote)
}

func Test_InvalidUserID(t *testing.T) {
	stub := shim.NewMockStub("test_submitVote", new(SmartContract))

	stub.MockTransactionStart("t123")
	PutState(stub, "state", VOTE)
	stub.MockTransactionEnd("t123")

	checkFailingInvoke(t, stub, [][]byte{[]byte("submitVote"), []byte("invalidUser"), []byte(strconv.Itoa(1))})
}

func Test_DuplicateVote(t *testing.T) {
	stub := shim.NewMockStub("test_submitVote", new(SmartContract))

	userID := "1"
	stub.MockTransactionStart("t123")
	PutState(stub, "state", VOTE)

	registered := map[string]bool{userID: true}
	PutState(stub, "registered", registered)

	votecast := map[string]bool{userID: true}
	PutState(stub, "votecast", votecast)
	stub.MockTransactionEnd("t123")

	checkFailingInvoke(t, stub, [][]byte{[]byte("submitVote"), []byte(userID), []byte(strconv.Itoa(1))})
}
*/
