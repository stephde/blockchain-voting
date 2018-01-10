package main

import (
	"crypto/sha256"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
)

// What do these parameters mean???
func (s *SmartContract) verifyZKP(userID string, xG []big.Int, r big.Int, vG []big.Int) bool {
	bitCurve := crypto.S256()

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

	// reflect.DeepEqual(*rGxcGx, vg[0])
	if rGxcGX.Cmp(&vG[0]) == 0 && rGxcGY.Cmp(&vG[1]) == 0 {
		return true
	} else {
		return false
	}
}
