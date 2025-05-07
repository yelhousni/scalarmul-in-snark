# Implementation of GLV + Fake GLV
# find small s1,s2,s3,s4 such that s = (s1+λs2)/(s3+λs4) mod r
load("curves.sage")

δ = 0.999  # For LLL
α = (δ-1/4)
max_bound_dic = {}
for curve in curves:
    r = curves[curve]['r']
    λ = curves[curve]['λ']
    max_bound = 0
    bound_lll = round((r**2 / α**5)**(1/6))
    for i in range(1000):
        k = randint(1, r)
        l = randint(1, r)
        M = Matrix(ZZ, [
            [r, 0, 0, 0, 0, 0],
            [0, r, 0, 0, 0, 0],
            [0, 0, r, 0, 0, 0],
            [0, 0, 0, r, 0, 0],
            [0, 0, 0, 0, r, 0],
            [0, 0, 0, 0, 0, r],
            [-λ, 1, 0, 0, 0, 0],
            [0, 0, -λ, 1, 0, 0],
            [0, 0, 0, 0, -λ, 1],
            [k, 0, l, 0, 1, 0],
        ])
        N = M.LLL()
        # find a good vector from LLL
        sol = -1
        (u, v, w, x, y, z) = N[sol]
        while y == 0:
            sol -= 1
            (u, v, w, x, y, z) = N[sol]
        assert k == ((u+v*λ)/(y+λ*z)) % r
        assert l == ((w+x*λ)/(y+λ*z)) % r
        print(max(abs(u), abs(v), abs(w), abs(x), abs(y), abs(z)).nbits())
        # in-circuit:
        # [u]P + [v]φ(P) + [w]Q + [x]φ(Q) - [y]R - [z]φ(R) == 0
        # where R = [k]P + [l]Q is a hint.
        # /!\ but this becomes slower!

        # Bound size
        tmp = max([abs(elt) for elt in N[sol]])
        print(tmp, tmp.nbits())
        print(bound_lll, tmp.nbits())
        print()
        assert tmp <= bound_lll
        if tmp > max_bound:
            max_bound = tmp
    max_bound_dic[curve] = {
        "computed": max_bound,
        "theoretic": bound_lll
    }


for curve in max_bound_dic:
    print(curve)
    print(max_bound_dic[curve]['computed'])
    print(max_bound_dic[curve]['theoretic'])
    print()
