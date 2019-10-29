package main

import(
  "os"
  "bufio"
  "strings"
  "strconv"
)

const startAddress uint16 = 0x0100
const mainAddress uint16 = 0x0150
const titleAddress uint16 = 0x0134
const ramSizeAddress uint16 = 0x0149
const checksumAddress uint16 = 0x014D
const numTitleChars int = 16
const nintendoLogoAddress uint16 = 0x0104
const topScopeLevel int = 0

var nintendoLogoData []uint8 = []uint8{
  0xce, 0xed, 0x66, 0x66, 0xcc, 0x0d, 0x00, 0x0b, 0x03, 0x73, 0x00, 0x83, 0x00, 0x0c, 0x00, 0x0d,
  0x00, 0x08, 0x11, 0x1f, 0x88, 0x89, 0x00, 0x0e, 0xdc, 0xcc, 0x6e, 0xe6, 0xdd, 0xdd, 0xd9, 0x99,
  0xbb, 0xbb, 0x67, 0x63, 0x6e, 0x0e, 0xec, 0xcc, 0xdd, 0xdc, 0x99, 0x9f, 0xbb, 0xb9, 0x33, 0x3e}

var pc uint16 = mainAddress
var eram_counter uint16 = 0xa000
var wram_counter uint16 = 0xc000
var labelsPtr []map[string]uint16 = make([]map[string]uint16, 0)
var unassignedLabelsPtr []map[uint16]string = make([]map[uint16]string, 0)
var scopeLevel int = topScopeLevel
var indentationLevel int = 0

var outfile *os.File
var opfile *os.File
var opsUsed []uint8
var infileQueue []string

//takes a variable directive and the current RAM counter as input and returns an updated counter
func allocateVariable(cmd []string, counter uint16) uint16 {
  labelsPtr[topScopeLevel][cmd[1]] = counter
  if cmd[2] == "byte" {
    counter += 1
  } else if cmd[2] == "word" {
    counter += 2
  } else {
    bailout(22)
  }
  return counter
}

func updateAddress(address uint16, file *os.File) {
  file.Seek(int64(address), 0)
  pc = address
}

func parseFile(filename string) {
  infile, err := os.Open(filename)
  check(err)
  defer infile.Close()

  rd := bufio.NewReader(infile)
  line, err := getLine(rd, outfile)

  for {
    lineType := getSectionType(line, err)
    switch lineType {
      case title:
        //read next line, move to titleAddress, insert title then jump back to pc
        line, err = getLine(rd, outfile)
        prevLocation := pc
        updateAddress(titleAddress, outfile)
        writeCode(outfile, titleToSlice(line))
        updateAddress(prevLocation, outfile)
        line, err = getLine(rd, outfile)
      case start:
        //move pc to startAddress
        updateAddress(startAddress, outfile)
        line, err = getLine(rd, outfile)
      case main_section:
        //move pc to mainAddress
        updateAddress(mainAddress, outfile)
        line, err = getLine(rd, outfile)
      case address:
        //move pc to address
        updateAddress(getUint16(getLabel(line)), outfile)
        line, err = getLine(rd, outfile)
      case label:
        //make a label at the current pc
        labelsPtr[scopeLevel][getLabel(line)] = pc
        line, err = getLine(rd, outfile)
      case comment:
        //ignore line and continue
        line, err = getLine(rd, outfile)
      case blank:
        //ignore line and continue
        line, err = getLine(rd, outfile)
      case alias:
        //make an arbitrary label
        cmd := strings.Fields(line)
        labelsPtr[scopeLevel][cmd[1]] = getUint16(cmd[2])
        line, err = getLine(rd, outfile)
      case savedVariable:
        //make a label to external RAM
        cmd := strings.Fields(line)
        eram_counter = allocateVariable(cmd, eram_counter)
        line, err = getLine(rd, outfile)
      case variable:
        //make a label to work RAM
        cmd := strings.Fields(line)
        wram_counter = allocateVariable(cmd, wram_counter)
        line, err = getLine(rd, outfile)
      case code:
        //insert raw data
        if isNum(line) {
          //FIXME: dataToSlice() returns nil for non-hex numbers
          //I need to fix the ambiguity of determining the number of bytes to write for decimal numbers
          rawData := dataToSlice(line)
          writeCode(outfile, rawData)
          updateAddress(pc + uint16(len(rawData)), outfile)
          line, err = getLine(rd, outfile)
        } else {
          //insert instruction at current pc
          byteCode := readCode(line)
          writeCode(outfile, byteCode)
          if (len(os.Args) == 4) {
            if !opAlreadyWritten(byteCode[0]) {
              opsUsed = append(opsUsed, byteCode[0])
              nn, err := opfile.WriteString(strconv.FormatInt(int64(byteCode[0]),16)+" - "+line+"\n")
              check(err)
              if nn < len(strconv.FormatInt(int64(byteCode[0]),16)+" - "+line+"\n") {
                bailout(3)
              }
            }
          }
          updateAddress(pc + uint16(len(byteCode)), outfile)
          line, err = getLine(rd, outfile)
        }
      case include:
        //add a file to the input file queue
        cmd := strings.Fields(line)
        infileQueue = append(infileQueue, cmd[1])
        line, err = getLine(rd, outfile)
    }
    //EOF should be handled outside the switch so we can break out of the for loop
    if lineType == eof {
      break
    }
  }
}

func main() {
  if (len(os.Args)) < 2 {
    bailout(1)
  }
  labelsPtr = append(labelsPtr, make(map[string]uint16, 0))
  unassignedLabelsPtr = append(unassignedLabelsPtr, make(map[uint16]string, 0))

  labelsPtr[topScopeLevel]["start"] = startAddress
  labelsPtr[topScopeLevel]["main"] = mainAddress

  var err error
  outfile, err = os.Create(os.Args[2])
  check(err)
  defer outfile.Close()

  if (len(os.Args) == 4) {
    opfile, err = os.Create(os.Args[3])
    check(err)
    defer opfile.Close()
  }


  infileQueue = append(infileQueue, os.Args[1])

  for len(infileQueue) != 0 {
    currentFile := infileQueue[0]
    infileQueue = infileQueue[1:]
    parseFile(currentFile)
  }

  //fill in jump and call instructions that used labels before the labels were defined
  //addr is the location of the jump/call instruction
  //fillInUnassignedLabels(topScopeLevel, outfile)
  for scopeLevel >= topScopeLevel {
    fillInUnassignedLabels(outfile)
    if scopeLevel != topScopeLevel {
      labelsPtr = labelsPtr[:len(labelsPtr)-1]
      unassignedLabelsPtr = unassignedLabelsPtr[:len(unassignedLabelsPtr)-1]
    }
    scopeLevel--
  }
  //add nintendo logo data to header
  outfile.Seek(int64(nintendoLogoAddress),0)
  writeCode(outfile, nintendoLogoData)
  //set RAM size to 8 kb by default
  outfile.Seek(int64(ramSizeAddress),0)
  writeCode(outfile, []byte{0x02})
  //compute checksum and write to header
  outfile.Seek(int64(titleAddress),0)
  checksum := []byte{0}
  temp := []byte{0}
  for i := titleAddress; i < checksumAddress; i++ {
    outfile.Read(temp)
    checksum[0] -= (1 + temp[0])
  }
  outfile.Seek(int64(checksumAddress),0)
  writeCode(outfile, checksum)

  outfile.Seek(0,2)
  var fill []byte = make([]byte,0x8000)
  writeCode(outfile, fill)
  outfile.Truncate(0x8000)
  outfile.Sync()
}
