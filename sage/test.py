from decompose import simple_decompose, double_decompose, simple_decompose_ext, double_decompose_ext
from curves import curves, curves_glv
from random import randint

reps = 200


def test_simple_decompose():
    """
    Decompose a scalar k = x/z and check that |x|,|r| ≤ (r/(δ-1/4))**(1/2)
    """
    print("Simple decomposition")
    for curve in curves:
        print("\t{}".format(curve), end='')
        r = curves[curve]['r']

        for i in range(reps):
            # Integer simple decomposition k = x/z mod r
            k = randint(1, r)
            x, z = simple_decompose(k, r)
            assert (x/z) % r == k
            assert max(abs(x), abs(z)) <= 1.16 * r**(1/2)  # (r/(δ-1/4))**(1/2)
        print("... OK")
    print()


def test_double_decompose():
    """
    Decompose simultaneously k1=x1/z and k2=x2/z and check that |x1|,|x2|,|z| ≤ (r**2/(δ-1/4)**2)**(1/3)
    """
    print("Double decomposition")
    for curve in curves:
        print("\t{}".format(curve), end='')
        r = curves[curve]['r']

        for i in range(reps):
            # Integer double decomposition k1 = x1/z mod r and k2 = x2/z mod r
            k1 = randint(1, r)
            k2 = randint(1, r)
            x1, x2, z = double_decompose(k1, k2, r)
            assert k1 == (x1/z) % r
            assert k2 == (x2/z) % r
            assert max(
                abs(x1), abs(x2), abs(z)
            ) <= 1.22 * r**(2/3)  # (r**2/(δ-1/4)**2)**(1/3)
        print("... OK")
    print()


def test_simple_decompose_ext():
    """
    Decompose a scalar k = (x+λy)/(z+λt) and check that |x|,|y|,|z|,|t| ≤ (r/(δ-1/4)**3)**(1/4)
    """
    print("Simple decomposition (GLV)")
    for curve in curves_glv:
        print("\t{}".format(curve), end='')
        r = curves[curve]['r']
        λ = curves[curve]['λ']

        for i in range(reps):
            # Quadratic extension simple decomposition k = (x+λy)/(z+λt) mod r
            k = randint(1, r)
            x, y, z, t = simple_decompose_ext(k, r, λ)
            assert (x + λ * y - k * z - k * λ * t) % r == 0
            assert max(
                abs(x), abs(y), abs(z), abs(t)
            ) <= 1.25 * r**(1/4)  # (r/(δ-1/4)**3)**(1/4)
        print("... OK")
    print()


def test_double_decompose_ext():
    """
    Decompose simultaneously k1 = (x1+λy1)/(z+λt) and k2 = (x2+λy2)/(z+λt) and check that |x1|,|y1|,|x2|,|y2|,|z|,|t| ≤ (r**2 / (δ-1/4)**5)**(1/6)
    """
    print("Double decomposition (GLV)")
    for curve in curves_glv:
        print("\t{}".format(curve), end='')
        r = curves[curve]['r']
        λ = curves[curve]['λ']

        for i in range(reps):
            # Quadratic extension double decomposition k1 = (x1+λy1)/(z+λt) mod r and k2 = (x2+λy2)/(z+λt) mod r
            k1 = randint(1, r)
            k2 = randint(1, r)
            x1, y1, x2, y2, z, t = double_decompose_ext(k1, k2, r, λ)
            assert k1 == ((x1+y1*λ)/(z+λ*t)) % r
            assert k2 == ((x2+y2*λ)/(z+λ*t)) % r
            assert max(
                abs(x1), abs(y1), abs(x2), abs(y2), abs(z), abs(t)
            ) < 1.28 * r**(1/3)  # (r**2 / (δ-1/4)**5)**(1/6)
        print("... OK")
    print()


test_simple_decompose()
test_double_decompose()
test_simple_decompose_ext()
test_double_decompose_ext()
