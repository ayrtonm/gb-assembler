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

    case isReg(dest) && isPtr(data) && !isGenericItr(getPtr(data)):
      destReg := getReg(dest)
      dataPtr := getPtr(data)
      destOffset,foundDest := regOffsets2[destReg]
      if dataPtr == "hl" && foundDest {
        output = append(output, 0x46 + (destOffset * 0x08))
      } else if destReg == "a" {
        var newAddress uint16
        var found bool
        if isNum(dataPtr) {
          newAddress = getUint16(dataPtr)
        } else {
          newAddress, found = labels[dataPtr]
          if !found {
            unassignedLabels[pc] = dataPtr
          }
        }
        /*
          make sure to always write the amount of data the instruction expects
          even if that means writing trash to an unassigned label
          since len(output) is used to increment the program counter
        */
        output = append(output, 0xfa, lowByte(newAddress), hiByte(newAddress))
      } else {
        bailout(10)
      }

    case isPtr(dest) && isReg(data) && !isGenericItr(getPtr(dest)):
      destPtr := getPtr(dest)
      dataReg := getReg(data)
      dataOffset,foundData := regOffsets2[dataReg]
      if destPtr  == "hl" && foundData {
        output = append(output, 0x70 + dataOffset)
      } else if dataReg == "a" {
        var newAddress uint16
        var found bool
        if isNum(destPtr) {
          newAddress = getUint16(destPtr)
        } else {
          newAddress, found = labels[destPtr]
          if !found {
            unassignedLabels[pc] = destPtr
          }
        }
        /*
          make sure to always write the amount of data the instruction expects
          even if that means writing trash to an unassigned label
          since len(output) is used to increment the program counter
        */
        output = append(output, 0xea, lowByte(newAddress), hiByte(newAddress))
      } else {
        bailout(11)
      }

    case isPtr(dest) && isNum(data) && !isGenericItr(getPtr(dest)):
      destPtr := getPtr(dest)
      if destPtr == "hl" {
        output = append(output, 0x36, getUint8(data))
      } else {
        bailout(15)
      }

    case isItrPtr(dest) && isReg(data):
      destItrPtr := getItrPtr(dest)
      dataReg := getReg(data)
      if destItrPtr == "hl" && dataReg == "a" {
        output = append(output, 0x22)
      } else {
        bailout(20)
      }

    case isRevItrPtr(dest) && isReg(data):
      destRevItrPtr := getRevItrPtr(dest)
      dataReg := getReg(data)
      if destRevItrPtr == "hl" && dataReg == "a" {
        output = append(output, 0x32)
      } else {
        bailout(20)
      }

    case isItrPtr(data) && isReg(dest):
      dataItrPtr := getItrPtr(data)
      destReg := getReg(dest)
      if dataItrPtr == "hl" && destReg == "a" {
        output = append(output, 0x2a)
      } else {
        bailout(20)
      }

    case isRevItrPtr(data) && isReg(dest):
      dataRevItrPtr := getRevItrPtr(data)
      destReg := getReg(dest)
      if dataRevItrPtr == "hl" && destReg == "a" {
        output = append(output, 0x3a)
      } else {
        bailout(20)
      }

    case isReg(dest) && regLength(dest) == 1 && isNum(data):
      destReg := getReg(dest)
      destOffset, foundDest := regOffsets2[destReg]
      if foundDest {
        output = append(output, 0x06 + (destOffset * 0x08), getUint8(data))
      } else {
        bailout(12)
      }

    case isReg(dest) && regLength(dest) == 2:
      destReg := getReg(dest)
      destOffset, foundDest := regOffsets1[destReg]
      if foundDest {
        var dataAddress uint16
        var found bool
        if isNum(data) {
          dataAddress = getUint16(data)
        } else {
          dataAddress, found = labels[data]
          if !found {
            unassignedLabels[pc] = data
          }
        }
        output = append(output, 0x01 + (destOffset * 0x10), lowByte(dataAddress), hiByte(dataAddress))
      } else {
        bailout(13)
      }

    default:
      bailout(16)
  }
  return output
}


