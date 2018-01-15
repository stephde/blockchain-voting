package main

import (
	"crypto/sha256"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
)

// What do these parameters mean???
func (s *SmartContract) verifyZKP(userID string, xG []*big.Int, r *big.Int, vG []*big.Int) bool {
	bitCurve := crypto.S256()

	logger.Info(fmt.Sprintf("%x", xG[0]))
	logger.Info(fmt.Sprintf("%x", xG[1]))

	logger.Info(bitCurve.IsOnCurve(xG[0], xG[1]))

	// Reference implementation is ignoring vG[2] as well
	if !bitCurve.IsOnCurve(xG[0], xG[1]) || !bitCurve.IsOnCurve(vG[0], vG[1]) {
		return false
	}

	logger.Info("Points are on curve")
	/*
			 * Get c = H(g, g^{x}, g^{v});
		   * bytes32 b_c = sha256(msg.sender, Gx, Gy, xG, vG);
	*/
	Gx := bitCurve.Params().Gx
	Gy := bitCurve.Params().Gy
	data := append(
		[]byte(userID)[:],
		append(Gx.Bytes()[:],
			append(Gy.Bytes()[:],
				append(xG[0].Bytes()[:],
					append(xG[1].Bytes()[:],
						append(vG[0].Bytes()[:],
							vG[1].Bytes()[:]...)...)...)...)...)...)
	hashBytes := sha256.Sum256(data)
	c := new(big.Int)
	c.SetBytes(hashBytes[:])

	// Get g^{r}, and g^{xc}
	rGX, rGY := bitCurve.ScalarMult(Gx, Gy, r.Bytes())
	xcGX, xcGY := bitCurve.ScalarMult(xG[0], xG[1], c.Bytes())

	// Add both points together
	rGxcGX, rGxcGY := bitCurve.Add(rGX, rGY, xcGX, xcGY)

	// reflect.DeepEqual(*rGxcGx, vg[0])
	return rGxcGX.Cmp(vG[0]) == 0 && rGxcGY.Cmp(vG[1]) == 0
}
