package main

import (
  "testing"
  "fmt"
)

var isNumTrue = []string{
  "0x0150",
  "0xffe0",
  "0xff0",
  "00150",
}
var isNumFalse = []string{
  "ff",
  "1b0",
}
var testSectionTypes = []string{
  ".data",
  ".title",
  "start:",
  "0x0150:",
  "150:",
  "function:",
  "ld some shit",
  "",
}

func TestIsNum(t *testing.T) {
  for _,test := range isNumTrue {
    if !isNum(test) {
      t.Errorf("failed test "+test)
    }
    fmt.Println(test, getNum(test))
  }
  for _,test := range isNumFalse {
    if isNum(test) {
      t.Errorf("failed test "+test)
    }
    fmt.Println(test, getNum(test))
  }
  for _,test := range testSectionTypes {
    fmt.Println(test, getSectionType(test,nil))
  }
  for i,j := range titleToSlice("very long demo game") {
    fmt.Println(i,string(j))
  }
}
