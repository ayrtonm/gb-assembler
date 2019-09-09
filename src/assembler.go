package main

import (
  "fmt"
  "os"
  "io"
  "strings"
  "bufio"
  "strconv"
  "encoding/binary"
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

const commentPrefix string = "//"
const labelSuffix string = ":"
const hexPrefix string = "0x"
const regPrefix string = "$"
const ptrPrefix string = "["
const ptrSuffix string = "]"
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

var labels map[string]uint16 = make(map[string]uint16, 0)
var pc uint16 = startAddress

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

func isValidRst(rstVector uint8) bool {
  return (rstVector & 0xc7) == 0
}

func isPtr(line string) bool {
  return strings.HasPrefix(line, ptrPrefix) && strings.HasSuffix(line, ptrSuffix)
}

func isReg(line string) bool {
  return strings.HasPrefix(line, regPrefix)
}

func isHex(line string) bool {
  return strings.HasPrefix(line, hexPrefix)
}

func isComment(line string) bool {
  return strings.HasPrefix(line, commentPrefix)
}

func isLabel(line string) bool {
  return strings.HasSuffix(line, labelSuffix)
}

func isAddress(line string) bool {
  return isHex(line) && isLabel(line)
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

func parseHex(input string, numChars int) []byte {
  if !isHex(input) {
    fmt.Println("failed to parse:", input)
    os.Exit(2)
  }
  output := make([]byte, numChars)
  const charsPerByte = 2
  input = input[2:]
  for i := 0; i < numChars; i++ {
    charByte,_ := strconv.ParseInt(input[charsPerByte*i:(charsPerByte*i)+charsPerByte],16,64)
    output[i] = byte(charByte)
  }
  return output
}

func parseWord(input string) uint16 {
  return binary.BigEndian.Uint16(parseHex(input, 2))
}

func parseByte(input string) uint8 {
  return parseHex(input, 1)[0]
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

func getNextline(rd *bufio.Reader) (line string, more bool) {
  line, err := rd.ReadString('\n')
  line = strings.TrimSuffix(line, "\n")
  if err != nil {
    if len(line) == 0 && err == io.EOF {
      return "", false 
    }
  } else if getSectionType(line) == comment {
    return getNextline(rd)
  } else if getSectionType(line) == blank {
    return getNextline(rd)
  }
  return line, true
}

func writeCode(file *os.File, data []byte) {
  _, err := file.Write(data)
  check(err)
  //fmt.Println(bytesWritten)
}

func readSection(outfile *os.File, rd *bufio.Reader) (line string, more bool) {
  line, more = getNextline(rd)
  for ; getSectionType(line) == srcCode; line, more = getNextline(rd) {
    byteCode := readCode(line)
    writeCode(outfile, byteCode)
    updatePc(pc + uint16(len(byteCode)), outfile)
    if !more {
      return "", more
    }
  }
  return line, more
}

func main() {
  if (len(os.Args)) < 2 {
    fmt.Println("not enough arguments")
    os.Exit(1)
  }

  infile, err := os.Open(os.Args[1])
  check(err)
  outfile, err := os.Create(os.Args[2])
  check(err)
  defer infile.Close()
  defer outfile.Close()

  rd := bufio.NewReader(infile)
  line, more := getNextline(rd)
  for {
    switch getSectionType(line) {
      case title:
        line, more = getNextline(rd)
        outfile.Seek(int64(titleAddress), 0)
        writeCode(outfile, readTitle(line))
        line, more = getNextline(rd)
      case start:
        updatePc(startAddress, outfile)
        line, more = readSection(outfile, rd)
      case address:
        updatePc(parseWord(line), outfile)
        line, more = readSection(outfile, rd)
      case label:
        labels[strings.TrimSuffix(line, labelSuffix)] = pc
        line, more = readSection(outfile, rd)
    }
    if !more {
      break
    }
  }
  outfile.Seek(int64(logoAddress), 0)
  _, err = outfile.Write(nintendoLogo)

  //compute checksum and write to header
  outfile.Seek(int64(titleAddress), 0)
  checksum := []byte{0}
  temp := []byte{0}
  for i := titleAddress; i < checksumAddress; i++ {
    outfile.Read(temp)
    checksum[0] -= (1 + temp[0])
  }
  outfile.Seek(int64(checksumAddress), 0)
  _, err = outfile.Write(checksum)
  check(err)

  //fill in rest of the file with zeros
  outfile.Seek(0,2)
  var fill []byte = make([]byte,0x8000)
  _, err = outfile.Write(fill)
  check(err)
  outfile.Truncate(0x8000)
  outfile.Sync()
}
