package main

import (
  "strings"
)

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

