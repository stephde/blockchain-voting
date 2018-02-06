package main

/*
func Test_BeginSignup(t *testing.T) {
	stub := shim.NewMockStub("test_beginSignup", new(SmartContract))

	question := "What is the question?"

	stub.MockTransactionStart("t123")
	PutState(stub, "state", SETUP)
	stub.MockTransactionEnd("t123")

	checkInvoke(t, stub, [][]byte{[]byte("beginSignUp"), []byte(question)})

	var state StateEnum
	GetState(stub, "state", &state)
	assert.Equal(t, SIGNUP, state)

	var questionFromState string
	GetState(stub, "question", &questionFromState)
	assert.Equal(t, question, questionFromState)
}

func Test_MissingQuestion(t *testing.T) {
	stub := shim.NewMockStub("test_missing_question", new(SmartContract))
	stub.MockTransactionStart("t123")
	PutState(stub, "state", SETUP)
	stub.MockTransactionEnd("t123")

	checkFailingInvoke(t, stub, [][]byte{[]byte("beginSignUp")})
}

func Test_WrongState(t *testing.T) {
	stub := shim.NewMockStub("Test_WrongState", new(SmartContract))

	question := "What is the question?"
	checkFailingInvoke(t, stub, [][]byte{[]byte("beginSignUp"), []byte(question)})
}
*/
