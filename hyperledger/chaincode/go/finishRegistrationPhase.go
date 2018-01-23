package main

import (
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

func (s *SmartContract) finishRegistrationPhase(stub shim.ChaincodeStubInterface) sc.Response {
	if !s.inState(stub, SIGNUP) {
		return shim.Error("Wrong state")
	}

	var totalRegistered int
	GetState(stub, "totalRegistered", &totalRegistered)
	if totalRegistered < 3 {
		return shim.Error("Too few voters registered, need at least 3")
	}

	var voters []Voter
	GetState(stub, "voters", &voters)

	voters = s.reconstructKeys(totalRegistered, voters)

	// We have computed each voter's special voting key.
	// Now we either enter the commitment phase (option) or voting phase.
	s.transitionToState(stub, VOTE)
	PutState(stub, "voters", voters)

	return shim.Success([]byte("Success"))
}

func (s *SmartContract) reconstructKeys(totalRegistered int, voters []Voter) []Voter {
	curve := crypto.S256()
	pp := curve.Params().P
	var temp [2]*big.Int
	var yG [2]*big.Int
	var beforei [2]*big.Int
	var afteri [2]*big.Int

	logger.Info("Total registered: ", totalRegistered)
	logger.Info("Number of voters: ", len(voters))

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
	logger.Info(voters[0])
	// Step 2 is to add to beforei, and subtract from afteri.
	for i := 1; i < totalRegistered; i++ {

		if i == 1 {
			beforei[0] = voters[0].registeredKey.X
			beforei[1] = voters[0].registeredKey.Y
		} else {
			beforei[0], beforei[1] = curve.Add(beforei[0], beforei[1], voters[i-1].registeredKey.X, voters[i-1].registeredKey.Y)
		}

		logger.Info("After i==1")

		// If we have reached the end... just store beforei
		// Otherwise, we need to compute a key.
		// Counting from 0 to n-1...
		voter := &voters[i]

		logger.Info("After i==1")

		if i == (totalRegistered - 1) {
			(*voter).reconstructedKey.X = beforei[0]
			(*voter).reconstructedKey.Y = beforei[1]
		} else {
			// Subtract 'i' from afteri
			temp[0], temp[1] = (*voter).registeredKey.X, big.NewInt(0).Sub(pp, (*voter).registeredKey.Y)

			// Grab negation of afteri (seems like it did not seem to work with Jacob co-ordinates)
			afteri[0], afteri[1] = curve.Add(afteri[0], afteri[1], temp[0], temp[1])

			temp[0], temp[1] = afteri[0], big.NewInt(0).Sub(pp, afteri[i])

			// Now we do beforei - afteri...
			yG[0], yG[1] = curve.Add(beforei[0], beforei[1], temp[0], temp[1])

			(*voter).reconstructedKey.X, (*voter).reconstructedKey.Y = yG[0], yG[1]
		}

		logger.Info("After if else")
	}

	return voters
}
