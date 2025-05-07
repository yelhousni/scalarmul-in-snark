# Fake GLV
# Find small u,v such that s = u/v mod r

load("curves.sage")

δ = 0.99  # For LLL
α = (δ-1/4)
max_bound_dic = {}
for curve in curves:
    r = curves[curve]['r']
    λ = curves[curve]['λ']
    max_bound = 0
    bound_lll = ceil((r/α)**(1/2))
    for i in range(1000):
        # We consider the scalar multiplication s*G
        s = Integer(randint(1, r))

        M = Matrix(ZZ, [
            [r, 0],
            [0, r],
            [s, 1],
        ]).LLL(delta=δ)
        sol = 0
        (u, v) = M[sol]
        tmp = max(abs(u), abs(v))
        while prod(M[sol]) == 0 or tmp > bound_lll:
            sol += 1
            (u, v) = M[sol]
            tmp = max(abs(u), abs(v))
        assert (u/v) % r == s
        # in-circuit:
        # [u]P - [v]Q == 0
        # where Q = [s]P is a hint.

        # Bound size
        assert tmp <= bound_lll
        if tmp > max_bound:
            max_bound = tmp
    max_bound_dic[curve] = {
        "computed": max_bound,
        "theoretic": bound_lll
    }
