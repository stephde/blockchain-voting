package main

import (
	"fmt"

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
	// Retrieve info needed for the update procedure
	txid := stub.GetTxID()
	compositeIndexName := "varName~userID~txID"
	name := "register"

	// Create the composite key that will allow us to query for all deltas on a particular variable
	compositeKey, compositeErr := stub.CreateCompositeKey(compositeIndexName, []string{name, userID, txid})
	if compositeErr != nil {
		return shim.Error(fmt.Sprintf("Could not create a registration composite key for %s: %s", userID, compositeErr.Error()))
	}

	var eligible map[string]bool
	GetState(stub, "eligible", &eligible)

	if eligible[userID] {
		// Save the composite key index
		compositePutErr := stub.PutState(compositeKey, []byte{0x00})
		if compositePutErr != nil {
			return shim.Error(fmt.Sprintf("Could not put registration for %s in the ledger: %s", userID, compositePutErr.Error()))
		}
	}

	return shim.Success(nil)
}
