package main

func testBit(dest string, data string) (output []byte) {
  output = make([]byte, 2)
  output[0] = 0xcb
  if isReg(dest) {
    reg := getReg(dest)
    offset, found := regOffsets2[reg]
    if found {
      if isNum(data) {
        bitNumber := byte(getNum(data))
        if bitNumber < 8 {
          var base byte = 0x40
          output[1] = base + offset + (bitNumber  * 8)
        } else {
          bailout(26)
        }
      } else {
        bailout(26)
      }
    } else {
      bailout(25)
    }
  } else if isPtr(dest) {
    reg := getPtr(dest)
    if reg != "hl" {
      bailout(25)
    }
    if isNum(data) {
      bitNumber := byte(getNum(data))
      if bitNumber < 8 {
        var base byte = 0x46
        output[1] = base + (bitNumber * 8)
      } else {
        bailout(26)
      }
    } else {
      bailout(26)
    }
  } else {
    bailout(25)
  }
  return output
}
func setBit(dest string, data string) (output []byte) {
  output = make([]byte, 2)
  output[0] = 0xcb
  if isReg(dest) {
    reg := getReg(dest)
    offset, found := regOffsets2[reg]
    if found {
      if isNum(data) {
        bitNumber := byte(getNum(data))
        if bitNumber < 8 {
          var base byte = 0xc0
          output[1] = base + offset + (bitNumber  * 8)
        } else {
          bailout(28)
        }
      } else {
        bailout(28)
      }
    } else {
      bailout(27)
    }
  } else if isPtr(dest) {
    reg := getPtr(dest)
    if reg != "hl" {
      bailout(27)
    }
    if isNum(data) {
      bitNumber := byte(getNum(data))
      if bitNumber < 8 {
        var base byte = 0xc6
        output[1] = base + (bitNumber * 8)
      } else {
        bailout(28)
      }
    } else {
      bailout(28)
    }
  } else {
    bailout(27)
  }
  return output
}
func clearBit(dest string, data string) (output []byte) {
  output = make([]byte, 2)
  output[0] = 0xcb
  if isReg(dest) {
    reg := getReg(dest)
    offset, found := regOffsets2[reg]
    if found {
      if isNum(data) {
        bitNumber := byte(getNum(data))
        if bitNumber < 8 {
          var base byte = 0x80
          output[1] = base + offset + (bitNumber  * 8)
        } else {
          bailout(30)
        }
      } else {
        bailout(30)
      }
    } else {
      bailout(29)
    }
  } else if isPtr(dest) {
    reg := getPtr(dest)
    if reg != "hl" {
      bailout(29)
    }
    if isNum(data) {
      bitNumber := byte(getNum(data))
      if bitNumber < 8 {
        var base byte = 0x86
        output[1] = base + (bitNumber * 8)
      } else {
        bailout(30)
      }
    } else {
      bailout(30)
    }
  } else {
    bailout(29)
  }
  return output
}
