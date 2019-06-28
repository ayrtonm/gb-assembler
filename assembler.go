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

//functions used to write opcodes
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

//push and pop
func pushPop(dest string, instruction string) (output byte) {
  if isReg(dest) {
    reg := stripReg(dest)
    var base byte
    if instruction == "push" {
      base = 0xc5
    } else if instruction == "pop" {
      base = 0xc1
    } else {
      os.Exit(7)
    }
    output = base + (regOffsets3[reg] * 0x10)
  } else {
    //argument to increment is not a register or pointer
    os.Exit(3)
  }
  return output
}

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
    os.Exit(7)
  }
  newAddress := parseWord(dest)
  output = append(output, lowByte(newAddress), hiByte(newAddress))
  return output
}

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

func readCode(line string) (byteCode []byte) {
  output := make([]byte,0)
  cmd := strings.Fields(line)
  instruction := cmd[0]
  //handle instructions with no arguments here
  switch instruction {
    case "nop":
      output = append(output, 0x00)
    case "stop":
      output = append(output, 0x10)
    case "halt":
      output = append(output, 0x76)
    case "ei":
      output = append(output, 0xfb)
    case "di":
      output = append(output, 0xf3)
    case "ret":
      output = append(output, 0xc9)
    case "retz":
      output = append(output, 0xc8)
    //if instruction not found, read one argument then switch case with one-argument instructions
    default:
      if len(cmd) < 1 {
        //instructions missing arguments
        os.Exit(5)
      }
      dest := cmd[1]
      switch instruction {
        case "jp":
          output = append(output, jumpCall(dest, "jp")...)
        case "jpz":
          output = append(output, jumpCall(dest, "jpz")...)
        case "jpnz":
          output = append(output, jumpCall(dest, "jpnz")...)
        case "jpc":
          output = append(output, jumpCall(dest, "jpc")...)
        case "jpnc":
          output = append(output, jumpCall(dest, "jpnc")...)
        case "call":
          output = append(output, jumpCall(dest, "call")...)
        case "callz":
          output = append(output, jumpCall(dest, "callz")...)
        case "callnz":
          output = append(output, jumpCall(dest, "callnz")...)
        case "callc":
          output = append(output, jumpCall(dest, "callc")...)
        case "callnc":
          output = append(output, jumpCall(dest, "callnc")...)
        case "rst":
          newAddress := parseByte(dest)
          if !isValidRst(newAddress) {
            //reset vector is not valid
            os.Exit(6)
          }
          output = append(output, 0xc7 + newAddress)
        case "push":
          output = append(output, pushPop(dest, "push"))
        case "pop":
          output = append(output, pushPop(dest, "pop"))
        case "inc":
          output = append(output, incDec(dest, "inc"))
        case "dec":
          output = append(output, incDec(dest, "dec"))
        case "add":
          output = append(output, arithmetic(dest, "add"))
        case "adc":
          output = append(output, arithmetic(dest, "adc"))
        case "sub":
          output = append(output, arithmetic(dest, "sub"))
        case "sbc":
          output = append(output, arithmetic(dest, "sbc"))
        case "and":
          output = append(output, arithmetic(dest, "and"))
        case "xor":
          output = append(output, arithmetic(dest, "xor"))
        case "or":
          output = append(output, arithmetic(dest, "or"))
        case "cp":
          output = append(output, arithmetic(dest, "cp"))
        //instruction not found, read second argument the switch case with two-argument instructions
        default:
          if len(cmd) < 2 {
            os.Exit(5)
          }
          data := cmd[2]
          switch instruction {
            case "ld":
            output = append(output, load(dest, data)...)
            default:
              output = append(output , 0xff)
          }
      }
  }
  return output
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
