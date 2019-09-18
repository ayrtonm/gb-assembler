package main

//push and pop
func pushPop(dest string, instruction string) (output byte) {
  if isReg(dest) {
    reg := getReg(dest)
    var base byte
    if instruction == "push" {
      base = 0xc5
    } else if instruction == "pop" {
      base = 0xc1
    } else {
      bailout(18)
    }
    regOffset, found := regOffsets3[reg]
    if found {
      output = base + (regOffset * 0x10)
    } else {
      bailout(17)
    }
  } else {
    //argument to increment is not a register or pointer
    bailout(17)
  }
  return output
}

