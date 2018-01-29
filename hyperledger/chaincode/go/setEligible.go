package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

func (s *SmartContract) setEligible(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if !s.inState(stub, SETUP) {
		return shim.Error("Wrong state, expected SETUP")
	}

	eligible := make(map[string]bool)

	for _, a := range args {
		eligible[a] = true
	}

	PutState(stub, "eligible", eligible)
	PutState(stub, "totalEligible", len(args))

	return shim.Success(nil)
}
