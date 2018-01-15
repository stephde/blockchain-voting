package main

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

func createZKP(userId string, x *big.Int, v *big.Int, xG []*big.Int) []*big.Int {

	bitCurve := crypto.S256()

	Gx := bitCurve.Params().Gx
	Gy := bitCurve.Params().Gy

	if !bitCurve.IsOnCurve(xG[0], xG[1]) {
		// raise exception
		logger.Error("xG is not on curve")
		return nil
	}

	vGx, vGy := bitCurve.ScalarMult(Gx, Gy, v.Bytes())

	// Get c = H(g, g^{x}, g^{v});
	data := append(
		[]byte(userId)[:],
		append(Gx.Bytes()[:],
			append(Gy.Bytes()[:],
				append(xG[0].Bytes()[:],
					append(xG[1].Bytes()[:],
						append(vGx.Bytes()[:],
							vGy.Bytes()[:]...)...)...)...)...)...)
	hashBytes := sha256.Sum256(data)
	c := new(big.Int)
	c.SetBytes(hashBytes[:])

	n := bitCurve.Params().N
	xc := mulMod(x, c, n)
	r := subMod(v, xc, n)

	return []*big.Int{r, vGx, vGy, v}
}

// from secp256k1.
func generateKeyPair() (pubkey *ecdsa.PublicKey, privkey *ecdsa.PrivateKey) {

	key, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}
	return &key.PublicKey, key
}

func Test_verifyZKP(t *testing.T) {
	publicKeyECDSA, privateKeyECDSA := generateKeyPair()
	userID := "someUserId"

	xG := []*big.Int{publicKeyECDSA.X, publicKeyECDSA.Y}

	// We abuse this as a PRNG (pseudo-random number generator)
	_, vECDSA := generateKeyPair()
	v := vECDSA.D

	zkp := createZKP(userID, privateKeyECDSA.D, v, xG)

	r := zkp[0]
	vG := []*big.Int{zkp[1], zkp[2], zkp[3]}

	scc := new(SmartContract)
	assert.True(t, scc.verifyZKP(userID, xG, r, vG))
}
