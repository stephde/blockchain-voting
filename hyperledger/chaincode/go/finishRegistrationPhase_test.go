package main

/*
func Test_WrongStateFinishRegistrationPhase(t *testing.T) {
	stub := shim.NewMockStub("Test_WrongState", new(SmartContract))

	checkFailingInvoke(t, stub, [][]byte{[]byte("finishRegistrationPhase")})
}

func Test_FinishRegistrationPhase(t *testing.T) {
	stub := shim.NewMockStub("test_beginVote", new(SmartContract))
	compositeIndexName := "varName~userID~txID"

	stub.MockTransactionStart("t123")
	PutState(stub, "state", SIGNUP)
	var i int
	for i = 0; i < 3; i++ {
		compositeKey, _ := stub.CreateCompositeKey(compositeIndexName, []string{"register", strconv.Itoa(i), strconv.Itoa(1)})
		stub.PutState(compositeKey, []byte{0x00})
	}

	stub.MockTransactionEnd("t123")

	checkInvoke(t, stub, [][]byte{[]byte("finishRegistrationPhase")})

	var state StateEnum
	GetState(stub, "state", &state)
	assert.Equal(t, VOTE, state)
}
*/
