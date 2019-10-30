package main

//arithmetic
func arithmetic(dest string, instruction string) (output []byte) {
  if isReg(dest) {
    reg := getReg(dest)
    offset, found := regOffsets2[reg]
    if found {
      var base byte = 0x80
      offsetInstr, foundInstr := opcodeOffsets1[instruction]
      if foundInstr {
        output = append(output, base + (offsetInstr * 8) + offset)
      } else {
        bailout(5)
      }
    } else {
      //reg is not a valid register
      bailout(4)
    }
  } else if isPtr(dest) {
    reg := getPtr(dest)
    if reg != "hl" {
      //reg is not a valid register
      bailout(4)
    }
    var base byte = 0x86
    offset, found := opcodeOffsets1[instruction]
    if found {
      output = append(output, base + (offset * 8))
    } else {
      bailout(5)
    }
  } else if isNum(dest) {
    var base byte = 0xc6
    offset, found := opcodeOffsets1[instruction]
    if found {
      output = append(output, base + (offset * 8))
      output = append(output, getUint8(dest))
    } else {
      bailout(5)
    }
  } else {
    //argument to add is not a register or pointer
    bailout(4)
  }
  return output
}

func addWords(dest string, data string) (output []byte) {
  output = make([]byte,1)
  if isReg(dest) {
    destReg := getReg(dest)
    if destReg == "hl" {
      if isReg(data) {
        dataReg := getReg(data)
        offset, found := regOffsets1[dataReg]
        if found {
          var base byte = 0x09
          output[0] = base + (16 * offset)
        } else {
          bailout(32)
        }
      } else {
        bailout(32)
      }
    } else {
      bailout(31)
    }
  } else {
    bailout(31)
  }
  return output
}
