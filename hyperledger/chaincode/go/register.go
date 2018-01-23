package main

import (
	"crypto/ecdsa"
	"encoding/json"
	"math/big"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

func (s *SmartContract) register(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if !s.inState(stub, SIGNUP) {
		return shim.Error("Wrong state")
	}

	if len(args) != 3 {
		return shim.Error("Wrong number of arguments, expected 3")
	}

	// TODO: what do these names mean?
	// TODO: error handling

	userIDBytes, _ := stub.GetCreator()
	userID := string(userIDBytes)

	// Public key of voter: xG
	var xG []*big.Int
	json.Unmarshal([]byte(args[0]), &xG)

	publicKey := new(ecdsa.PublicKey)
	publicKey.X = xG[0]
	publicKey.Y = xG[1]

	var vG []*big.Int
	json.Unmarshal([]byte(args[1]), &vG)

	var r big.Int
	json.Unmarshal([]byte(args[2]), &r)

	var eligible map[string]bool
	GetState(stub, "eligible", &eligible)

	var registered map[string]bool
	GetState(stub, "registered", &registered)

	isEligible := eligible[userID]
	isRegistered := registered[userID]

	var emptyVote []*big.Int
	emptyReconstructedKey := new(ecdsa.PublicKey)
	voter := Voter{userID, publicKey, emptyReconstructedKey, emptyVote}

	if isEligible && !isRegistered && s.verifyZKP(voter, &r, vG) {
		registered[userID] = true
		PutState(stub, "registered", registered)

		// voter := Voter{userID, xG, {}, {}}

		var totalRegistered int
		GetState(stub, "totalRegistered", &totalRegistered)
		totalRegistered = totalRegistered + 1
		PutState(stub, "totalRegistered", totalRegistered)
	}

	return shim.Success(nil)
}
