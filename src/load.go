package main

func load(dest string, data string) (output []byte) {
  switch {
    case isReg(dest) && isReg(data):
      destReg := getReg(dest)
      dataReg := getReg(data)
      destOffset,foundDest := regOffsets2[destReg]
      dataOffset,foundData := regOffsets2[dataReg]
      if foundDest && foundData {
        output = append(output, 0x40 + (destOffset * 0x08) + dataOffset)
      } else {
        bailout(9)
      }

    case isReg(dest) && isPtr(data):
      destReg := getReg(dest)
      dataPtr := getPtr(data)
      destOffset,foundDest := regOffsets2[destReg]
      if dataPtr == "hl" && foundDest {
        output = append(output, 0x46 + (destOffset * 0x08))
      } else {
        bailout(10)
      }

    case isPtr(dest) && isReg(data):
      destPtr := getPtr(dest)
      dataReg := getReg(data)
      dataOffset,foundData := regOffsets2[dataReg]
      if destPtr  == "hl" && foundData {
        output = append(output, 0x70 + dataOffset)
      } else {
        bailout(11)
      }

    case isReg(dest) && isHex(data):
      destReg := getReg(dest)
      switch regLength(dest) {
        case 1:
          destOffset, foundDest := regOffsets2[destReg]
          if foundDest {
            output = append(output, 0x06 + (destOffset * 0x08), getUint8(data))
          } else {
            bailout(12)
          }
        case 2:
          dataAddress := getUint16(data)
          destOffset, foundDest := regOffsets1[destReg]
          if foundDest {
            output = append(output, 0x01 + (destOffset * 0x10), lowByte(dataAddress), hiByte(dataAddress))
          } else {
            bailout(13)
          }
        default:
          bailout(14)
      }

    case isPtr(dest) && isHex(data):
      destPtr := getPtr(dest)
      if destPtr != "hl" {
        bailout(15)
      } else {
        output = append(output, 0x36, getUint8(data))
      }

    default:
      bailout(16)
  }
  return output
}


