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
	UserId string
	Vote   int
}

type Result struct {
	Voters int
	Votes  map[int]int
}

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
		return s.initVote(APIstub)
	case "beginSignUp":
		return s.beginSignUp(APIstub, args)
	case "finishRegistrationPhase":
		return s.finishRegistrationPhase(APIstub)
	case "submitVote":
		return s.submitVote(APIstub, args)
	case "setEligible":
		return s.setEligible(APIstub, args)
	case "register":
		return s.register(APIstub, args)
	case "computeTally":
		return s.computeTally(APIstub)
	case "question":
		return s.question(APIstub)
	default:
		return shim.Error("Invalid Smart Contract function name: " + function)
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
