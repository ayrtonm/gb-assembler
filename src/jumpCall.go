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
  } else {
    bailout(8)
  }
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
  output = append(output, lowByte(newAddress), hiByte(newAddress))
  return output
}

