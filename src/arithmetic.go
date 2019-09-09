package main

import(
  "os"
)

//arithmetic
func arithmetic(dest string, instruction string) (output byte) {
  if isReg(dest) {
    reg := stripReg(dest)
    if regLength(dest) == 1 {
      var base byte
      if instruction == "add" {
        base = 0x80
      } else if instruction == "adc" {
        base = 0x88
      } else if instruction == "sub" {
        base = 0x90
      } else if instruction == "sbc" {
        base = 0x98
      } else if instruction == "and" {
        base = 0xa0
      } else if instruction == "xor" {
        base = 0xa8
      } else if instruction == "or" {
        base = 0xb0
      } else if instruction == "cp" {
        base = 0xb8
      } else {
        os.Exit(7)
      }
      output = base + regOffsets2[reg]
    } else {
      //reg is not a valid register
      os.Exit(4)
    }
  } else if isPtr(dest) {
    reg := stripPtr(dest)
    if reg != "hl" {
      //reg is not a valid register
      os.Exit(4)
    }
    if instruction == "add" {
      output = 0x86
    } else if instruction == "adc" {
      output = 0x8e
    } else if instruction == "sub" {
      output = 0x96
    } else if instruction == "sbc" {
      output = 0x9e
    } else if instruction == "and" {
      output = 0xa6
    } else if instruction == "xor" {
      output = 0xae
    } else if instruction == "or" {
      output = 0xb6
    } else if instruction == "cp" {
      output = 0xbe
    } else {
      os.Exit(7)
    }
  } else {
    //argument to add is not a register or pointer
    os.Exit(3)
  }
  return output
}
