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
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

var logger = shim.NewLogger("myChaincode")

// Define the Smart Contract structure
type SmartContract struct {
}

type Vote struct {
	count int
}

type State struct {
	state StateEnum
}

type StateEnum int

const (
	SETUP StateEnum = iota
	SIGNUP
	VOTE
	FINISHED
)

var states = [...]string{
	"SETUP",
	"SIGNUP",
	"VOTE",
	"FINISHED",
}

func (state StateEnum) String() string {
	return states[state]
}

/*
 * Hyperledger Chaincode Interface
 */

/*
 * The Init method is called when the Smart Contract is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
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
	case "submitVote":
		return s.submitVote(APIstub, args)
	case "queryVotes":
		return s.queryVotes(APIstub)
	case "queryOptions":
		return s.queryOptions(APIstub)
	default:
		return shim.Error("Invalid Smart Contract function name.")
	}
}

/*
 * Custom functions
 */

func (s *SmartContract) computeTally(APIstub shim.ChaincodeStubInterface) sc.Response {

	return shim.Error("Foooo")
}

func (s *SmartContract) submitVote(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if !s.inState(stub, VOTE) {
		return shim.Error("Wrong state")
	}

	// get sender address
	creator, err := stub.GetCreator()
	publicKey := "SomePublicKey"

	logger.Info("Creator is ", creator)
	logger.Info("Err is ", err)

	// Make sure the sender can vote and hasn't already voted
	registered := make(map[string]bool)
	votecast := make(map[string]bool)
	GetState(stub, "registered", registered)
	GetState(stub, "votecast", votecast)

	_, ok1 := registered[publicKey]
	_, ok2 := votecast[publicKey]

	if ok1 && !ok2 {
		// User is registered and did not cast vote yet
		logger.Info("User is allowed to vote")
	} else {
		logger.Info("User is not allowed to vote")
	}

	return shim.Error("Not implemented yet")
}

func GetState(stub shim.ChaincodeStubInterface, name string, v interface{}) error {
	bytes, err := stub.GetState(name)
	if err != nil {
		return errors.New("Failed to get state")
	}
	if bytes == nil {
		return errors.New("Entity not found")
	}

	return json.Unmarshal(bytes, &v)
}

func (s *SmartContract) inState(APIstub shim.ChaincodeStubInterface, expectedState StateEnum) bool {
	stateAsBytes, _ := APIstub.GetState("state")
	state := State{}

	json.Unmarshal(stateAsBytes, &state)

	logger.Info("State is " + state.state.String())
	return expectedState == state.state
}

// What do these parameters mean???
func (s *SmartContract) verifyZKP(xG [2]int, r int, vG [3]int) bool {
	return true
}

func (s *SmartContract) vote(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	voteAsBytes, _ := APIstub.GetState(args[0])
	vote := Vote{}

	json.Unmarshal(voteAsBytes, &vote)
	vote.count = vote.count + 1

	logger.Info("Voted for", args[0], "-", vote)

	voteAsBytes, _ = json.Marshal(vote)
	APIstub.PutState(args[0], voteAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) queryVotes(APIstub shim.ChaincodeStubInterface) sc.Response {
	// buffer is a JSON array containing QueryResults
	buffer, err := s.stateToJSON(APIstub)
	if err != nil {
		return shim.Error(err.Error())
	}
	logger.Info("QueryVotes", buffer.String())
	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) queryOptions(APIstub shim.ChaincodeStubInterface) sc.Response {
	resultsIterator, err := APIstub.GetStateByRange("", "")
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, _ := resultsIterator.Next()

		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	logger.Info(buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) stateToJSON(APIstub shim.ChaincodeStubInterface) (bytes.Buffer, error) {
	resultsIterator, err := APIstub.GetStateByRange("", "")
	var buffer bytes.Buffer
	if err != nil {
		return buffer, err
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, _ := resultsIterator.Next()

		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"value\":")
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return buffer, nil
}

func (s *SmartContract) initVote(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	vote := Vote{0}

	for _, party := range args {
		voteAsBytes, _ := json.Marshal(vote)
		APIstub.PutState(party, voteAsBytes)
		fmt.Println("Added", party)
	}

	return shim.Success(nil)
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
