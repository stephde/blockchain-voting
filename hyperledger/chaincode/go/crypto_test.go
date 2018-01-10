package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

// func (s *SmartContract) createZKP(userID string, xG []big.Int, x big.Int, v big.Int) (big.Int, []big.Int) {
// 	bitCurve := crypto.S256()
//
// 	// Reference implementation is ignoring vG[2] as well
// 	if !bitCurve.IsOnCurve(&xG[0], &xG[1]) {
// 		return
// 	}
//
// 	/*
// 			 * Get c = H(g, g^{x}, g^{v});
// 		   * bytes32 b_c = sha256(msg.sender, Gx, Gy, xG, vG);
// 	*/
// 	Gx := bitCurve.Params().Gx
// 	Gy := bitCurve.Params().Gy
// 	data := append([]byte(userID)[:], append(Gx.Bytes()[:], append(Gy.Bytes()[:], append(xG[0].Bytes()[:], xG[1].Bytes()[:]...)...)...)...)
// 	hashBytes := sha256.Sum256(data)
// 	c := new(big.Int)
// 	c.SetBytes(hashBytes[:])
//
// 	// Get g^{r}, and g^{xc}
// 	rGX, rGY := bitCurve.ScalarMult(Gx, Gy, r.Bytes())
// 	xcGX, xcGY := bitCurve.ScalarMult(&xG[0], &xG[1], c.Bytes())
//
// 	// Add both points together
// 	rGxcGX, rGxcGY := bitCurve.Add(rGX, rGY, xcGX, xcGY)
//
// 	// reflect.DeepEqual(*rGxcGx, vg[0])
// 	if rGxcGX.Cmp(&vG[0]) == 0 && rGxcGY.Cmp(&vG[1]) == 0 {
// 		return true
// 	} else {
// 		return false
// 	}
// }

// from secp256k1.
func generateKeyPair() (pubkey, privkey []byte) {
	key, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	pubkey = elliptic.Marshal(crypto.S256(), key.X, key.Y)
	return pubkey, math.PaddedBigBytes(key.D, 32)
}

func Test_verifyZKP(t *testing.T) {
	publicKey, _ := generateKeyPair()
	userID := ""

	xGx := new(big.Int).SetBytes(publicKey[0:31])
	xGy := new(big.Int).SetBytes(publicKey[32:63])
	xG := []big.Int{*xGx, *xGy}
	r := big.NewInt(3)
	// vG := []big.Int{3, 4, 5}
	vG := []big.Int{}

	scc := new(SmartContract)
	assert.True(t, scc.verifyZKP(userID, xG, *r, vG))
}
