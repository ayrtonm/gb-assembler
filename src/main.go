package main

import(
  "os"
  "bufio"
)

const startAddress uint16 = 0x0100
const titleAddress uint16 = 0x0134
const checksumAddress uint16 = 0x014D
const numTitleChars int = 16

var pc uint16 = startAddress
var labels map[string]uint16 = make(map[string]uint16, 0)
var unassignedLabels map[uint16]string = make(map[uint16]string, 0)

func updateAddress(address uint16, file *os.File) {
  file.Seek(int64(address), 0)
  pc = address
}

func main() {
  if (len(os.Args)) < 2 {
    bailout(1)
  }
  infile, err := os.Open(os.Args[1])
  check(err)
  outfile, err := os.Create(os.Args[2])
  check(err)
  defer infile.Close()
  defer outfile.Close()

  rd := bufio.NewReader(infile)
  line, err := getLine(rd)
  labels["start"] = startAddress
  var stop bool = false
  var dataDirective = false
  for {
    switch getSectionType(line, err) {
      case title:
        //read next line, move pc to titleAddress and insert title
        line, err = getLine(rd)
        updateAddress(titleAddress, outfile)
        writeCode(outfile, titleToSlice(line))
        dataDirective = false;
        line, err = getLine(rd)
      case start:
        //move pc to startAddress and continue
        updateAddress(startAddress, outfile)
        dataDirective = false;
        line, err = getLine(rd)
      case address:
        //move pc to address and continue
        updateAddress(getUint16(getLabel(line)), outfile)
        dataDirective = false;
        line, err = getLine(rd)
      case label:
        //make a label at the current pc and continue
        labels[getLabel(line)] = pc
        dataDirective = false;
        line, err = getLine(rd)
      case data:
        dataDirective = true;
        line, err = getLine(rd)
      case comment:
        //ignore line and continue
        line, err = getLine(rd)
      case blank:
        //ignore line and continue
        line, err = getLine(rd)
      case code:
        if dataDirective {
          rawData := dataToSlice(line)
          if len(rawData) != 0 {
            writeCode(outfile, rawData)
            updateAddress(pc + uint16(len(rawData)), outfile)
            line, err = getLine(rd)
          } else {
            dataDirective = false
          }
        } else {
          //insert instruction at current pc and continue
          byteCode := readCode(line)
          writeCode(outfile, byteCode)
          updateAddress(pc + uint16(len(byteCode)), outfile)
          line, err = getLine(rd)
        }
      case eof:
        stop = true
    }
    if stop {
      break
    }
  }
  //fill in jump and call instructions that used labels before the labels were defined
  //addr is the location of the jump/call instruction
  for addr, labelName := range unassignedLabels {
    //assignedAddr is the value we want to write in
    assignedAddr, found := labels[labelName]
    if !found {
      bailout(2)
    }
    outfile.Seek(int64(addr + 1), 0)
    writeCode(outfile, uint16ToSlice(assignedAddr))
  }
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
