package main

//jumps and calls
func jumpCall(dest string, instruction string) (output []byte) {
  output = make([]byte,0)
  if instruction == "jp" {
    output = append(output, 0xc3)
  } else if instruction == "jpz" {
    output = append(output, 0xca)
  } else if instruction == "jpnz" {
    output = append(output, 0xc2)
  } else if instruction == "jpc" {
    output = append(output, 0xda)
  } else if instruction == "jpnc" {
    output = append(output, 0xd2)
  } else if instruction == "call" {
    output = append(output, 0xcd)
  } else if instruction == "callz" {
    output = append(output, 0xcc)
  } else if instruction == "callnz" {
    output = append(output, 0xc4)
  } else if instruction == "callc" {
    output = append(output, 0xdc)
  } else if instruction == "callnc" {
    output = append(output, 0xd4)
  } else if instruction == "jr" {
    output = append(output, 0x18)
  } else if instruction == "jrnz" {
    output = append(output, 0x20)
  } else if instruction == "jrz" {
    output = append(output, 0x28)
  } else if instruction == "jrnc" {
    output = append(output, 0x30)
  } else if instruction == "jrc" {
    output = append(output, 0x38)
  } else {
    bailout(8)
  }
  if isRelativeJump(instruction) {
    if isNum(dest) {
      output = append(output, getUint8(dest))
    }
  } else {
    var newAddress uint16
    var found bool
    if isNum(dest) {
      newAddress = getUint16(dest)
    } else {
      newAddress, found = labels[dest]
      if !found {
        unassignedLabels[pc] = dest
      }
    }
    /*
      make sure to always write the amount of data the instruction expects
      even if that means writing trash to an unassigned label
      since len(output) is used to increment the program counter
    */
    output = append(output, lowByte(newAddress), hiByte(newAddress))
  }
  return output
}

