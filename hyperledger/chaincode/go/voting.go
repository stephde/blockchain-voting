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

/*
 * The sample smart contract for documentation topic:
 * Writing Your First Blockchain Application
 */

package main

/* Imports
 * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

type Vote struct {
	count int    `json:"count"`
}

/*
 * The Init method is called when the Smart Contract "fabcar" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "fabcar"
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "vote" {
		return s.vote(APIstub, args)
	} else if function == "queryVotes" {
		return s.queryVotes(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) vote(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

		if len(args) != 1 {
			return shim.Error("Incorrect number of arguments. Expecting 1")
		}

		voteAsBytes, _ := APIstub.GetState(args[0])
		vote := Vote{}

		json.Unmarshal(voteAsBytes, &vote)
		vote.count = vote.count + 1

		voteAsBytes, _ = json.Marshal(vote)
		APIstub.PutState(args[0], voteAsBytes)

		return shim.Success(nil)
}

func (s *SmartContract) registerParty(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

		if len(args) != 1 {
			return shim.Error("Incorrect number of arguments. Expecting 1")
		}

		vote := Vote{}
		voteAsBytes, _ = json.Marshal(vote)
		APIstub.PutState(args[0], voteAsBytes)

		return shim.Success(nil)
}


func (s *SmartContract) queryVotes(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

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
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllVotes:\n%s\n", buffer.String())

	// TODO: switch to buffer.String()
	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	vote := Vote{}

	for index, party := range args {
		fmt.Println(index, " is ", party)
		APIstub.PutState(party, vote)
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
