package main

/*
func Test_Register(t *testing.T) {
	stub := shim.NewMockStub("test_register", new(SmartContract))

	userID := "userId"

	stub.MockTransactionStart("t123")
	PutState(stub, "state", SIGNUP)
	eligible := map[string]bool{userID: true}
	PutState(stub, "eligible", eligible)
	registered := map[string]bool{userID: false}
	PutState(stub, "registered", registered)
	voters := map[string]Voter{}
	PutState(stub, "voters", voters)
	stub.MockTransactionEnd("t123")

	checkInvoke(t, stub, [][]byte{
		[]byte("register"),
		[]byte(userID),
	})

	var totalRegistered int
	GetState(stub, "totalRegistered", &totalRegistered)
	assert.Equal(t, 1, totalRegistered)

	GetState(stub, "registered", &registered)
	assert.True(t, registered[userID])

	GetState(stub, "voters", &voters)
	assert.Equal(t, 1, len(voters))
	_, found := voters[userID]
	assert.True(t, found)
}
*/
