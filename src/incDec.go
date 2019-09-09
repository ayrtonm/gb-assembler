package main

import(
  "os"
)

//increment and decrement
func incDec(dest string, instruction string) (output byte) {
  if isReg(dest) {
    reg := stripReg(dest)
    if regLength(dest) == 2 {
      if instruction == "inc" {
        //two byte register increment
        output = 0x03 + (regOffsets1[reg] * 0x10)
      } else if instruction == "dec" {
        //two byte register decrement
        output = 0x0b + (regOffsets1[reg] * 0x10)
      } else {
        os.Exit(7)
      }
    } else if regLength(dest) == 1 {
      if instruction == "inc" {
        //one byte register increment
        output = 0x04 + (regOffsets2[reg] * 0x08)
      } else if instruction == "dec" {
        //one byte register decrement
        output = 0x05 + (regOffsets2[reg] * 0x08)
      } else {
        os.Exit(7)
      }
    } else {
      //reg is not a valid register
      os.Exit(4)
    }
  } else if isPtr(dest) {
    reg := stripPtr(dest)
    if reg != "hl" {
      os.Exit(4)
    }
    if instruction == "inc" {
      //increment address in hl
      output = 0x34
    } else if instruction == "dec" {
      //decrement address in hl
      output = 0x35
    } else {
      os.Exit(7)
    }
  } else {
    os.Exit(3)
  }
  return output
}
