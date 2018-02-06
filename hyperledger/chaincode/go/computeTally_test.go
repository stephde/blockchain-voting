package main

/*
func Test_ComputeTally(t *testing.T) {
	stub := shim.NewMockStub("test_computeTally", new(SmartContract))

	stub.MockTransactionStart("t123")
	PutState(stub, "state", VOTE)
	PutState(stub, "totalRegistered", 2)
	compositeIndexName := "varName~userID~vote~txID"

	var i int
	for i = 0; i < 2; i++ {
		compositeKey, _ := stub.CreateCompositeKey(compositeIndexName, []string{"vote", strconv.Itoa(i), strconv.Itoa(i), strconv.Itoa(1)})
		stub.PutState(compositeKey, []byte{0x00})
	}

	stub.MockTransactionEnd("t123")

	result := Result{2, map[int]int{0: 1, 1: 1}}
	bytes, _ := json.Marshal(result)

	checkQuery(t, stub, "computeTally", string(bytes))
}
*/
