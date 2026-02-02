# scalarmul-in-snark
Related code to the research paper "Fast elliptic curve scalar multiplications in SN(T)ARK circuits" (https://ia.cr/2025/933), published at Latincrypt 2025 (https://link.springer.com/chapter/10.1007/978-3-032-06754-8_4).

*Authors: [Liam Eagen](https://github.com/Liam-Eagen), [Youssef El Housni](https://github.com/yelhousni), [Simon Masson](https://github.com/simonmasson) and [Thomas Piellard](https://github.com/ThomasPiellard).*

## organization
- `sage/` contains SageMath scripts to check the formulas presented in the paper.
- `go/` contains the SNARK circuits written in gnark framework, and lattice/eisenstein reductions in Go.
