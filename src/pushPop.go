package main

import(
  "os"
)

//push and pop
func pushPop(dest string, instruction string) (output byte) {
  if isReg(dest) {
    reg := stripReg(dest)
    var base byte
    if instruction == "push" {
      base = 0xc5
    } else if instruction == "pop" {
      base = 0xc1
    } else {
      os.Exit(7)
    }
    output = base + (regOffsets3[reg] * 0x10)
  } else {
    //argument to increment is not a register or pointer
    os.Exit(3)
  }
  return output
}

