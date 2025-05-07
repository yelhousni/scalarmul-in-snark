# Lattice-based scalar decomposition

## Description

The following code demonstrates (in [test.py](test.py)) the scalar bounds we obtained in the different decompositions (in [decompose.py](decompose.py)). We consider GLV (BN254, BLS12-381, BW6-761, secp256k1, Bandersnatch) and large discriminant (p256, p384, Jubjub) curves, as listed in [curves.py](curves.py).

* `test_simple_decompose`.
  Integer fraction decomposition of $k = x/z\bmod r$.
  Sub-scalars $x,z$ are expected to be of size $1.16\cdot r^{1/2}$.
  This can be applied for a simple scalar multiplication verification $[k]P=Q$ for a generic elliptic curve.

* `test_double_decompose`.
  Simultaneous fraction decomposition of $k_1 = x_1/z\bmod r$ and $k_2 = x_2/z\bmod r$.
  Sub-scalars $x_1,x_2,z$ are expected to be of size $1.22\cdot r^{2/3}$.
  This can be applied for a double scalar multiplication verification $[k_1]P_1+[k_2]P_2=Q$ for a generic elliptic curve.

* `test_simple_decompose_ext`.
  Quadratic extension fraction decomposition of $k = (x+\lambda y)/(z+\lambda t)\bmod r$.
  Sub-scalars $x,y,z,t$ are expected to be of size $1.25\cdot r^{1/4}$.
  This can be applied for a simple scalar multiplication verification $[k]P=Q$ for a GLV curve.

* `test_double_decompose_ext`: 
  Simultaneous quadratic extension fraction decomposition of $k_1 = (x_1+\lambda y_1)/(z+\lambda t)\bmod r$ and $k_2 = (x_2+\lambda y_2)/(z+\lambda t)\bmod r$.
  Sub-scalars $x_1,x_2,y_1,y_2,z,t$ are expected to be of size $1.28\cdot r^{1/3}$.
  This can be applied for a double scalar multiplication verification $[k_1]P_1+[k_2]P_2=Q$ for a GLV curve.

## How to use
In order to use this code, run the command:
```
sage test.py
```
It should output:
```
Simple decomposition
	p256... OK
	p384... OK
	jubjub... OK
	bn254... OK
	bls12-381... OK
	bw6-761... OK
	secp256k1... OK
	bandersnatch... OK

Double decomposition
	p256... OK
	p384... OK
	jubjub... OK
	bn254... OK
	bls12-381... OK
	bw6-761... OK
	secp256k1... OK
	bandersnatch... OK

Simple decomposition (GLV)
	bn254... OK
	bls12-381... OK
	bw6-761... OK
	secp256k1... OK
	bandersnatch... OK

Double decomposition (GLV)
	bn254... OK
	bls12-381... OK
	bw6-761... OK
	secp256k1... OK
	bandersnatch... OK
```