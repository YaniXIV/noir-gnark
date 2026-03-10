# noir-gnark

Go deserves to be a first class citizen in the ZK ecosystem. noir-gnark is another step toward that.

A pure Go backend for proving and verifying Noir circuits using [gnark](https://github.com/Consensys/gnark). No FFI, no Rust, no external toolchain. Just Go.

This is a companion to [noir-go](https://github.com/YaniXIV/noir-go). Once you have an ACIR artifact from compilation, noir-gnark takes it from there.

## The pipeline
```
Noir circuit → ACIR (noir-go) → R1CS → Groth16 proof (noir-gnark)
```
