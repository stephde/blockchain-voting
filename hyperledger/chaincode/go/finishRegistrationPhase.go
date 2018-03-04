package main

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

func (s *SmartContract) finishRegistrationPhase(stub shim.ChaincodeStubInterface) sc.Response {
	if !s.inState(stub, SIGNUP) {
		return shim.Error("Wrong state, expected SIGNUP")
	}

	// Retrieve all registrations
	name := "register"
	deltaResultsIterator, deltaErr := stub.GetStateByPartialCompositeKey("varName~userID~txID", []string{name})
	if deltaErr != nil {
		return shim.Error(fmt.Sprintf("Could not retrieve value for %s: %s", name, deltaErr.Error()))
	}
	defer deltaResultsIterator.Close()

	// Check the variable existed
	if !deltaResultsIterator.HasNext() {
		return shim.Error(fmt.Sprintf("No variable by the name %s exists", name))
	}

	// Calculate personal voting keys
	voters := s.reconstructKeys(deltaResultsIterator)

	// registered := map[string]struct{}{}
	// var i int
	// for i = 0; deltaResultsIterator.HasNext(); i++ {
	// 	// Get the next row
	// 	responseRange, nextErr := deltaResultsIterator.Next()
	// 	if nextErr != nil {
	// 		return shim.Error(nextErr.Error())
	// 	}
	//
	// 	// Split the composite key into its component parts
	// 	_, keyParts, splitKeyErr := stub.SplitCompositeKey(responseRange.Key)
	// 	if splitKeyErr != nil {
	// 		return shim.Error(splitKeyErr.Error())
	// 	}
	//
	// 	// Retrieve the userID
	// 	userID := keyParts[1]
	// 	registered[userID] = struct{}{}
	// }

	if len(voters) < 3 {
		// Legacy from Anonymous Voting Protocol, but makes sense since for two voters, each would know what the other one has voted
		return shim.Error("Too few voters registered, need at least 3")
	}
	PutState(stub, "totalRegistered", i)
	// Store slice here, since it won't be updated anymore
	PutState(stub, "registered", voters)

	// Now we either enter the voting phase.
	s.transitionToState(stub, VOTE)

	return shim.Success(nil)
}

func (s *SmartContract) reconstructKeys(voterIterator shim.StateQueryIteratorInterface) []Voter {
	curve := crypto.S256()
	pp := curve.Params().P
	var temp [2]*big.Int
	var yG [2]*big.Int
	var beforei [2]*big.Int
	var afteri [2]*big.Int

	voters := []Voter{}
	// Step 1 is to compute the index 1 reconstructed key
	voterOne := voters[1]
	voterOneRegisteredKey := voterOne.registeredKey
	afteri[0] = voterOneRegisteredKey.X
	afteri[1] = voterOneRegisteredKey.Y

	for i := 2; i < totalRegistered; i++ {
		voter := voters[i]
		registeredKey := voter.registeredKey
		afteri[0], afteri[1] = curve.Add(afteri[0], afteri[1], registeredKey.X, registeredKey.Y)
	}

	voterZero := voters[0]

	voterZero.reconstructedKey.X = afteri[0]
	voterZero.reconstructedKey.Y = big.NewInt(0).Sub(pp, afteri[1])

	voters[0] = voterZero

	// Step 2 is to add to beforei, and subtract from afteri.
	for i := 1; i < totalRegistered; i++ {
		if i == 1 {
			beforei[0] = voters[0].registeredKey.X
			beforei[1] = voters[0].registeredKey.Y
		} else {
			beforei[0], beforei[1] = curve.Add(beforei[0], beforei[1], voters[i-1].registeredKey.X, voters[i-1].registeredKey.Y)
		}

		// If we have reached the end... just store beforei
		// Otherwise, we need to compute a key.
		// Counting from 0 to n-1...
		voter := &voters[i]

		if i == (totalRegistered - 1) {
			(*voter).reconstructedKey.X = beforei[0]
			(*voter).reconstructedKey.Y = beforei[1]
		} else {
			// Subtract 'i' from afteri
			temp[0], temp[1] = (*voter).registeredKey.X, big.NewInt(0).Sub(pp, (*voter).registeredKey.Y)

			// Grab negation of afteri (seems like it did not seem to work with Jacob co-ordinates)
			afteri[0], afteri[1] = curve.Add(afteri[0], afteri[1], temp[0], temp[1])

			temp[0], temp[1] = afteri[0], big.NewInt(0).Sub(pp, afteri[1])

			// Now we do beforei - afteri...
			yG[0], yG[1] = curve.Add(beforei[0], beforei[1], temp[0], temp[1])

			(*voter).reconstructedKey.X, (*voter).reconstructedKey.Y = yG[0], yG[1]
		}
	}

	return voters
}
