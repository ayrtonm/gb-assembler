package main

func rotateShiftSwap(dest string, instruction string) (output []byte) {
  output = make([]byte, 2)
  output[0] = 0xcb
  if isReg(dest) {
    reg := getReg(dest)
    offset, found := regOffsets2[reg]
    if found {
      offsetInstr, foundInstr := opcodeOffsets2[instruction]
      if foundInstr {
        output[1] = (offsetInstr * 8) + offset
      } else {
        bailout(24)
      }
    } else {
      bailout(23)
    }
  } else if isPtr(dest) {
    reg := getPtr(dest)
    if reg != "hl" {
      bailout(23)
    }
    var base byte = 0x06
    offset, found := opcodeOffsets2[instruction]
    if found {
      output[1] = base + (offset * 8)
    } else {
      bailout(24)
    }
  } else {
    bailout(23)
  }
  return output
}
