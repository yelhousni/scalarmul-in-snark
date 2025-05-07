package circuits

import (
	"crypto/elliptic"
	"crypto/rand"
	"math/big"
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	bls12381 "github.com/consensys/gnark-crypto/ecc/bls12-381"
	fr_bls381 "github.com/consensys/gnark-crypto/ecc/bls12-381/fr"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	fr_bn "github.com/consensys/gnark-crypto/ecc/bn254/fr"
	bw6761 "github.com/consensys/gnark-crypto/ecc/bw6-761"
	fr_bw6761 "github.com/consensys/gnark-crypto/ecc/bw6-761/fr"
	"github.com/consensys/gnark-crypto/ecc/secp256k1"
	fp_secp "github.com/consensys/gnark-crypto/ecc/secp256k1/fp"
	fr_secp "github.com/consensys/gnark-crypto/ecc/secp256k1/fr"
	stark_curve "github.com/consensys/gnark-crypto/ecc/stark-curve"
	fr_stark "github.com/consensys/gnark-crypto/ecc/stark-curve/fr"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/algebra/algopts"
	"github.com/consensys/gnark/std/math/emulated"
	"github.com/consensys/gnark/test"
)

var testCurve = ecc.BN254

type NegTest[T, S emulated.FieldParams] struct {
	P, Q AffinePoint[T]
}

func (c *NegTest[T, S]) Define(api frontend.API) error {
	cr, err := New[T, S](api, GetCurveParams[T]())
	if err != nil {
		return err
	}
	res := cr.Neg(&c.P)
	cr.AssertIsEqual(res, &c.Q)
	return nil
}

func TestNeg(t *testing.T) {
	assert := test.NewAssert(t)
	_, g := secp256k1.Generators()
	var yn fp_secp.Element
	yn.Neg(&g.Y)
	circuit := NegTest[emulated.Secp256k1Fp, emulated.Secp256k1Fr]{}
	witness := NegTest[emulated.Secp256k1Fp, emulated.Secp256k1Fr]{
		P: AffinePoint[emulated.Secp256k1Fp]{
			X: emulated.ValueOf[emulated.Secp256k1Fp](g.X),
			Y: emulated.ValueOf[emulated.Secp256k1Fp](g.Y),
		},
		Q: AffinePoint[emulated.Secp256k1Fp]{
			X: emulated.ValueOf[emulated.Secp256k1Fp](g.X),
			Y: emulated.ValueOf[emulated.Secp256k1Fp](yn),
		},
	}
	err := test.IsSolved(&circuit, &witness, testCurve.ScalarField())
	assert.NoError(err)
}

type AddTest[T, S emulated.FieldParams] struct {
	P, Q, R AffinePoint[T]
}

func (c *AddTest[T, S]) Define(api frontend.API) error {
	cr, err := New[T, S](api, GetCurveParams[T]())
	if err != nil {
		return err
	}
	res1 := cr.add(&c.P, &c.Q)
	res2 := cr.AddUnified(&c.P, &c.Q)
	cr.AssertIsEqual(res1, &c.R)
	cr.AssertIsEqual(res2, &c.R)
	return nil
}

func TestAdd(t *testing.T) {
	assert := test.NewAssert(t)
	var dJac, aJac secp256k1.G1Jac
	g, _ := secp256k1.Generators()
	dJac.Double(&g)
	aJac.Set(&dJac).
		AddAssign(&g)
	var dAff, aAff secp256k1.G1Affine
	dAff.FromJacobian(&dJac)
	aAff.FromJacobian(&aJac)
	circuit := AddTest[emulated.Secp256k1Fp, emulated.Secp256k1Fr]{}
	witness := AddTest[emulated.Secp256k1Fp, emulated.Secp256k1Fr]{
		P: AffinePoint[emulated.Secp256k1Fp]{
			X: emulated.ValueOf[emulated.Secp256k1Fp](g.X),
			Y: emulated.ValueOf[emulated.Secp256k1Fp](g.Y),
		},
		Q: AffinePoint[emulated.Secp256k1Fp]{
			X: emulated.ValueOf[emulated.Secp256k1Fp](dAff.X),
			Y: emulated.ValueOf[emulated.Secp256k1Fp](dAff.Y),
		},
		R: AffinePoint[emulated.Secp256k1Fp]{
			X: emulated.ValueOf[emulated.Secp256k1Fp](aAff.X),
			Y: emulated.ValueOf[emulated.Secp256k1Fp](aAff.Y),
		},
	}
	err := test.IsSolved(&circuit, &witness, testCurve.ScalarField())
	assert.NoError(err)
}

type DoubleTest[T, S emulated.FieldParams] struct {
	P, Q AffinePoint[T]
}

func (c *DoubleTest[T, S]) Define(api frontend.API) error {
	cr, err := New[T, S](api, GetCurveParams[T]())
	if err != nil {
		return err
	}
	res1 := cr.double(&c.P)
	res2 := cr.AddUnified(&c.P, &c.P)
	cr.AssertIsEqual(res1, &c.Q)
	cr.AssertIsEqual(res2, &c.Q)
	return nil
}

func TestDouble(t *testing.T) {
	assert := test.NewAssert(t)
	g, _ := secp256k1.Generators()
	var dJac secp256k1.G1Jac
	dJac.Double(&g)
	var dAff secp256k1.G1Affine
	dAff.FromJacobian(&dJac)
	circuit := DoubleTest[emulated.Secp256k1Fp, emulated.Secp256k1Fr]{}
	witness := DoubleTest[emulated.Secp256k1Fp, emulated.Secp256k1Fr]{
		P: AffinePoint[emulated.Secp256k1Fp]{
			X: emulated.ValueOf[emulated.Secp256k1Fp](g.X),
			Y: emulated.ValueOf[emulated.Secp256k1Fp](g.Y),
		},
		Q: AffinePoint[emulated.Secp256k1Fp]{
			X: emulated.ValueOf[emulated.Secp256k1Fp](dAff.X),
			Y: emulated.ValueOf[emulated.Secp256k1Fp](dAff.Y),
		},
	}
	err := test.IsSolved(&circuit, &witness, testCurve.ScalarField())
	assert.NoError(err)
}

type TripleTest[T, S emulated.FieldParams] struct {
	P, Q AffinePoint[T]
}

func (c *TripleTest[T, S]) Define(api frontend.API) error {
	cr, err := New[T, S](api, GetCurveParams[T]())
	if err != nil {
		return err
	}
	res := cr.triple(&c.P)
	cr.AssertIsEqual(res, &c.Q)
	return nil
}

func TestTriple(t *testing.T) {
	assert := test.NewAssert(t)
	g, _ := secp256k1.Generators()
	var dJac secp256k1.G1Jac
	dJac.Double(&g).AddAssign(&g)
	var dAff secp256k1.G1Affine
	dAff.FromJacobian(&dJac)
	circuit := TripleTest[emulated.Secp256k1Fp, emulated.Secp256k1Fr]{}
	witness := TripleTest[emulated.Secp256k1Fp, emulated.Secp256k1Fr]{
		P: AffinePoint[emulated.Secp256k1Fp]{
			X: emulated.ValueOf[emulated.Secp256k1Fp](g.X),
			Y: emulated.ValueOf[emulated.Secp256k1Fp](g.Y),
		},
		Q: AffinePoint[emulated.Secp256k1Fp]{
			X: emulated.ValueOf[emulated.Secp256k1Fp](dAff.X),
			Y: emulated.ValueOf[emulated.Secp256k1Fp](dAff.Y),
		},
	}
	err := test.IsSolved(&circuit, &witness, testCurve.ScalarField())
	assert.NoError(err)
}

type DoubleAndAddTest[T, S emulated.FieldParams] struct {
	P, Q, R AffinePoint[T]
}

func (c *DoubleAndAddTest[T, S]) Define(api frontend.API) error {
	cr, err := New[T, S](api, GetCurveParams[T]())
	if err != nil {
		return err
	}
	res := cr.doubleAndAdd(&c.P, &c.Q)
	cr.AssertIsEqual(res, &c.R)
	return nil
}

func TestDoubleAndAdd(t *testing.T) {
	assert := test.NewAssert(t)
	var pJac, qJac, rJac secp256k1.G1Jac
	g, _ := secp256k1.Generators()
	pJac.Double(&g)
	qJac.Set(&g)
	rJac.Double(&pJac).
		AddAssign(&qJac)
	var pAff, qAff, rAff secp256k1.G1Affine
	pAff.FromJacobian(&pJac)
	qAff.FromJacobian(&qJac)
	rAff.FromJacobian(&rJac)
	circuit := DoubleAndAddTest[emulated.Secp256k1Fp, emulated.Secp256k1Fr]{}
	witness := DoubleAndAddTest[emulated.Secp256k1Fp, emulated.Secp256k1Fr]{
		P: AffinePoint[emulated.Secp256k1Fp]{
			X: emulated.ValueOf[emulated.Secp256k1Fp](pAff.X),
			Y: emulated.ValueOf[emulated.Secp256k1Fp](pAff.Y),
		},
		Q: AffinePoint[emulated.Secp256k1Fp]{
			X: emulated.ValueOf[emulated.Secp256k1Fp](qAff.X),
			Y: emulated.ValueOf[emulated.Secp256k1Fp](qAff.Y),
		},
		R: AffinePoint[emulated.Secp256k1Fp]{
			X: emulated.ValueOf[emulated.Secp256k1Fp](rAff.X),
			Y: emulated.ValueOf[emulated.Secp256k1Fp](rAff.Y),
		},
	}
	err := test.IsSolved(&circuit, &witness, testCurve.ScalarField())
	assert.NoError(err)
}

type AddUnifiedEdgeCasesTest[T, S emulated.FieldParams] struct {
	P, Q, R AffinePoint[T]
}

func (c *AddUnifiedEdgeCasesTest[T, S]) Define(api frontend.API) error {
	cr, err := New[T, S](api, GetCurveParams[T]())
	if err != nil {
		return err
	}
	res := cr.AddUnified(&c.P, &c.Q)
	cr.AssertIsEqual(res, &c.R)
	return nil
}

func TestAddUnifiedEdgeCases(t *testing.T) {
	assert := test.NewAssert(t)
	var infinity bn254.G1Affine
	_, _, g, _ := bn254.Generators()
	var r fr_bn.Element
	_, _ = r.SetRandom()
	s := new(big.Int)
	r.BigInt(s)
	var S, Sn bn254.G1Affine
	S.ScalarMultiplication(&g, s)
	Sn.Neg(&S)

	circuit := AddUnifiedEdgeCasesTest[emulated.BN254Fp, emulated.BN254Fr]{}

	// (0,0) + (0,0) == (0,0)
	witness1 := AddUnifiedEdgeCasesTest[emulated.BN254Fp, emulated.BN254Fr]{
		P: AffinePoint[emulated.BN254Fp]{
			X: emulated.ValueOf[emulated.BN254Fp](infinity.X),
			Y: emulated.ValueOf[emulated.BN254Fp](infinity.Y),
		},
		Q: AffinePoint[emulated.BN254Fp]{
			X: emulated.ValueOf[emulated.BN254Fp](infinity.X),
			Y: emulated.ValueOf[emulated.BN254Fp](infinity.Y),
		},
		R: AffinePoint[emulated.BN254Fp]{
			X: emulated.ValueOf[emulated.BN254Fp](infinity.X),
			Y: emulated.ValueOf[emulated.BN254Fp](infinity.Y),
		},
	}
	err := test.IsSolved(&circuit, &witness1, testCurve.ScalarField())
	assert.NoError(err)

	// S + (0,0) == S
	witness2 := AddUnifiedEdgeCasesTest[emulated.BN254Fp, emulated.BN254Fr]{
		P: AffinePoint[emulated.BN254Fp]{
			X: emulated.ValueOf[emulated.BN254Fp](S.X),
			Y: emulated.ValueOf[emulated.BN254Fp](S.Y),
		},
		Q: AffinePoint[emulated.BN254Fp]{
			X: emulated.ValueOf[emulated.BN254Fp](infinity.X),
			Y: emulated.ValueOf[emulated.BN254Fp](infinity.Y),
		},
		R: AffinePoint[emulated.BN254Fp]{
			X: emulated.ValueOf[emulated.BN254Fp](S.X),
			Y: emulated.ValueOf[emulated.BN254Fp](S.Y),
		},
	}
	err = test.IsSolved(&circuit, &witness2, testCurve.ScalarField())
	assert.NoError(err)

	// (0,0) + S == S
	witness3 := AddUnifiedEdgeCasesTest[emulated.BN254Fp, emulated.BN254Fr]{
		P: AffinePoint[emulated.BN254Fp]{
			X: emulated.ValueOf[emulated.BN254Fp](infinity.X),
			Y: emulated.ValueOf[emulated.BN254Fp](infinity.Y),
		},
		Q: AffinePoint[emulated.BN254Fp]{
			X: emulated.ValueOf[emulated.BN254Fp](S.X),
			Y: emulated.ValueOf[emulated.BN254Fp](S.Y),
		},
		R: AffinePoint[emulated.BN254Fp]{
			X: emulated.ValueOf[emulated.BN254Fp](S.X),
			Y: emulated.ValueOf[emulated.BN254Fp](S.Y),
		},
	}
	err = test.IsSolved(&circuit, &witness3, testCurve.ScalarField())
	assert.NoError(err)

	// S + (-S) == (0,0)
	witness4 := AddUnifiedEdgeCasesTest[emulated.BN254Fp, emulated.BN254Fr]{
		P: AffinePoint[emulated.BN254Fp]{
			X: emulated.ValueOf[emulated.BN254Fp](S.X),
			Y: emulated.ValueOf[emulated.BN254Fp](S.Y),
		},
		Q: AffinePoint[emulated.BN254Fp]{
			X: emulated.ValueOf[emulated.BN254Fp](Sn.X),
			Y: emulated.ValueOf[emulated.BN254Fp](Sn.Y),
		},
		R: AffinePoint[emulated.BN254Fp]{
			X: emulated.ValueOf[emulated.BN254Fp](infinity.X),
			Y: emulated.ValueOf[emulated.BN254Fp](infinity.Y),
		},
	}
	err = test.IsSolved(&circuit, &witness4, testCurve.ScalarField())
	assert.NoError(err)

	// (-S) + S == (0,0)
	witness5 := AddUnifiedEdgeCasesTest[emulated.BN254Fp, emulated.BN254Fr]{
		P: AffinePoint[emulated.BN254Fp]{
			X: emulated.ValueOf[emulated.BN254Fp](Sn.X),
			Y: emulated.ValueOf[emulated.BN254Fp](Sn.Y),
		},
		Q: AffinePoint[emulated.BN254Fp]{
			X: emulated.ValueOf[emulated.BN254Fp](S.X),
			Y: emulated.ValueOf[emulated.BN254Fp](S.Y),
		},
		R: AffinePoint[emulated.BN254Fp]{
			X: emulated.ValueOf[emulated.BN254Fp](infinity.X),
			Y: emulated.ValueOf[emulated.BN254Fp](infinity.Y),
		},
	}
	err = test.IsSolved(&circuit, &witness5, testCurve.ScalarField())
	assert.NoError(err)
}

type ScalarMulTest[T, S emulated.FieldParams] struct {
	P, Q AffinePoint[T]
	S    emulated.Element[S]
}

func (c *ScalarMulTest[T, S]) Define(api frontend.API) error {
	cr, err := New[T, S](api, GetCurveParams[T]())
	if err != nil {
		return err
	}
	res := cr.ScalarMul(&c.P, &c.S)
	cr.AssertIsEqual(res, &c.Q)
	return nil
}

func TestScalarMul(t *testing.T) {
	assert := test.NewAssert(t)
	_, g := secp256k1.Generators()
	var r fr_secp.Element
	_, _ = r.SetRandom()
	s := new(big.Int)
	r.BigInt(s)
	var S secp256k1.G1Affine
	S.ScalarMultiplication(&g, s)

	circuit := ScalarMulTest[emulated.Secp256k1Fp, emulated.Secp256k1Fr]{}
	witness := ScalarMulTest[emulated.Secp256k1Fp, emulated.Secp256k1Fr]{
		S: emulated.ValueOf[emulated.Secp256k1Fr](s),
		P: AffinePoint[emulated.Secp256k1Fp]{
			X: emulated.ValueOf[emulated.Secp256k1Fp](g.X),
			Y: emulated.ValueOf[emulated.Secp256k1Fp](g.Y),
		},
		Q: AffinePoint[emulated.Secp256k1Fp]{
			X: emulated.ValueOf[emulated.Secp256k1Fp](S.X),
			Y: emulated.ValueOf[emulated.Secp256k1Fp](S.Y),
		},
	}
	err := test.IsSolved(&circuit, &witness, testCurve.ScalarField())
	assert.NoError(err)
}

func TestScalarMul2(t *testing.T) {
	assert := test.NewAssert(t)
	var r fr_bn.Element
	_, _ = r.SetRandom()
	s := new(big.Int)
	r.BigInt(s)
	var res bn254.G1Affine
	_, _, gen, _ := bn254.Generators()
	res.ScalarMultiplication(&gen, s)

	circuit := ScalarMulTest[emulated.BN254Fp, emulated.BN254Fr]{}
	witness := ScalarMulTest[emulated.BN254Fp, emulated.BN254Fr]{
		S: emulated.ValueOf[emulated.BN254Fr](s),
		P: AffinePoint[emulated.BN254Fp]{
			X: emulated.ValueOf[emulated.BN254Fp](gen.X),
			Y: emulated.ValueOf[emulated.BN254Fp](gen.Y),
		},
		Q: AffinePoint[emulated.BN254Fp]{
			X: emulated.ValueOf[emulated.BN254Fp](res.X),
			Y: emulated.ValueOf[emulated.BN254Fp](res.Y),
		},
	}
	err := test.IsSolved(&circuit, &witness, testCurve.ScalarField())
	assert.NoError(err)
}

func TestScalarMul3(t *testing.T) {
	assert := test.NewAssert(t)
	var r fr_bls381.Element
	_, _ = r.SetRandom()
	s := new(big.Int)
	r.BigInt(s)
	var res bls12381.G1Affine
	_, _, gen, _ := bls12381.Generators()
	res.ScalarMultiplication(&gen, s)

	circuit := ScalarMulTest[emulated.BLS12381Fp, emulated.BLS12381Fr]{}
	witness := ScalarMulTest[emulated.BLS12381Fp, emulated.BLS12381Fr]{
		S: emulated.ValueOf[emulated.BLS12381Fr](s),
		P: AffinePoint[emulated.BLS12381Fp]{
			X: emulated.ValueOf[emulated.BLS12381Fp](gen.X),
			Y: emulated.ValueOf[emulated.BLS12381Fp](gen.Y),
		},
		Q: AffinePoint[emulated.BLS12381Fp]{
			X: emulated.ValueOf[emulated.BLS12381Fp](res.X),
			Y: emulated.ValueOf[emulated.BLS12381Fp](res.Y),
		},
	}
	err := test.IsSolved(&circuit, &witness, testCurve.ScalarField())
	assert.NoError(err)
}

func TestScalarMul4(t *testing.T) {
	assert := test.NewAssert(t)
	p256 := elliptic.P256()
	s, err := rand.Int(rand.Reader, p256.Params().N)
	assert.NoError(err)
	px, py := p256.ScalarBaseMult(s.Bytes())

	circuit := ScalarMulTest[emulated.P256Fp, emulated.P256Fr]{}
	witness := ScalarMulTest[emulated.P256Fp, emulated.P256Fr]{
		S: emulated.ValueOf[emulated.P256Fr](s),
		P: AffinePoint[emulated.P256Fp]{
			X: emulated.ValueOf[emulated.P256Fp](p256.Params().Gx),
			Y: emulated.ValueOf[emulated.P256Fp](p256.Params().Gy),
		},
		Q: AffinePoint[emulated.P256Fp]{
			X: emulated.ValueOf[emulated.P256Fp](px),
			Y: emulated.ValueOf[emulated.P256Fp](py),
		},
	}
	err = test.IsSolved(&circuit, &witness, testCurve.ScalarField())
	assert.NoError(err)
}

func TestScalarMul5(t *testing.T) {
	assert := test.NewAssert(t)
	p384 := elliptic.P384()
	s, err := rand.Int(rand.Reader, p384.Params().N)
	assert.NoError(err)
	px, py := p384.ScalarBaseMult(s.Bytes())

	circuit := ScalarMulTest[emulated.P384Fp, emulated.P384Fr]{}
	witness := ScalarMulTest[emulated.P384Fp, emulated.P384Fr]{
		S: emulated.ValueOf[emulated.P384Fr](s),
		P: AffinePoint[emulated.P384Fp]{
			X: emulated.ValueOf[emulated.P384Fp](p384.Params().Gx),
			Y: emulated.ValueOf[emulated.P384Fp](p384.Params().Gy),
		},
		Q: AffinePoint[emulated.P384Fp]{
			X: emulated.ValueOf[emulated.P384Fp](px),
			Y: emulated.ValueOf[emulated.P384Fp](py),
		},
	}
	err = test.IsSolved(&circuit, &witness, testCurve.ScalarField())
	assert.NoError(err)
}

func TestScalarMul6(t *testing.T) {
	assert := test.NewAssert(t)
	var r fr_bw6761.Element
	_, _ = r.SetRandom()
	s := new(big.Int)
	r.BigInt(s)
	var res bw6761.G1Affine
	_, _, gen, _ := bw6761.Generators()
	res.ScalarMultiplication(&gen, s)

	circuit := ScalarMulTest[emulated.BW6761Fp, emulated.BW6761Fr]{}
	witness := ScalarMulTest[emulated.BW6761Fp, emulated.BW6761Fr]{
		S: emulated.ValueOf[emulated.BW6761Fr](s),
		P: AffinePoint[emulated.BW6761Fp]{
			X: emulated.ValueOf[emulated.BW6761Fp](gen.X),
			Y: emulated.ValueOf[emulated.BW6761Fp](gen.Y),
		},
		Q: AffinePoint[emulated.BW6761Fp]{
			X: emulated.ValueOf[emulated.BW6761Fp](res.X),
			Y: emulated.ValueOf[emulated.BW6761Fp](res.Y),
		},
	}
	err := test.IsSolved(&circuit, &witness, testCurve.ScalarField())
	assert.NoError(err)
}

func TestScalarMul7(t *testing.T) {
	assert := test.NewAssert(t)
	var r fr_stark.Element
	_, _ = r.SetRandom()
	s := new(big.Int)
	r.BigInt(s)
	var res stark_curve.G1Affine
	_, gen := stark_curve.Generators()
	res.ScalarMultiplication(&gen, s)

	circuit := ScalarMulTest[emulated.STARKCurveFp, emulated.STARKCurveFr]{}
	witness := ScalarMulTest[emulated.STARKCurveFp, emulated.STARKCurveFr]{
		S: emulated.ValueOf[emulated.STARKCurveFr](s),
		P: AffinePoint[emulated.STARKCurveFp]{
			X: emulated.ValueOf[emulated.STARKCurveFp](gen.X),
			Y: emulated.ValueOf[emulated.STARKCurveFp](gen.Y),
		},
		Q: AffinePoint[emulated.STARKCurveFp]{
			X: emulated.ValueOf[emulated.STARKCurveFp](res.X),
			Y: emulated.ValueOf[emulated.STARKCurveFp](res.Y),
		},
	}
	err := test.IsSolved(&circuit, &witness, testCurve.ScalarField())
	assert.NoError(err)
}

type ScalarMulEdgeCasesTest[T, S emulated.FieldParams] struct {
	P, R AffinePoint[T]
	S    emulated.Element[S]
}

func (c *ScalarMulEdgeCasesTest[T, S]) Define(api frontend.API) error {
	cr, err := New[T, S](api, GetCurveParams[T]())
	if err != nil {
		return err
	}
	res := cr.ScalarMul(&c.P, &c.S, algopts.WithCompleteArithmetic())
	cr.AssertIsEqual(res, &c.R)
	return nil
}

func TestScalarMulEdgeCasesEdgeCases(t *testing.T) {
	assert := test.NewAssert(t)
	var infinity bn254.G1Affine
	_, _, g, _ := bn254.Generators()
	var r fr_bn.Element
	_, _ = r.SetRandom()
	s := new(big.Int)
	r.BigInt(s)
	var S bn254.G1Affine
	S.ScalarMultiplication(&g, s)

	circuit := ScalarMulEdgeCasesTest[emulated.BN254Fp, emulated.BN254Fr]{}

	// s * (0,0) == (0,0)
	witness1 := ScalarMulEdgeCasesTest[emulated.BN254Fp, emulated.BN254Fr]{
		S: emulated.ValueOf[emulated.BN254Fr](s),
		P: AffinePoint[emulated.BN254Fp]{
			X: emulated.ValueOf[emulated.BN254Fp](infinity.X),
			Y: emulated.ValueOf[emulated.BN254Fp](infinity.Y),
		},
		R: AffinePoint[emulated.BN254Fp]{
			X: emulated.ValueOf[emulated.BN254Fp](infinity.X),
			Y: emulated.ValueOf[emulated.BN254Fp](infinity.Y),
		},
	}
	err := test.IsSolved(&circuit, &witness1, testCurve.ScalarField())
	assert.NoError(err)

	// 0 * S == (0,0)
	witness2 := ScalarMulEdgeCasesTest[emulated.BN254Fp, emulated.BN254Fr]{
		S: emulated.ValueOf[emulated.BN254Fr](new(big.Int)),
		P: AffinePoint[emulated.BN254Fp]{
			X: emulated.ValueOf[emulated.BN254Fp](g.X),
			Y: emulated.ValueOf[emulated.BN254Fp](g.Y),
		},
		R: AffinePoint[emulated.BN254Fp]{
			X: emulated.ValueOf[emulated.BN254Fp](infinity.X),
			Y: emulated.ValueOf[emulated.BN254Fp](infinity.Y),
		},
	}
	err = test.IsSolved(&circuit, &witness2, testCurve.ScalarField())
	assert.NoError(err)

	// 0 * (0,0) == (0,0)
	witness3 := ScalarMulEdgeCasesTest[emulated.BN254Fp, emulated.BN254Fr]{
		S: emulated.ValueOf[emulated.BN254Fr](new(big.Int)),
		P: AffinePoint[emulated.BN254Fp]{
			X: emulated.ValueOf[emulated.BN254Fp](0),
			Y: emulated.ValueOf[emulated.BN254Fp](0),
		},
		R: AffinePoint[emulated.BN254Fp]{
			X: emulated.ValueOf[emulated.BN254Fp](infinity.X),
			Y: emulated.ValueOf[emulated.BN254Fp](infinity.Y),
		},
	}
	err = test.IsSolved(&circuit, &witness3, testCurve.ScalarField())
	assert.NoError(err)
}

type ScalarMulJoyeTest[T, S emulated.FieldParams] struct {
	P, Q AffinePoint[T]
	S    emulated.Element[S]
}

func (c *ScalarMulJoyeTest[T, S]) Define(api frontend.API) error {
	cr, err := New[T, S](api, GetCurveParams[T]())
	if err != nil {
		return err
	}
	res := cr.scalarMulJoye(&c.P, &c.S)
	cr.AssertIsEqual(res, &c.Q)
	return nil
}

func TestScalarMulJoye(t *testing.T) {
	assert := test.NewAssert(t)
	p256 := elliptic.P256()
	s, err := rand.Int(rand.Reader, p256.Params().N)
	assert.NoError(err)
	px, py := p256.ScalarBaseMult(s.Bytes())

	circuit := ScalarMulJoyeTest[emulated.P256Fp, emulated.P256Fr]{}
	witness := ScalarMulJoyeTest[emulated.P256Fp, emulated.P256Fr]{
		S: emulated.ValueOf[emulated.P256Fr](s),
		P: AffinePoint[emulated.P256Fp]{
			X: emulated.ValueOf[emulated.P256Fp](p256.Params().Gx),
			Y: emulated.ValueOf[emulated.P256Fp](p256.Params().Gy),
		},
		Q: AffinePoint[emulated.P256Fp]{
			X: emulated.ValueOf[emulated.P256Fp](px),
			Y: emulated.ValueOf[emulated.P256Fp](py),
		},
	}
	err = test.IsSolved(&circuit, &witness, testCurve.ScalarField())
	assert.NoError(err)
}

func TestScalarMulJoye2(t *testing.T) {
	assert := test.NewAssert(t)
	_, g := secp256k1.Generators()
	var r fr_secp.Element
	_, _ = r.SetRandom()
	s := new(big.Int)
	r.BigInt(s)
	var S secp256k1.G1Affine
	S.ScalarMultiplication(&g, s)

	circuit := ScalarMulJoyeTest[emulated.Secp256k1Fp, emulated.Secp256k1Fr]{}
	witness := ScalarMulJoyeTest[emulated.Secp256k1Fp, emulated.Secp256k1Fr]{
		S: emulated.ValueOf[emulated.Secp256k1Fr](s),
		P: AffinePoint[emulated.Secp256k1Fp]{
			X: emulated.ValueOf[emulated.Secp256k1Fp](g.X),
			Y: emulated.ValueOf[emulated.Secp256k1Fp](g.Y),
		},
		Q: AffinePoint[emulated.Secp256k1Fp]{
			X: emulated.ValueOf[emulated.Secp256k1Fp](S.X),
			Y: emulated.ValueOf[emulated.Secp256k1Fp](S.Y),
		},
	}
	err := test.IsSolved(&circuit, &witness, testCurve.ScalarField())
	assert.NoError(err)
}

type ScalarMulFakeGLVTest[T, S emulated.FieldParams] struct {
	Q, R AffinePoint[T]
	S    emulated.Element[S]
}

func (c *ScalarMulFakeGLVTest[T, S]) Define(api frontend.API) error {
	cr, err := New[T, S](api, GetCurveParams[T]())
	if err != nil {
		return err
	}
	res := cr.scalarMulFakeGLV(&c.Q, &c.S)
	cr.AssertIsEqual(res, &c.R)
	return nil
}

func TestScalarMulFakeGLV(t *testing.T) {
	assert := test.NewAssert(t)
	p256 := elliptic.P256()
	s, err := rand.Int(rand.Reader, p256.Params().N)
	assert.NoError(err)
	px, py := p256.ScalarBaseMult(s.Bytes())

	circuit := ScalarMulFakeGLVTest[emulated.P256Fp, emulated.P256Fr]{}
	witness := ScalarMulFakeGLVTest[emulated.P256Fp, emulated.P256Fr]{
		S: emulated.ValueOf[emulated.P256Fr](s),
		Q: AffinePoint[emulated.P256Fp]{
			X: emulated.ValueOf[emulated.P256Fp](p256.Params().Gx),
			Y: emulated.ValueOf[emulated.P256Fp](p256.Params().Gy),
		},
		R: AffinePoint[emulated.P256Fp]{
			X: emulated.ValueOf[emulated.P256Fp](px),
			Y: emulated.ValueOf[emulated.P256Fp](py),
		},
	}
	err = test.IsSolved(&circuit, &witness, testCurve.ScalarField())
	assert.NoError(err)
}

func TestScalarMulFakeGLV2(t *testing.T) {
	assert := test.NewAssert(t)
	p384 := elliptic.P384()
	s, err := rand.Int(rand.Reader, p384.Params().N)
	assert.NoError(err)
	px, py := p384.ScalarBaseMult(s.Bytes())

	circuit := ScalarMulFakeGLVTest[emulated.P384Fp, emulated.P384Fr]{}
	witness := ScalarMulFakeGLVTest[emulated.P384Fp, emulated.P384Fr]{
		S: emulated.ValueOf[emulated.P384Fr](s),
		Q: AffinePoint[emulated.P384Fp]{
			X: emulated.ValueOf[emulated.P384Fp](p384.Params().Gx),
			Y: emulated.ValueOf[emulated.P384Fp](p384.Params().Gy),
		},
		R: AffinePoint[emulated.P384Fp]{
			X: emulated.ValueOf[emulated.P384Fp](px),
			Y: emulated.ValueOf[emulated.P384Fp](py),
		},
	}
	err = test.IsSolved(&circuit, &witness, testCurve.ScalarField())
	assert.NoError(err)
}

func TestScalarMulFakeGLV3(t *testing.T) {
	assert := test.NewAssert(t)
	_, g := stark_curve.Generators()
	var r fr_stark.Element
	_, _ = r.SetRandom()
	s := new(big.Int)
	r.BigInt(s)
	var S stark_curve.G1Affine
	S.ScalarMultiplication(&g, s)

	circuit := ScalarMulFakeGLVTest[emulated.STARKCurveFp, emulated.STARKCurveFr]{}
	witness := ScalarMulFakeGLVTest[emulated.STARKCurveFp, emulated.STARKCurveFr]{
		S: emulated.ValueOf[emulated.STARKCurveFr](s),
		Q: AffinePoint[emulated.STARKCurveFp]{
			X: emulated.ValueOf[emulated.STARKCurveFp](g.X),
			Y: emulated.ValueOf[emulated.STARKCurveFp](g.Y),
		},
		R: AffinePoint[emulated.STARKCurveFp]{
			X: emulated.ValueOf[emulated.STARKCurveFp](S.X),
			Y: emulated.ValueOf[emulated.STARKCurveFp](S.Y),
		},
	}
	err := test.IsSolved(&circuit, &witness, testCurve.ScalarField())
	assert.NoError(err)
}

type ScalarMulFakeGLVEdgeCasesTest[T, S emulated.FieldParams] struct {
	P, R AffinePoint[T]
	S    emulated.Element[S]
}

func (c *ScalarMulFakeGLVEdgeCasesTest[T, S]) Define(api frontend.API) error {
	cr, err := New[T, S](api, GetCurveParams[T]())
	if err != nil {
		return err
	}
	res := cr.scalarMulFakeGLV(&c.P, &c.S, algopts.WithCompleteArithmetic())
	cr.AssertIsEqual(res, &c.R)
	return nil
}

func TestScalarMulFakeGLVEdgeCasesEdgeCases(t *testing.T) {
	assert := test.NewAssert(t)
	p256 := elliptic.P256()
	s, err := rand.Int(rand.Reader, p256.Params().N)
	assert.NoError(err)
	px, py := p256.ScalarBaseMult(s.Bytes())
	_, _ = p256.ScalarMult(px, py, s.Bytes())

	circuit := ScalarMulFakeGLVEdgeCasesTest[emulated.P256Fp, emulated.P256Fr]{}

	// s * (0,0) == (0,0)
	witness1 := ScalarMulFakeGLVEdgeCasesTest[emulated.P256Fp, emulated.P256Fr]{
		S: emulated.ValueOf[emulated.P256Fr](s),
		P: AffinePoint[emulated.P256Fp]{
			X: emulated.ValueOf[emulated.P256Fp](0),
			Y: emulated.ValueOf[emulated.P256Fp](0),
		},
		R: AffinePoint[emulated.P256Fp]{
			X: emulated.ValueOf[emulated.P256Fp](0),
			Y: emulated.ValueOf[emulated.P256Fp](0),
		},
	}
	err = test.IsSolved(&circuit, &witness1, testCurve.ScalarField())
	assert.NoError(err)

	// 0 * P == (0,0)
	witness2 := ScalarMulFakeGLVEdgeCasesTest[emulated.P256Fp, emulated.P256Fr]{
		S: emulated.ValueOf[emulated.P256Fr](new(big.Int)),
		P: AffinePoint[emulated.P256Fp]{
			X: emulated.ValueOf[emulated.P256Fp](px),
			Y: emulated.ValueOf[emulated.P256Fp](py),
		},
		R: AffinePoint[emulated.P256Fp]{
			X: emulated.ValueOf[emulated.P256Fp](0),
			Y: emulated.ValueOf[emulated.P256Fp](0),
		},
	}
	err = test.IsSolved(&circuit, &witness2, testCurve.ScalarField())
	assert.NoError(err)

	// 0 * (0,0) == (0,0)
	witness3 := ScalarMulFakeGLVEdgeCasesTest[emulated.P256Fp, emulated.P256Fr]{
		S: emulated.ValueOf[emulated.P256Fr](new(big.Int)),
		P: AffinePoint[emulated.P256Fp]{
			X: emulated.ValueOf[emulated.P256Fp](0),
			Y: emulated.ValueOf[emulated.P256Fp](0),
		},
		R: AffinePoint[emulated.P256Fp]{
			X: emulated.ValueOf[emulated.P256Fp](0),
			Y: emulated.ValueOf[emulated.P256Fp](0),
		},
	}
	err = test.IsSolved(&circuit, &witness3, testCurve.ScalarField())
	assert.NoError(err)
}

func TestScalarMulFakeGLVEdgeCasesEdgeCases2(t *testing.T) {
	assert := test.NewAssert(t)
	p384 := elliptic.P384()
	s, err := rand.Int(rand.Reader, p384.Params().N)
	assert.NoError(err)
	px, py := p384.ScalarBaseMult(s.Bytes())
	_, _ = p384.ScalarMult(px, py, s.Bytes())

	circuit := ScalarMulFakeGLVEdgeCasesTest[emulated.P384Fp, emulated.P384Fr]{}

	// s * (0,0) == (0,0)
	witness1 := ScalarMulFakeGLVEdgeCasesTest[emulated.P384Fp, emulated.P384Fr]{
		S: emulated.ValueOf[emulated.P384Fr](s),
		P: AffinePoint[emulated.P384Fp]{
			X: emulated.ValueOf[emulated.P384Fp](0),
			Y: emulated.ValueOf[emulated.P384Fp](0),
		},
		R: AffinePoint[emulated.P384Fp]{
			X: emulated.ValueOf[emulated.P384Fp](0),
			Y: emulated.ValueOf[emulated.P384Fp](0),
		},
	}
	err = test.IsSolved(&circuit, &witness1, testCurve.ScalarField())
	assert.NoError(err)

	// 0 * P == (0,0)
	witness2 := ScalarMulFakeGLVEdgeCasesTest[emulated.P384Fp, emulated.P384Fr]{
		S: emulated.ValueOf[emulated.P384Fr](new(big.Int)),
		P: AffinePoint[emulated.P384Fp]{
			X: emulated.ValueOf[emulated.P384Fp](px),
			Y: emulated.ValueOf[emulated.P384Fp](py),
		},
		R: AffinePoint[emulated.P384Fp]{
			X: emulated.ValueOf[emulated.P384Fp](0),
			Y: emulated.ValueOf[emulated.P384Fp](0),
		},
	}
	err = test.IsSolved(&circuit, &witness2, testCurve.ScalarField())
	assert.NoError(err)

	// 0 * (0,0) == (0,0)
	witness3 := ScalarMulFakeGLVEdgeCasesTest[emulated.P384Fp, emulated.P384Fr]{
		S: emulated.ValueOf[emulated.P384Fr](new(big.Int)),
		P: AffinePoint[emulated.P384Fp]{
			X: emulated.ValueOf[emulated.P384Fp](0),
			Y: emulated.ValueOf[emulated.P384Fp](0),
		},
		R: AffinePoint[emulated.P384Fp]{
			X: emulated.ValueOf[emulated.P384Fp](0),
			Y: emulated.ValueOf[emulated.P384Fp](0),
		},
	}
	err = test.IsSolved(&circuit, &witness3, testCurve.ScalarField())
	assert.NoError(err)
}

func TestScalarMulFakeGLVEdgeCasesEdgeCases3(t *testing.T) {
	assert := test.NewAssert(t)
	_, g := stark_curve.Generators()
	var r fr_stark.Element
	_, _ = r.SetRandom()
	s := new(big.Int)
	r.BigInt(s)
	var S stark_curve.G1Affine
	S.ScalarMultiplication(&g, s)

	circuit := ScalarMulFakeGLVEdgeCasesTest[emulated.STARKCurveFp, emulated.STARKCurveFr]{}

	// s * (0,0) == (0,0)
	witness1 := ScalarMulFakeGLVEdgeCasesTest[emulated.STARKCurveFp, emulated.STARKCurveFr]{
		S: emulated.ValueOf[emulated.STARKCurveFr](s),
		P: AffinePoint[emulated.STARKCurveFp]{
			X: emulated.ValueOf[emulated.STARKCurveFp](0),
			Y: emulated.ValueOf[emulated.STARKCurveFp](0),
		},
		R: AffinePoint[emulated.STARKCurveFp]{
			X: emulated.ValueOf[emulated.STARKCurveFp](0),
			Y: emulated.ValueOf[emulated.STARKCurveFp](0),
		},
	}
	err := test.IsSolved(&circuit, &witness1, testCurve.ScalarField())
	assert.NoError(err)

	// 0 * P == (0,0)
	witness2 := ScalarMulFakeGLVEdgeCasesTest[emulated.STARKCurveFp, emulated.STARKCurveFr]{
		S: emulated.ValueOf[emulated.STARKCurveFr](new(big.Int)),
		P: AffinePoint[emulated.STARKCurveFp]{
			X: emulated.ValueOf[emulated.STARKCurveFp](S.X),
			Y: emulated.ValueOf[emulated.STARKCurveFp](S.X),
		},
		R: AffinePoint[emulated.STARKCurveFp]{
			X: emulated.ValueOf[emulated.STARKCurveFp](0),
			Y: emulated.ValueOf[emulated.STARKCurveFp](0),
		},
	}
	err = test.IsSolved(&circuit, &witness2, testCurve.ScalarField())
	assert.NoError(err)

	// 0 * (0,0) == (0,0)
	witness3 := ScalarMulFakeGLVEdgeCasesTest[emulated.STARKCurveFp, emulated.STARKCurveFr]{
		S: emulated.ValueOf[emulated.STARKCurveFr](new(big.Int)),
		P: AffinePoint[emulated.STARKCurveFp]{
			X: emulated.ValueOf[emulated.STARKCurveFp](0),
			Y: emulated.ValueOf[emulated.STARKCurveFp](0),
		},
		R: AffinePoint[emulated.STARKCurveFp]{
			X: emulated.ValueOf[emulated.STARKCurveFp](0),
			Y: emulated.ValueOf[emulated.STARKCurveFp](0),
		},
	}
	err = test.IsSolved(&circuit, &witness3, testCurve.ScalarField())
	assert.NoError(err)
}

type ScalarMulGLVAndFakeGLVTest[T, S emulated.FieldParams] struct {
	Q, R AffinePoint[T]
	S    emulated.Element[S]
}

func (c *ScalarMulGLVAndFakeGLVTest[T, S]) Define(api frontend.API) error {
	cr, err := New[T, S](api, GetCurveParams[T]())
	if err != nil {
		return err
	}
	res1 := cr.scalarMulGLVAndFakeGLV(&c.Q, &c.S)
	res2 := cr.scalarMulGLVAndFakeGLV(&c.Q, &c.S, algopts.WithCompleteArithmetic())
	cr.AssertIsEqual(res1, &c.R)
	cr.AssertIsEqual(res2, &c.R)
	return nil
}

func TestScalarMulGLVAndFakeGLV(t *testing.T) {
	assert := test.NewAssert(t)
	_, g := secp256k1.Generators()
	var r fr_secp.Element
	_, _ = r.SetRandom()
	s := new(big.Int)
	r.BigInt(s)
	var S secp256k1.G1Affine
	S.ScalarMultiplication(&g, s)

	circuit := ScalarMulGLVAndFakeGLVTest[emulated.Secp256k1Fp, emulated.Secp256k1Fr]{}
	witness := ScalarMulGLVAndFakeGLVTest[emulated.Secp256k1Fp, emulated.Secp256k1Fr]{
		S: emulated.ValueOf[emulated.Secp256k1Fr](s),
		Q: AffinePoint[emulated.Secp256k1Fp]{
			X: emulated.ValueOf[emulated.Secp256k1Fp](g.X),
			Y: emulated.ValueOf[emulated.Secp256k1Fp](g.Y),
		},
		R: AffinePoint[emulated.Secp256k1Fp]{
			X: emulated.ValueOf[emulated.Secp256k1Fp](S.X),
			Y: emulated.ValueOf[emulated.Secp256k1Fp](S.Y),
		},
	}
	err := test.IsSolved(&circuit, &witness, testCurve.ScalarField())
	assert.NoError(err)
}

type ScalarMulGLVAndFakeGLVEdgeCasesTest[T, S emulated.FieldParams] struct {
	P, R AffinePoint[T]
	S    emulated.Element[S]
}

func (c *ScalarMulGLVAndFakeGLVEdgeCasesTest[T, S]) Define(api frontend.API) error {
	cr, err := New[T, S](api, GetCurveParams[T]())
	if err != nil {
		return err
	}
	res := cr.scalarMulGLVAndFakeGLV(&c.P, &c.S, algopts.WithCompleteArithmetic())
	cr.AssertIsEqual(res, &c.R)
	return nil
}

func TestScalarMulGLVAndFakeGLVEdgeCasesEdgeCases(t *testing.T) {
	assert := test.NewAssert(t)
	_, g := secp256k1.Generators()
	var r fr_secp.Element
	_, _ = r.SetRandom()
	s := new(big.Int)
	r.BigInt(s)
	var S secp256k1.G1Affine
	S.ScalarMultiplication(&g, s)

	circuit := ScalarMulGLVAndFakeGLVEdgeCasesTest[emulated.Secp256k1Fp, emulated.Secp256k1Fr]{}

	// s * (0,0) == (0,0)
	witness1 := ScalarMulGLVAndFakeGLVEdgeCasesTest[emulated.Secp256k1Fp, emulated.Secp256k1Fr]{
		S: emulated.ValueOf[emulated.Secp256k1Fr](s),
		P: AffinePoint[emulated.Secp256k1Fp]{
			X: emulated.ValueOf[emulated.Secp256k1Fp](0),
			Y: emulated.ValueOf[emulated.Secp256k1Fp](0),
		},
		R: AffinePoint[emulated.Secp256k1Fp]{
			X: emulated.ValueOf[emulated.Secp256k1Fp](0),
			Y: emulated.ValueOf[emulated.Secp256k1Fp](0),
		},
	}
	err := test.IsSolved(&circuit, &witness1, testCurve.ScalarField())
	assert.NoError(err)

	// 0 * P == (0,0)
	witness2 := ScalarMulGLVAndFakeGLVEdgeCasesTest[emulated.Secp256k1Fp, emulated.Secp256k1Fr]{
		S: emulated.ValueOf[emulated.Secp256k1Fr](new(big.Int)),
		P: AffinePoint[emulated.Secp256k1Fp]{
			X: emulated.ValueOf[emulated.Secp256k1Fp](g.X),
			Y: emulated.ValueOf[emulated.Secp256k1Fp](g.Y),
		},
		R: AffinePoint[emulated.Secp256k1Fp]{
			X: emulated.ValueOf[emulated.Secp256k1Fp](0),
			Y: emulated.ValueOf[emulated.Secp256k1Fp](0),
		},
	}
	err = test.IsSolved(&circuit, &witness2, testCurve.ScalarField())
	assert.NoError(err)

	// 0 * (0,0) == (0,0)
	witness3 := ScalarMulGLVAndFakeGLVEdgeCasesTest[emulated.Secp256k1Fp, emulated.Secp256k1Fr]{
		S: emulated.ValueOf[emulated.Secp256k1Fr](new(big.Int)),
		P: AffinePoint[emulated.Secp256k1Fp]{
			X: emulated.ValueOf[emulated.Secp256k1Fp](0),
			Y: emulated.ValueOf[emulated.Secp256k1Fp](0),
		},
		R: AffinePoint[emulated.Secp256k1Fp]{
			X: emulated.ValueOf[emulated.Secp256k1Fp](0),
			Y: emulated.ValueOf[emulated.Secp256k1Fp](0),
		},
	}
	err = test.IsSolved(&circuit, &witness3, testCurve.ScalarField())
	assert.NoError(err)

	// 1 * P == P
	witness4 := ScalarMulGLVAndFakeGLVEdgeCasesTest[emulated.Secp256k1Fp, emulated.Secp256k1Fr]{
		S: emulated.ValueOf[emulated.Secp256k1Fr](big.NewInt(1)),
		P: AffinePoint[emulated.Secp256k1Fp]{
			X: emulated.ValueOf[emulated.Secp256k1Fp](g.X),
			Y: emulated.ValueOf[emulated.Secp256k1Fp](g.Y),
		},
		R: AffinePoint[emulated.Secp256k1Fp]{
			X: emulated.ValueOf[emulated.Secp256k1Fp](g.X),
			Y: emulated.ValueOf[emulated.Secp256k1Fp](g.Y),
		},
	}
	err = test.IsSolved(&circuit, &witness4, testCurve.ScalarField())
	assert.NoError(err)

	// -1 * P == -P
	witness5 := ScalarMulGLVAndFakeGLVEdgeCasesTest[emulated.Secp256k1Fp, emulated.Secp256k1Fr]{
		S: emulated.ValueOf[emulated.Secp256k1Fr](big.NewInt(-1)),
		P: AffinePoint[emulated.Secp256k1Fp]{
			X: emulated.ValueOf[emulated.Secp256k1Fp](g.X),
			Y: emulated.ValueOf[emulated.Secp256k1Fp](g.Y),
		},
		R: AffinePoint[emulated.Secp256k1Fp]{
			X: emulated.ValueOf[emulated.Secp256k1Fp](g.X),
			Y: emulated.ValueOf[emulated.Secp256k1Fp](g.Y.Neg(&g.Y)),
		},
	}
	err = test.IsSolved(&circuit, &witness5, testCurve.ScalarField())
	assert.NoError(err)
}

func TestScalarMulGLVAndFakeGLVEdgeCasesEdgeCases2(t *testing.T) {
	assert := test.NewAssert(t)
	_, _, g, _ := bn254.Generators()
	var r fr_bn.Element
	_, _ = r.SetRandom()
	s := new(big.Int)
	r.BigInt(s)
	var S bn254.G1Affine
	S.ScalarMultiplication(&g, s)

	circuit := ScalarMulGLVAndFakeGLVEdgeCasesTest[emulated.BN254Fp, emulated.BN254Fr]{}

	// s * (0,0) == (0,0)
	witness1 := ScalarMulGLVAndFakeGLVEdgeCasesTest[emulated.BN254Fp, emulated.BN254Fr]{
		S: emulated.ValueOf[emulated.BN254Fr](s),
		P: AffinePoint[emulated.BN254Fp]{
			X: emulated.ValueOf[emulated.BN254Fp](0),
			Y: emulated.ValueOf[emulated.BN254Fp](0),
		},
		R: AffinePoint[emulated.BN254Fp]{
			X: emulated.ValueOf[emulated.BN254Fp](0),
			Y: emulated.ValueOf[emulated.BN254Fp](0),
		},
	}
	err := test.IsSolved(&circuit, &witness1, testCurve.ScalarField())
	assert.NoError(err)

	// 0 * P == (0,0)
	witness2 := ScalarMulGLVAndFakeGLVEdgeCasesTest[emulated.BN254Fp, emulated.BN254Fr]{
		S: emulated.ValueOf[emulated.BN254Fr](new(big.Int)),
		P: AffinePoint[emulated.BN254Fp]{
			X: emulated.ValueOf[emulated.BN254Fp](g.X),
			Y: emulated.ValueOf[emulated.BN254Fp](g.Y),
		},
		R: AffinePoint[emulated.BN254Fp]{
			X: emulated.ValueOf[emulated.BN254Fp](0),
			Y: emulated.ValueOf[emulated.BN254Fp](0),
		},
	}
	err = test.IsSolved(&circuit, &witness2, testCurve.ScalarField())
	assert.NoError(err)

	// 0 * (0,0) == (0,0)
	witness3 := ScalarMulGLVAndFakeGLVEdgeCasesTest[emulated.BN254Fp, emulated.BN254Fr]{
		S: emulated.ValueOf[emulated.BN254Fr](new(big.Int)),
		P: AffinePoint[emulated.BN254Fp]{
			X: emulated.ValueOf[emulated.BN254Fp](0),
			Y: emulated.ValueOf[emulated.BN254Fp](0),
		},
		R: AffinePoint[emulated.BN254Fp]{
			X: emulated.ValueOf[emulated.BN254Fp](0),
			Y: emulated.ValueOf[emulated.BN254Fp](0),
		},
	}
	err = test.IsSolved(&circuit, &witness3, testCurve.ScalarField())
	assert.NoError(err)

	// 1 * P == P
	witness4 := ScalarMulGLVAndFakeGLVEdgeCasesTest[emulated.BN254Fp, emulated.BN254Fr]{
		S: emulated.ValueOf[emulated.BN254Fr](big.NewInt(1)),
		P: AffinePoint[emulated.BN254Fp]{
			X: emulated.ValueOf[emulated.BN254Fp](g.X),
			Y: emulated.ValueOf[emulated.BN254Fp](g.Y),
		},
		R: AffinePoint[emulated.BN254Fp]{
			X: emulated.ValueOf[emulated.BN254Fp](g.X),
			Y: emulated.ValueOf[emulated.BN254Fp](g.Y),
		},
	}
	err = test.IsSolved(&circuit, &witness4, testCurve.ScalarField())
	assert.NoError(err)

	// -1 * P == -P
	witness5 := ScalarMulGLVAndFakeGLVEdgeCasesTest[emulated.BN254Fp, emulated.BN254Fr]{
		S: emulated.ValueOf[emulated.BN254Fr](big.NewInt(-1)),
		P: AffinePoint[emulated.BN254Fp]{
			X: emulated.ValueOf[emulated.BN254Fp](g.X),
			Y: emulated.ValueOf[emulated.BN254Fp](g.Y),
		},
		R: AffinePoint[emulated.BN254Fp]{
			X: emulated.ValueOf[emulated.BN254Fp](g.X),
			Y: emulated.ValueOf[emulated.BN254Fp](g.Y.Neg(&g.Y)),
		},
	}
	err = test.IsSolved(&circuit, &witness5, testCurve.ScalarField())
	assert.NoError(err)
}
