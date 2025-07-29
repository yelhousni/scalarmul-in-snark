# SNARK circuits

## Pre-requisites
Install Golang (https://go.dev/doc/install). This code was tested with the last 2 major releases of Go (currently 1.19 and 1.20).

## Organization
The directory `circuits/` contains the SNARK circuits written in gnark framework. The following files are:
- `points.go`: contains the circuits for points arithmetic. In particular, `scalarMulFakeGLV` corresponds to 3.1 and `scalarMulGLVAndFakeGLV` to 3.3.
- `points_test.go`: contains the tests for the circuits
- `hints.go` corresponds to the out-circuit hints. In particular, `decomposeScalarG1Subscalars` corresponds to half-GCD in Z and `halfGCDEisenstein` to a half-GCD in Eisentein integers Z[j].
- `eisenstein.go`: contains the Eisenstein integers arithmetic out-circuit.
- `eisenstein_test.go`: contains tests of Eisenstein integers arithmetic.
- `params.go` and `params_compute.go`: contain the targeted elliptic curves parameters.
- `utils.go`: contains some utility methods.

All tests are run using `go test -v .` command. For a particular test, e.g. hinted GLV scalar multiplication, run `go test -v -run TestScalarMulGLVAndFakeGLV`.
