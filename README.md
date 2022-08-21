# package arm

```go
import "github.com/wdamron/arm"
```

Package `arm` implements an ARMv8 (AArch64) instruction assembler in Go, for runtime or ahead-of-time generation of executable code. SVE/SME instructions are not yet supported.

This library is mostly adapted from the [CensoredUsername/dynasm-rs](https://github.com/CensoredUsername/dynasm-rs) (Rust) project, and is not heavily tested.

## Brief Overview

The `Assembler` type encodes executable instructions to a code buffer.

Some instructions support label offset arguments, which may be resolved by the `Assembler`
and encoded after all label addresses are assigned.

The following are argument types:
- `Reg`: integer, SP, SIMD scalar, or SIMD vector register (with optional element index)
- `RegList`: list of sequential registers
- `Ref`: memory reference with register base, optionally followed by X register or immediate for post-indexing
- `RefOffset`: memory reference with register base and immediate offset
- `RefPreIndexed`: pre-indexed memory reference with register base and immediate offset
- `RefIndexed`: memory index with register base, register index, and optional index modifier
- `Imm`: 32-bit immediate integer
- `Float`: 32-bit immediate float
- `Wide`: 64-bit immediate integer
- `Mod`: modifier with optional immediate shift/rotate
- `Label`: label reference with optional offset from label address
- `Symbol`: constant identifier

## Additional References

- [Package Documentation (pkg.go.dev)](https://pkg.go.dev/github.com/wdamron/arm)
- [Full Instruction Listing (INSTRUCTIONS.md)](./INSTRUCTIONS.md)
- [Basic Test Examples (asm_test.go)](./asm_test.go)
- [AArch64 Instruction Set Architecture (developer.arm.com)](https://developer.arm.com/documentation/102374/latest/)
