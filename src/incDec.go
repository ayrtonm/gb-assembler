package main

//increment and decrement
func incDec(dest string, instruction string) (output byte) {
  if isReg(dest) {
    reg := getReg(dest)
    offset, found := regOffsets1[reg]
    if found {
      if instruction == "inc" {
        //two byte register increment
        output = 0x03 + (offset * 0x10)
      } else if instruction == "dec" {
        //two byte register decrement
        output = 0x0b + (offset * 0x10)
      } else {
        bailout(6)
      }
    } else {
      offset, found := regOffsets2[reg]
      if found {
        if instruction == "inc" {
          //one byte register increment
          output = 0x04 + (offset * 0x08)
        } else if instruction == "dec" {
          //one byte register decrement
          output = 0x05 + (offset * 0x08)
        } else {
          bailout(6)
        }
      } else {
        //reg is not a valid register
        bailout(7)
        }
      }
  } else if isPtr(dest) {
    reg := getPtr(dest)
    if reg != "hl" {
      bailout(7)
    }
    if instruction == "inc" {
      //increment address in hl
      output = 0x34
    } else if instruction == "dec" {
      //decrement address in hl
      output = 0x35
    } else {
      bailout(6)
    }
  } else {
    bailout(7)
  }
  return output
}
