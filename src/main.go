package main

import(
  "os"
  "bufio"
)

const startAddress uint16 = 0x0100
const titleAddress uint16 = 0x0134
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
  for {
    switch getSectionType(line, err) {
      case title:
        //read next line, move pc to titleAddress and insert title
        line, err = getLine(rd)
        updateAddress(titleAddress, outfile)
        writeCode(outfile, titleToSlice(line))
        line, err = getLine(rd)
      case start:
        //move pc to startAddress and continue
        updateAddress(startAddress, outfile)
        line, err = getLine(rd)
      case address:
        //move pc to address and continue
        updateAddress(getUint16(getLabel(line)), outfile)
        line, err = getLine(rd)
      case label:
        //make a label at the current pc and continue
        labels[getLabel(line)] = pc
        line, err = getLine(rd)
      case data:
        //how do I handle this?
        stop = true
      case comment:
        //ignore line and continue
        line, err = getLine(rd)
      case blank:
        //ignore line and continue
        line, err = getLine(rd)
      case code:
        //insert instruction at current pc and continue
        byteCode := readCode(line)
        writeCode(outfile, byteCode)
        updateAddress(pc + uint16(len(byteCode)), outfile)
        line, err = getLine(rd)
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
  outfile.Seek(0,2)
  var fill []byte = make([]byte,0x8000)
  writeCode(outfile, fill)
  outfile.Truncate(0x8000)
  outfile.Sync()
}
