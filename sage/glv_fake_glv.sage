# Implementation of GLV + Fake GLV
# find small s1,s2,s3,s4 such that s = (s1+λs2)/(s3+λs4) mod r
from lll import reduction
load("curves.sage")
δ = 0.999  # For LLL
α = (δ-1/4)
max_bound_dic = {}
for curve in curves:
    r = curves[curve]['r']
    λ = curves[curve]['λ']
    max_bound = 0
    bound_lll = Integer((r/α**3)**(1/4))
    for i in range(100):
        # We consider the scalar multiplication s*G
        s = Integer(randint(r//2, r))

        basis = [
            [r, 0, 0, 0],
            # [0, r, 0, 0],
            # [0, 0, r, 0],
            # [0, 0, 0, r],
            [-λ, 1, 0, 0],
            [s, 0, 1, 0],
            [0, 0, -λ, 1],
        ]
        M = Matrix(ZZ, basis).LLL(delta=δ)
        print(M)
        print(reduction(basis, δ))
        sol = 0
        (k1, k2, k3, k4) = M[sol]
        while prod(M[sol]) == 0:
            sol += 1
            (k1, k2, k3, k4) = M[sol]
        assert (k1 + λ * k2 - s * k3 - s * λ * k4) % r == 0
        # in-circuit:
        # [k1]P + [k2]φ(P) - [k3] Q - [k4]φ(Q) == 0
        # where Q = [s]P is a hint.

        # Bound size
        tmp = max([abs(elt) for elt in M[sol]])
        assert tmp <= bound_lll
        if tmp > max_bound:
            max_bound = tmp
    max_bound_dic[curve] = {
        "computed": max_bound,
        "theoretic": bound_lll
    }
