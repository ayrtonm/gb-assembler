package main

import (
  "os"
  "strings"
)

//may need to clean up these error codes
/*error code explanations
  1 - main program missing arguments
  2 - called parseHex on argument without hex prefix
  3 - argument to opcode is not a register or pointer
  4 - argument to opcode is not a valid register (more than two characters)
  5 - instruction missing arguments
  6 - argument to opcode is not a valid reset vector
  7 - opcode function called with invalid instruction
  8 - argument to opcode is an invalid two-character register, currently unused
  9 - argument to opcode is an invalid one-character register, currently unused
  10 - error in arguments to load instruction, for temporary use
*/

type section int

const (
  title section = iota
  start
  address
  label
  comment
  srcCode
  blank
)

const numTitleChars int = 16
const titleAddress uint16 = 0x134
const startAddress uint16 = 0x0100
const logoAddress uint16 = 0x0104
const checksumAddress uint16 = 0x014D
var nintendoLogo []byte = []byte{0xCE, 0xED, 0x66, 0x66, 0xCC, 0x0D, 0x00, 0x0B, 0x03, 0x73, 0x00, 0x83, 0x00, 0x0C, 0x00, 0x0D, 0x00, 0x08, 0x11, 0x1F, 0x88, 0x89, 0x00, 0x0E, 0xDC, 0xCC, 0x6E, 0xE6, 0xDD, 0xDD, 0xD9, 0x99, 0xBB, 0xBB, 0x67, 0x63, 0x6E, 0x0E, 0xEC, 0xCC, 0xDD, 0xDC, 0x99, 0x9F, 0xBB, 0xB9, 0x33, 0x3E}

//offset pattern used in inc/dec for 16 bit registers
var regOffsets1 = map[string]byte{"bc":0, "de":1, "hl":2, "sp":3}
//offset pattern used in inc/dec for 8 bit registers
//also used for ld between two 8 bit registers
var regOffsets2 = map[string]byte{"b":0, "c":1, "d":2, "e":3, "h":4, "l":5, "a":7}
//offset pattern used in push/pop
var regOffsets3 = map[string]byte{"bc":0, "de":1, "hl":2, "af":3}


func lowByte(x uint16) uint8 {
  return uint8(x & 0x00ff)
}

func hiByte(x uint16) uint8 {
  return uint8(x >> 8)
}

func check(e error) {
  if e != nil {
    panic(e)
  }
}
func updatePc(newPc uint16, file *os.File) {
  file.Seek(int64(newPc), 0)
  pc = newPc
}

func regLength(dest string) int {
  return len(dest)-len(regPrefix)
}

func stripReg(dest string) string {
  return dest[len(regPrefix):]
}

func stripPtr(dest string) string {
  return dest[len(ptrPrefix):len(dest)-len(ptrSuffix)]
}

func stripLabel(dest string) string {
  return strings.TrimSuffix(dest, labelSuffix)
}

func getSectionType(line string) section {
  switch line {
    case "":
      return blank
    case "title"+labelSuffix:
      return title
    case "start"+labelSuffix:
      return start
    default:
      if isComment(line) {
        return comment
      } else if isAddress(line) {
        return address
      } else if isLabel(line) {
        return label
      } else {
        return srcCode
      }
  }
}

func readTitle(line string) []byte {
  var endIndex int = numTitleChars
  if len(line) < numTitleChars {
    endIndex = len(line)
  }
  titleBytes := make([]byte, endIndex)
  if isHex(line) {
    titleBytes = parseHex(line, endIndex)
  } else {
    for i := 0; i < endIndex; i++ {
      titleBytes[i] = line[i]
    }
  }
  return titleBytes
}
