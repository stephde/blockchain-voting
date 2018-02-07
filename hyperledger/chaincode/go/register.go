package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

func (s *SmartContract) register(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if !s.inState(stub, SIGNUP) {
		return shim.Error("Wrong state, expected SIGNUP")
	}

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments, expecting userID")
	}

	userID := args[0]

	var eligible map[string]bool
	GetState(stub, "eligible", &eligible)
	logger.Info("eligible {}", eligible)

	var registered map[string]bool
	GetState(stub, "registered", &registered)
	logger.Info("registered {}", registered)

	isEligible := eligible[userID]
	isRegistered := registered[userID]

	logger.Info("iseligible {}", isEligible)
	logger.Info("isregistered {}", isRegistered)

	voter := Voter{userID, -1}

	if isEligible && !isRegistered {
		registered[userID] = true
		logger.Info("registered {}", registered)
		PutState(stub, "registered", registered)

		var totalRegistered int
		GetState(stub, "totalRegistered", &totalRegistered)
		totalRegistered++
		logger.Info("totalRegistered {}", totalRegistered)
		PutState(stub, "totalRegistered", totalRegistered)

		var voters map[string]Voter
		GetState(stub, "voters", &voters)
		logger.Info("voters {}", voters)
		voters[userID] = voter
		logger.Info("voters {}", voters)
		PutState(stub, "voters", voters)
		return shim.Success(nil)
	}

	return shim.Error("not eligible")

}
