package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

func (s *SmartContract) register(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if !s.inState(stub, SIGNUP) {
		return shim.Error("Wrong state")
	}

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments, expecting userId")
	}

	userID := args[0]

	var eligible map[string]bool
	GetState(stub, "eligible", &eligible)

	var registered map[string]bool
	GetState(stub, "registered", &registered)

	isEligible := eligible[userID]
	isRegistered := registered[userID]

	voter := Voter{userID, -1}

	if isEligible && !isRegistered {
		registered[userID] = true
		PutState(stub, "registered", registered)

		var totalRegistered int
		GetState(stub, "totalRegistered", &totalRegistered)
		totalRegistered = totalRegistered + 1
		PutState(stub, "totalRegistered", totalRegistered)

		var voters map[string]Voter
		GetState(stub, "voters", &voters)
		voters[userID] = voter
		PutState(stub, "voters", voters)
	}

	return shim.Success(nil)
}
