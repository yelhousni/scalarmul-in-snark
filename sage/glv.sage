# GLV technique using only a LLL computation
# Find s1,s2 such that s = s1 + λ s2 mod r

# NB: Usually, GLV compute the CVP of (k,0). Here, we compute it using LLL as CVP can be obtained using SVP.

load("curves.sage")

δ = 0.999  # For LLL
α = (δ-1/4)
max_bound_dic = {}
for curve in curves:
    r = curves[curve]['r']
    λ = curves[curve]['λ']
    max_bound = 0
    bound_lll = round(r**(1/2)) // round(α**(2/3))
    T = Matrix([
        [r, 0],
        [0, r],
        [-λ, 1],
    ])
    for i in range(1000):
        # We consider the scalar multiplication s*G
        s = Integer(randint(r//2, r))

        w = vector([s, 0])

        # We can find (s1,s2) by solving CVP, but we prefer SVP (LLL).
        M = Matrix(4, 3)
        M[0:3, 0:2] = T
        M[3, 0:2] = w
        M[3, 2] = floor(sqrt(r))  # or 2¹²⁸
        N = M.LLL()
        # e = (w - l1 v1 - l2 v2, 1) is in M as it is (-l1, -l2, 1) * M
        # and e is small, so we can find it in N!
        sol = -1
        (s1, s2) = N[sol][0:2]
        while (s1+s2*λ) % r != s:
            sol -= 1
            (s1, s2) = N[sol][0:2]
        assert (s1 + s2 * λ) % r == s
        tmp = max(abs(s1), abs(s2))
        assert tmp <= bound_lll
        if tmp > max_bound:
            max_bound = tmp
    max_bound_dic[curve] = {
        "computed": max_bound,
        "theoretic": bound_lll
    }
