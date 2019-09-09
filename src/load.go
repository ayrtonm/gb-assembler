package main

import(
  "os"
  "fmt"
)

func load(dest string, data string) (output []byte) {
  switch {
    case isReg(dest) && isReg(data):
      destReg := stripReg(dest)
      dataReg := stripReg(data)
      output = append(output, 0x40 + (regOffsets2[destReg] * 0x08) + regOffsets2[dataReg])

    case isReg(dest) && isPtr(data):
      destReg := stripReg(dest)
      dataPtr := stripPtr(data)
      if dataPtr != "hl" {
        os.Exit(4)
      }
      output = append(output, 0x46 + (regOffsets2[destReg] * 0x08))

    case isPtr(dest) && isReg(data):
      destPtr := stripPtr(dest)
      dataReg := stripReg(data)
      if destPtr != "hl" {
        os.Exit(4)
      }
      output = append(output, 0x70 + regOffsets2[dataReg])

    case isReg(dest) && isHex(data):
      destReg := stripReg(dest)
      switch regLength(dest) {
        case 1:
          output = append(output, 0x06 + (regOffsets2[destReg] * 0x08), parseByte(data))
        case 2:
          dataAddress := parseWord(data)
          output = append(output, 0x01 + (regOffsets1[destReg] * 0x10), lowByte(dataAddress), hiByte(dataAddress))
        default:
          os.Exit(4)
      }

    case isPtr(dest) && isHex(data):
      destPtr := stripPtr(dest)
      if destPtr != "hl" {
        os.Exit(4)
      } else {
        output = append(output, 0x36, parseByte(data))
      }

    default:
      fmt.Println("failed to parse: ld", dest, data)
      os.Exit(10)
  }
  return output
}


