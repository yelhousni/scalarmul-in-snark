# SageMath scripts

## Pre-requisites
Install SageMath (https://doc.sagemath.org/html/en/installation/index.html).

## Organization
The directory `sage/` contains the scripts to check the formulas presented in the paper. The following files are:
- `curves.sage`: contains the targeted elliptic curves parameters.
- `glv.sage`: corresponds to the (traditional) GLV technique using a LLL decomposition.
- `fake_glv.sage`: corresponds to the hinted simple scalar multiplication (Sec. 3.1).
- `fake_glv_multi.sage`: corresponds to the hinted double scalar multiplication (Sec. 3.2).
- `glv_fake_glv.sage`: corresponds to the hinted GLV scalar multiplication (Sec. 3.3).
- `glv_fake_glv_higher.sage`: corresponds to the hinted GLV double scalar multiplication (Sec. 3.4).
