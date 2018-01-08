/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package main

/* Imports
 * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

var logger = shim.NewLogger("myChaincode")

// Define the Smart Contract structure
type SmartContract struct {
}

type Voter struct {
	address          string
	registeredKey    [2]int
	reconstructedKey [2]int
	vote             [2]int
}

// TODO: get from store
// func getVoter(address string) ([2]int, [2]int) {
// 	index := addressid[address]
// 	return voters[index].registeredKey, voters[index].reconstructedKey
// }

/*
 * Hyperledger Chaincode Interface
 */

/*
 * The Init method is called when the Smart Contract is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initVote()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract
 * The calling application program has also specified the particular smart contract function to be called,
 * with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	switch function {
	case "initVote":
		return s.initVote(APIstub, args)
	case "beginSignUp":
		return s.beginSignUp(APIstub, args)
	case "submitVote":
		return s.submitVote(APIstub, args)
	case "setEligible":
		return s.setEligible(APIstub, args)
	case "register":
		return s.register(APIstub, args)
	default:
		return shim.Error("Invalid Smart Contract function name.")
	}
}

/*
 * Custom functions
 */

func (s *SmartContract) register(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if !s.inState(stub, SIGNUP) {
		return shim.Error("Wrong state")
	}

	if len(args) != 4 {
		return shim.Error("Wrong number of arguments, expected 4")
	}

	// TODO: what do these names mean?
	userId := args[0]
	xG := args[1]
	vG := args[2]
	r := args[3]

	var eligible map[string]bool
	GetState(stub, "eligible", &eligible)

	var registered map[string]bool
	GetState(stub, "registered", &registered)

	isEligible := eligible[userId]
	isRegistered := registered[userId]

	if isEligible && !isRegistered && s.verifyZKPString(xG, r, vG) {

	}

	return shim.Error("not implemented yet")
}

func (s *SmartContract) beginSignUp(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if !s.inState(stub, SETUP) {
		return shim.Error("Wrong state")
	}

	question := args[0]
	PutState(stub, "question", question)
	s.transitionToState(stub, SIGNUP)

	return shim.Success(nil)
}

func (s *SmartContract) computeTally(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Error("Not implemented yet")
}

// What do these parameters mean???
func (s *SmartContract) verifyZKPString(xG string, r string, vG string) bool {
	// xG [2]int, r int, vG [3]int
	return true
}

func (s *SmartContract) verifyZKP(xG [2]int, r int, vG [3]int) bool {
	return true
}

// func (s *SmartContract) vote(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
//
// 	if len(args) != 1 {
// 		return shim.Error("Incorrect number of arguments. Expecting 1")
// 	}
//
// 	voteAsBytes, _ := APIstub.GetState(args[0])
// 	vote := Vote{}
//
// 	json.Unmarshal(voteAsBytes, &vote)
// 	vote.count = vote.count + 1
//
// 	logger.Info("Voted for", args[0], "-", vote)
//
// 	voteAsBytes, _ = json.Marshal(vote)
// 	APIstub.PutState(args[0], voteAsBytes)
//
// 	return shim.Success(nil)
// }

// func (s *SmartContract) queryVotes(APIstub shim.ChaincodeStubInterface) sc.Response {
// 	// buffer is a JSON array containing QueryResults
// 	buffer, err := s.stateToJSON(APIstub)
// 	if err != nil {
// 		return shim.Error(err.Error())
// 	}
// 	logger.Info("QueryVotes", buffer.String())
// 	return shim.Success(buffer.Bytes())
// }
//
// func (s *SmartContract) queryOptions(APIstub shim.ChaincodeStubInterface) sc.Response {
// 	resultsIterator, err := APIstub.GetStateByRange("", "")
// 	if err != nil {
// 		return shim.Error(err.Error())
// 	}
// 	defer resultsIterator.Close()
//
// 	// buffer is a JSON array containing QueryResults
// 	var buffer bytes.Buffer
// 	buffer.WriteString("[")
//
// 	bArrayMemberAlreadyWritten := false
// 	for resultsIterator.HasNext() {
// 		queryResponse, _ := resultsIterator.Next()
//
// 		// Add a comma before array members, suppress it for the first array member
// 		if bArrayMemberAlreadyWritten == true {
// 			buffer.WriteString(",")
// 		}
// 		buffer.WriteString("\"")
// 		buffer.WriteString(queryResponse.Key)
// 		buffer.WriteString("\"")
// 		bArrayMemberAlreadyWritten = true
// 	}
// 	buffer.WriteString("]")
//
// 	logger.Info(buffer.String())
//
// 	return shim.Success(buffer.Bytes())
// }

// func (s *SmartContract) stateToJSON(APIstub shim.ChaincodeStubInterface) (bytes.Buffer, error) {
// 	resultsIterator, err := APIstub.GetStateByRange("", "")
// 	var buffer bytes.Buffer
// 	if err != nil {
// 		return buffer, err
// 	}
// 	defer resultsIterator.Close()
//
// 	// buffer is a JSON array containing QueryResults
// 	buffer.WriteString("[")
//
// 	bArrayMemberAlreadyWritten := false
// 	for resultsIterator.HasNext() {
// 		queryResponse, _ := resultsIterator.Next()
//
// 		// Add a comma before array members, suppress it for the first array member
// 		if bArrayMemberAlreadyWritten == true {
// 			buffer.WriteString(",")
// 		}
// 		buffer.WriteString("{\"key\":")
// 		buffer.WriteString("\"")
// 		buffer.WriteString(queryResponse.Key)
// 		buffer.WriteString("\"")
//
// 		buffer.WriteString(", \"value\":")
// 		buffer.WriteString(string(queryResponse.Value))
// 		buffer.WriteString("}")
// 		bArrayMemberAlreadyWritten = true
// 	}
// 	buffer.WriteString("]")
//
// 	return buffer, nil
// }

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
