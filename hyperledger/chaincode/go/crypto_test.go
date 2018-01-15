package main

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

// from https://github.com/vsergeev/btckeygenie/blob/master/btckey/elliptic.go
// addMod computes z = (x + y) % p.
func addMod(x *big.Int, y *big.Int, p *big.Int) (z *big.Int) {
	z = new(big.Int).Add(x, y)
	z.Mod(z, p)
	return z
}

// subMod computes z = (x - y) % p.
func subMod(x *big.Int, y *big.Int, p *big.Int) (z *big.Int) {
	z = new(big.Int).Sub(x, y)
	z.Mod(z, p)
	return z
}

// mulMod computes z = (x * y) % p.
func mulMod(x *big.Int, y *big.Int, p *big.Int) (z *big.Int) {
	n := new(big.Int).Set(x)
	z = big.NewInt(0)

	for i := 0; i < y.BitLen(); i++ {
		if y.Bit(i) == 1 {
			z = addMod(z, n, p)
		}
		n = addMod(n, n, p)
	}

	return z
}

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
