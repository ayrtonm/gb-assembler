package main

import (
  "testing"
  "fmt"
)

var testLd = []string{
  "ld $b $b",
  "ld $b $c",
  "ld $b $d",
  "ld $b $e",
  "ld $b $h",
  "ld $b $l",
  "ld $b [hl]",
  "ld $b $a",
}
var testJp = []string{
  "ld [hl] 0x3c",
  "jp 0x0100",
  "rst 0x08",
  "call 0x0100"}
var testPush = []string {
  "push $bc",
  "push $de",
  "push $hl",
  "push $af"}
var testPop = []string {
  "pop $bc",
  "pop $de",
  "pop $hl",
  "pop $af"}
var testInc = []string {
  "inc $bc",
  "inc $de",
  "inc $hl",
  "inc $sp"}
var testDec = []string {
  "dec $bc",
  "dec $de",
  "dec $hl",
  "dec $sp"}
var testMap = map[string][]string{
  "ld":testLd,
  "jp":testJp,
  "push":testPush,
  "pop":testPop,
  "inc":testInc,
  "dec":testDec}
var allTests = []string{
  "ld",
  "jp",
  "push",
  "pop",
  "inc",
  "dec"}

func TestReadCode(t *testing.T) {
  for _,testName := range allTests {
    for _,i := range testMap[testName] {
      fmt.Println("testing", i)
      bytes := readCode(i)
      if testName == "push" {
        if bytes[0] & 0x0F != 0x05 {
          t.Errorf(i+" wrote %x want x5",bytes)
        }
    } else if testName == "pop" {
        if bytes[0] & 0x0F != 0x01 {
          t.Errorf(i+" wrote %x want x1",bytes)
        }
    } else if testName == "inc" {
        if bytes[0] & 0x0F != 0x03 {
          t.Errorf(i+" wrote %x want x3",bytes)
        }
    } else if testName == "dec" {
        if bytes[0] & 0x0F != 0x0b {
          t.Errorf(i+" wrote %x want xb",bytes)
        }
    } else if testName == "jp" || testName == "ld" {
      fmt.Printf("%x\n",bytes)
    }
    }
  }
}
