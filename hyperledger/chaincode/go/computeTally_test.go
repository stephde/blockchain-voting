package main

/*

func Test_ComputeTally(t *testing.T) {
	stub := shim.NewMockStub("test_computeTally", new(SmartContract))

	stub.MockTransactionStart("t123")
	PutState(stub, "state", VOTE)
	PutState(stub, "totalRegistered", 2)
	voters := map[string]Voter{"1": Voter{"1", 0}, "2": Voter{"2", 1}}
	PutState(stub, "voters", voters)
	votecast := map[string]bool{"1": true, "2": true}
	PutState(stub, "votecast", votecast)
	stub.MockTransactionEnd("t123")

	result := Result{2, map[int]int{0: 1, 1: 1}}
	bytes, _ := json.Marshal(result)

	checkQuery(t, stub, "computeTally", string(bytes))
}

*/
