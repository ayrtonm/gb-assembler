package main

import (
  "io"
  "fmt"
  "strings"
  "strconv"
  "os"
  "bufio"
)

const hexPrefix string = "0x"
const directivePrefix string = "."
const regPrefix string = "$"
const ptrPrefix string = "["
const ptrSuffix string = "]"
const commentPrefix string = "//"
const labelSuffix string = ":"

type section int

const (
  eof section = iota
  blank
  title
  start
  address
  label
  comment
  code
  data
)
//offset pattern used in inc/dec for 16 bit registers
var regOffsets1 = map[string]byte{"bc":0, "de":1, "hl":2, "sp":3}
//offset pattern used in inc/dec for 8 bit registers
//also used for ld between two 8 bit registers
var regOffsets2 = map[string]byte{"b":0, "c":1, "d":2, "e":3, "h":4, "l":5, "a":7}
//offset pattern used in push/pop
var regOffsets3 = map[string]byte{"bc":0, "de":1, "hl":2, "af":3}
//can the string represent a number in hexadecimal or otherwise
func isNum(line string) bool {
  if isHex(line) {
    _, err := strconv.ParseInt(getHex(line), 16, 64)
    return err == nil
  } else {
    _, err := strconv.ParseInt(line, 10, 64)
    return err == nil
  }
}
func getNum(line string) int64 {
  if isHex(line) {
    val, _ := strconv.ParseInt(getHex(line), 16, 64)
    return val
  } else {
    val, _ := strconv.ParseInt(line, 10, 64)
    return val
  }
}
func getUint16(line string) uint16 {
  return uint16(getNum(line))
}
func getUint8(line string) uint8 {
  return uint8(getNum(line))
}
//can the string represente a 16-bit number (note that 0x0000ffff would return true)
func isUint16(line string) bool {
  if isNum(line) {
    return getNum(line) < 0x10000
  } else {
    return false
  }
}
func isEof(err error) bool {
  return err == io.EOF
}
func isBlank(line string) bool {
  return line == ""
}
func isTitle(line string) bool {
  return isDirective(line) && getDirective(line) == "title"
}
func isStartAddress(line string) bool {
  return isLabel(line) && getLabel(line) == "start"
}
func isDataDirective(line string) bool {
  return isDirective(line) && getDirective(line) == "data"
}
func isHex(line string) bool {
  return strings.HasPrefix(line, hexPrefix)
}
func isDirective(line string) bool {
  return strings.HasPrefix(line, directivePrefix)
}
func isReg(line string) bool {
  return strings.HasPrefix(line, regPrefix)
}
func isPtr(line string) bool {
  return strings.HasPrefix(line, ptrPrefix) && strings.HasSuffix(line, ptrSuffix)
}
func isComment(line string) bool {
  return strings.HasPrefix(line, commentPrefix)
}
func isLabel(line string) bool {
  return strings.HasSuffix(line, labelSuffix)
}
func isAddress(line string) bool {
  return isLabel(line) && isUint16(getLabel(line))
}
func isValidRst(rstVector uint8) bool {
  return (rstVector & 0xc7) == 0
}
func getHex(line string) string {
  return strings.TrimPrefix(line, hexPrefix)
}
func getDirective(line string) string {
  return strings.TrimPrefix(line, directivePrefix)
}
func getReg(line string) string {
  return strings.TrimPrefix(line, regPrefix)
}
func getPtr(line string) string {
  return strings.TrimPrefix(strings.TrimSuffix(line, ptrSuffix), ptrPrefix)
}
func getLabel(line string) string {
  return strings.TrimSuffix(line, labelSuffix)
}

func getSectionType(line string, e error) section {
  if isEof(e) {
    return eof
  } else if isBlank(line) {
    return blank
  } else if isTitle(line) {
    return title
  } else if isStartAddress(line) {
    return start
  } else if isAddress(line) {
    return address
  } else if isLabel(line) {
    return label
  } else if isComment(line) {
    return comment
  } else if isDataDirective(line) {
    return data
  } else {
    return code
  }
}

func getLine(rd *bufio.Reader) (line string, e error) {
  line, err := rd.ReadString('\n')
  line = strings.TrimSuffix(line, "\n")
  return strings.ToLower(line), err
}

func regLength(line string) int {
  return len(getReg(line))
}

func writeCode(file *os.File, data []byte) {
  nn, err := file.Write(data)
  check(err)
  if nn < len(data) {
    bailout(3)
  }
}

func lowByte(data uint16) uint8 {
  return uint8(data & 0x00ff)
}

func hiByte(data uint16) uint8 {
  return uint8(data >> 8)
}

func titleToSlice(line string) []byte {
  lastIndex := numTitleChars
  if len(line) < numTitleChars {
    lastIndex = len(line)
  }
  output := make([]byte, lastIndex)
  for i := 0; i < lastIndex; i++ {
    output[i] = line[i]
  }
  return output
}

func uint16ToSlice(data uint16) []byte {
  output := make([]byte, 2)
  output[0] = lowByte(data)
  output[1] = hiByte(data)
  return output
}

func stringInList(s string, list []string) bool {
  for _,test := range list {
    if s == test {
      return true
    }
  }
  return false
}

func bailout(code int) {
  switch code {
    case 1:
      fmt.Println("not enough arguments")
    case 2:
      fmt.Println("unassigned labels")
    case 3:
      fmt.Println("Write() wrote less than expected")
    case 4:
      fmt.Println("called arithmetic(dest, instruction) with an invalid dest")
    case 5:
      fmt.Println("called arithmetic(dest, instruction) with an invalid instruction")
    case 6:
      fmt.Println("called incDec(dest, instruction) with invalid instruction")
    case 7:
      fmt.Println("called incDec(dest, instruction) with invalid dest")
    case 8:
      fmt.Println("called jumpCall(dest, instruction) with invalid instruction")
    case 9:
      fmt.Println("load(dest, data) failed in case 1")
    case 10:
      fmt.Println("load(dest, data) failed in case 2")
    case 11:
      fmt.Println("load(dest, data) failed in case 3")
    case 12:
      fmt.Println("load(dest, data) failed in case 4 subcase 1")
    case 13:
      fmt.Println("load(dest, data) failed in case 4 subcase 2")
    case 14:
      fmt.Println("load(dest, data) failed in case 4 default subcase")
    case 15:
      fmt.Println("load(dest, data) failed in case 5")
    case 16:
      fmt.Println("load(dest, data) failed in default case")
    case 17:
      fmt.Println("called pushPop(dest, instruction) with invalid dest")
    case 18:
      fmt.Println("called pushPop(dest, instruction) with invalid instruction")
  }
  fmt.Println("bailing out")
  os.Exit(code)
}

func check(e error) {
  if e != nil {
    panic(e)
  }
}

