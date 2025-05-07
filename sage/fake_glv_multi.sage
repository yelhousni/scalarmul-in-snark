# Fake GLV
# Find small u,v,w such that k = u/w mod r and l = v/w mod r simultaneously

load("curves.sage")

δ = 0.99  # For LLL
α = (δ-1/4)
max_bound_dic = {}
for curve in curves:
    r = curves[curve]['r']
    λ = curves[curve]['λ']
    max_bound = 0
    bound_lll = ceil((r/α)**(2/3))
    for i in range(1000):
        k = Integer(randint(1, r))
        l = Integer(randint(1, r))

        M = Matrix(ZZ, [
            [r, 0, 0],
            [0, r, 0],
            [0, 0, r],
            [k, l, 1],
        ]).LLL(delta=δ)
        sol = 0
        (u, v, w) = M[sol]
        tmp = max(abs(u), abs(v), abs(w))
        while prod(M[sol]) == 0 or tmp > bound_lll:
            sol += 1
            (u, v, w) = M[sol]
            tmp = max(abs(u), abs(v), abs(w))
        print(tmp.nbits())
        print(bound_lll.nbits())
        print()
        assert k == (u/w) % r
        assert l == (v/w) % r

        # in-circuit:
        # [u]P + [v]Q - [w]R == 0
        # where R = [k]P+[l]Q is a hint.

        # Bound size
        assert tmp <= bound_lll
        if tmp > max_bound:
            max_bound = tmp
    max_bound_dic[curve] = {
        "computed": max_bound,
        "theoretic": bound_lll
    }
