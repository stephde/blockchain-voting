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
	"crypto/sha256"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

var logger = shim.NewLogger("myChaincode")

// Define the Smart Contract structure
type SmartContract struct {
}

type Voter struct {
	address          string
	registeredKey    XG
	reconstructedKey [2]int
	vote             [2]int
}

type XG struct {
	xG1 int
	xG2 int
}

type VG struct {
	vG1 int
	vG2 int
	vG3 int
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
	case "computeTally":
		return s.computeTally(APIstub)
	default:
		return shim.Error("Invalid Smart Contract function name.")
	}
}

/*
 * Custom functions
 */

func (s *SmartContract) computeTally(stub shim.ChaincodeStubInterface) sc.Response {

	var totalRegistered int
	GetState(stub, "totalregistered", &totalRegistered)

	// for (i := 0; i < totalRegistered; i++) {
	// TODO: confirm that all votes have been cast...
	// }

	return shim.Error("Not implemented yet")

}

// What do these parameters mean???
func (s *SmartContract) verifyZKP(userID string, xG []big.Int, r big.Int, vG []big.Int) bool {

	bitCurve := crypto.S256()
	isOnCurve := bitCurve.IsOnCurve(&xG[0], &xG[1])

	// Reference implementation is ignoring vG[2] as well
	if !bitCurve.IsOnCurve(&xG[0], &xG[1]) || !bitCurve.IsOnCurve(&vG[0], &vG[1]) {
		return false
	}

	/*
			 * Get c = H(g, g^{x}, g^{v});
		   * bytes32 b_c = sha256(msg.sender, Gx, Gy, xG, vG);
	*/
	Gx := bitCurve.Params().Gx
	Gy := bitCurve.Params().Gy
	data := append([]byte(userID)[:], append(Gx.Bytes()[:], append(Gy.Bytes()[:], append(xG[0].Bytes()[:], xG[1].Bytes()[:]...)...)...)...)
	hashBytes := sha256.Sum256(data)
	c := new(big.Int)
	c.SetBytes(hashBytes[:])

	// Get g^{r}, and g^{xc}
	rGX, rGY := bitCurve.ScalarMult(Gx, Gy, r.Bytes())
	xcGX, xcGY := bitCurve.ScalarMult(&xG[0], &xG[1], c.Bytes())

	// Add both points together
	rGxcGX, rGxcGY := bitCurve.Add(rGX, rGY, xcGX, xcGY)

	logger.Info(isOnCurve)

	// reflect.DeepEqual(*rGxcGx, vg[0])
	if rGxcGX.Cmp(&vG[0]) == 0 && rGxcGY.Cmp(&vG[1]) == 0 {
		return true
	} else {
		return false
	}
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
