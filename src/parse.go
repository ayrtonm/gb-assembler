package main

import (
  "fmt"
  "os"
  "strconv"
  "encoding/binary"
)

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

