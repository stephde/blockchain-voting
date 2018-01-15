package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

func (s *SmartContract) computeTally(stub shim.ChaincodeStubInterface) sc.Response {

	var totalRegistered int
	GetState(stub, "totalregistered", &totalRegistered)

	// for (i := 0; i < totalRegistered; i++) {
	// TODO: confirm that all votes have been cast...
	// }

	return shim.Error("Not implemented yet")

}
