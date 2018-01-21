package main

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
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

/*
 * xG: public key
 * r: ?
 * vG: ZKP
 */
func (s *SmartContract) verifyZKP(userID string, xG *ecdsa.PublicKey, r *big.Int, vG []*big.Int) bool {
	bitCurve := crypto.S256()

	logger.Info(bitCurve.IsOnCurve(xG.X, xG.Y))

	// Reference implementation is ignoring vG[2] as well
	if !bitCurve.IsOnCurve(xG.X, xG.Y) || !bitCurve.IsOnCurve(vG[0], vG[1]) {
		return false
	}

	logger.Info("Points are on curve")
	/*
			 * Get c = H(g, g^{x}, g^{v});
		   * bytes32 b_c = sha256(msg.sender, Gx, Gy, xG, vG);
	*/
	Gx := bitCurve.Params().Gx
	Gy := bitCurve.Params().Gy

	data := Append([]byte(userID), Gx, Gy, xG.X, xG.Y, vG[0], vG[1])
	hashBytes := sha256.Sum256(data)
	c := new(big.Int)
	c.SetBytes(hashBytes[:])

	// Get g^{r}, and g^{xc}
	rGX, rGY := bitCurve.ScalarMult(Gx, Gy, r.Bytes())
	xcGX, xcGY := bitCurve.ScalarMult(xG.X, xG.Y, c.Bytes())

	// Add both points together
	rGxcGX, rGxcGY := bitCurve.Add(rGX, rGY, xcGX, xcGY)

	// reflect.DeepEqual(*rGxcGx, vg[0])
	return rGxcGX.Cmp(vG[0]) == 0 && rGxcGY.Cmp(vG[1]) == 0
}

func (s *SmartContract) verify1outOf2ZKP(
	v Voter,
	params []*big.Int,
	y *ecdsa.PublicKey,
	a1 *ecdsa.PublicKey,
	b1 *ecdsa.PublicKey,
	a2 *ecdsa.PublicKey,
	b2 *ecdsa.PublicKey) bool {
	curve := crypto.S256()

	var temp1 []*big.Int
	var temp2 []*big.Int
	var temp3 []*big.Int

	yG := v.reconstructedKey
	publicKey := v.registeredKey // xG in OpenVote

	// make sure we are only dealing with valid public keys
	if !curve.IsOnCurve(publicKey.X, publicKey.Y) ||
		!curve.IsOnCurve(yG.X, yG.Y) ||
		!curve.IsOnCurve(y.X, y.Y) ||
		!curve.IsOnCurve(a1.X, a1.Y) ||
		!curve.IsOnCurve(b1.X, b1.Y) ||
		!curve.IsOnCurve(a2.X, a2.Y) ||
		!curve.IsOnCurve(b2.X, b2.Y) {
		logger.Error("some input is not valid public key")
		return false
	}

	data := Append([]byte("")[:], publicKey.X, publicKey.Y, y.X, y.Y, a1.X, a1.Y, b1.X, b1.Y, a2.X, a2.Y, b2.X, b2.Y)
	hashBytes := sha256.Sum256(data)
	c := new(big.Int)
	c.SetBytes(hashBytes[:])

	// Does c =? d1 + d2 (mod n)
	if c != addMod(params[0], params[1], curve.Params().N) {
		return false
	}

	// a1 =? g^{r1} * x^{d1}
	temp2[0], temp2[1] = curve.ScalarMult(curve.Params().Gx, curve.Params().Gy, params[2].Bytes())
	tempX, tempY := curve.ScalarMult(publicKey.X, publicKey.Y, params[0].Bytes())
	temp3[0], temp3[1] = curve.Add(temp2[0], temp2[1], tempX, tempY)

	if a1.X != temp3[0] || a1.Y != temp3[1] {
		return false
	}

	//b1 =? h^{r1} * y^{d1} (temp = affine 'y')
	temp2[0], temp2[1] = curve.ScalarMult(yG.X, yG.Y, params[2].Bytes())
	tempX, tempY = curve.ScalarMult(y.X, y.Y, params[0].Bytes())
	temp3[0], temp3[1] = curve.Add(temp2[0], temp2[1], tempX, tempY)

	if b1.X != temp3[0] || b1.Y != temp3[1] {
		return false
	}

	//a2 =? g^{r2} * x^{d2}
	temp2[0], temp2[1] = curve.ScalarMult(curve.Params().Gx, curve.Params().Gy, params[2].Bytes())
	tempX, tempY = curve.ScalarMult(publicKey.X, publicKey.Y, params[1].Bytes())
	temp3[0], temp3[1] = curve.Add(temp2[0], temp2[1], tempX, tempY)

	if a2.X != temp3[0] || a2.Y != temp3[1] {
		return false
	}

	// Negate the 'y' co-ordinate of g
	temp1[0] = curve.Params().Gx
	temp1[1] = new(big.Int).Sub(curve.Params().P, curve.Params().Gy)

	// get 'y'
	temp3[0] = y.X
	temp3[1] = y.Y
	temp3[2] = big.NewInt(1)

	temp2[0], temp2[1] = curve.Add(temp3[0], temp3[1], temp1[0], temp1[1])
	temp1[0] = temp2[0]
	temp1[1] = temp2[1]

	// (y-g)^{d2}
	temp2[0], temp2[1] = curve.ScalarMult(temp1[0], temp1[1], params[1].Bytes())

	// Now... it is h^{r2} + temp2..
	foo, bar := curve.ScalarMult(yG.X, yG.Y, params[3].Bytes())
	temp3[0], temp3[1] = curve.Add(foo, bar, temp2[0], temp2[1])

	if b2.X != temp3[0] || b2.Y != temp3[1] {
		return false
	}

	return true
}

func (s *SmartContract) create1outof2ZKPYesVote(
	xG *ecdsa.PublicKey,
	yG *ecdsa.PublicKey,
	w *big.Int,
	r1 *big.Int,
	d1 *big.Int,
	x *big.Int) ([]*big.Int, []*big.Int) {
	// Return values
	var res []*big.Int
	var res2 []*big.Int

	// Curve
	curve := crypto.S256()

	// Curve base point
	Gx := curve.Params().Gx
	Gy := curve.Params().Gy

	var temp []*big.Int
	var temp1 []*big.Int
	var temp2 []*big.Int

	// y = h^{x} * g
	temp1[0], temp1[1] = curve.ScalarMult(yG.X, yG.Y, x.Bytes())
	temp1[0], temp1[1] = curve.Add(temp1[0], temp1[1], Gx, Gy)
	res[0] = temp1[0]
	res[1] = temp1[1]

	// a1 = g^{r1} * x^{d1}
	temp1[0], temp1[1] = curve.ScalarMult(Gx, Gy, r1.Bytes())
	temp2[0], temp2[1] = curve.ScalarMult(xG.X, xG.Y, d1.Bytes())
	temp1[0], temp1[1] = curve.Add(temp1[0], temp1[1], temp2[0], temp2[1])
	res[2] = temp1[0]
	res[3] = temp1[1]

	// b1 = h^{r1} * y^{d1} (temp = affine 'y')
	temp1[0], temp1[1] = curve.ScalarMult(yG.X, yG.Y, r1.Bytes())

	// Setting temp to 'y'
	temp[0] = res[0]
	temp[1] = res[1]

	temp2[0], temp2[1] = curve.ScalarMult(temp[0], temp[1], d1.Bytes())
	temp1[0], temp1[1] = curve.Add(temp1[0], temp1[1], temp2[0], temp2[1])
	res[4] = temp1[0]
	res[5] = temp1[1]

	// a2 = g^{w}
	temp1[0], temp1[1] = curve.ScalarMult(Gx, Gy, w.Bytes())
	res[6] = temp1[0]
	res[7] = temp1[1]

	// b2 = h^{w} (where h = g^{y})
	temp1[0], temp1[1] = curve.ScalarMult(yG.X, yG.Y, w.Bytes())
	res[8] = temp1[0]
	res[9] = temp1[1]

	return res, res2
}

func Append(slice []byte, values ...*big.Int) []byte {
	for _, r := range values {
		slice = append(slice, r.Bytes()...)
	}
	return slice
}
