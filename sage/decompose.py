from sage.all import ZZ, Matrix


def simple_decompose(k, r, δ=0.99):
    """
    Decompose k = x/z mod r
    """
    M = Matrix(ZZ, [
        [r, 0],
        [0, r],
        [k, 1],
    ]).LLL(delta=δ)
    sol = 0
    (x, z) = M[sol]
    while z == 0:
        sol += 1
        (x, z) = M[sol]
    return x, z


def double_decompose(k1, k2, r, δ=0.99):
    """
    Decompose k1 = x1/z mod r and k2 = x2/z mod r
    """
    M = Matrix(ZZ, [
        [r, 0, 0],
        [0, r, 0],
        [0, 0, r],
        [k1, k2, 1],
    ]).LLL(delta=δ)
    sol = 0
    (x1, x2, z) = M[sol]
    while z == 0:
        sol += 1
        (x1, x2, z) = M[sol]
    return (x1, x2, z)


def simple_decompose_ext(k, r, λ, δ=0.99):
    """
    Decompose k = (x+λy)/(z+λt) mod r
    """
    M = Matrix(ZZ, [
        [r, 0, 0, 0],
        [0, r, 0, 0],
        [0, 0, r, 0],
        [0, 0, 0, r],
        [-λ, 1, 0, 0],
        [k, 0, 1, 0],
        [0, 0, -λ, 1],
    ]).LLL(delta=δ)
    sol = 0
    (x, y, z, t) = M[sol]
    while (z, t) == (0, 0):
        sol += 1
        (x, y, z, t) = M[sol]
    return x, y, z, t


def double_decompose_ext(k1, k2, r, λ, δ=0.99):
    """
    Decompose k1 = (x1+λy1)/(z+λt) mod r and k2 = (x2+λy2)/(z+λt) mod r
    """
    M = Matrix(ZZ, [
        [r, 0, 0, 0, 0, 0],
        [0, r, 0, 0, 0, 0],
        [0, 0, r, 0, 0, 0],
        [0, 0, 0, r, 0, 0],
        [0, 0, 0, 0, r, 0],
        [0, 0, 0, 0, 0, r],
        [-λ, 1, 0, 0, 0, 0],
        [0, 0, -λ, 1, 0, 0],
        [k1, 0, k2, 0, 1, 0],
        [0, 0, 0, 0, -λ, 1],
    ]).LLL(delta=δ)
    sol = 0
    (x1, y1, x2, y2, z, t) = M[0]
    while (z, t) == (0, 0):
        sol += 1
        (x1, y1, x2, y2, z, t) = M[sol]
    return x1, y1, x2, y2, z, t
