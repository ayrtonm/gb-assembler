package main

import (
  "fmt"
  "os"
  "io"
  "strings"
  "bufio"
)

var labels map[string]uint16 = make(map[string]uint16, 0)
var pc uint16 = startAddress

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
