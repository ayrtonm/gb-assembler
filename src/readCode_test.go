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
"ld $c $b",
"ld $c $c",
"ld $c $d",
"ld $c $e",
"ld $c $h",
"ld $c $l",
"ld $c [hl]",
"ld $c $a",

"ld $d $b",
"ld $d $c",
"ld $d $d",
"ld $d $e",
"ld $d $h",
"ld $d $l",
"ld $d [hl]",
"ld $d $a",
"ld $e $b",
"ld $e $c",
"ld $e $d",
"ld $e $e",
"ld $e $h",
"ld $e $l",
"ld $e [hl]",
"ld $e $a",

"ld $h $b",
"ld $h $c",
"ld $h $d",
"ld $h $e",
"ld $h $h",
"ld $h $l",
"ld $h [hl]",
"ld $h $a",
"ld $l $b",
"ld $l $c",
"ld $l $d",
"ld $l $e",
"ld $l $h",
"ld $l $l",
"ld $l [hl]",
"ld $l $a",

"ld [hl] $b",
"ld [hl] $c",
"ld [hl] $d",
"ld [hl] $e",
"ld [hl] $h",
"ld [hl] $l",
"halt",
"ld [hl] $a",
"ld $a $b",
"ld $a $c",
"ld $a $d",
"ld $a $e",
"ld $a $h",
"ld $a $l",
"ld $a [hl]",
"ld $a $a"}
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
  "push",
  "pop",
  "inc",
  "dec"}

func TestReadCode(t *testing.T) {
  for _,testName := range allTests {
    for i,v := range testMap[testName] {
      fmt.Println("testing", v)
      bytes := readCode(v)
      if testName == "push" {
        if bytes[0] & 0x0F != 0x05 {
          t.Errorf(v+" wrote %x want x5",bytes)
        }
    } else if testName == "pop" {
        if bytes[0] & 0x0F != 0x01 {
          t.Errorf(v+" wrote %x want x1",bytes)
        }
    } else if testName == "inc" {
        if bytes[0] & 0x0F != 0x03 {
          t.Errorf(v+" wrote %x want x3",bytes)
        }
    } else if testName == "dec" {
        if bytes[0] & 0x0F != 0x0b {
          t.Errorf(v+" wrote %x want xb",bytes)
        }
    } else if testName == "jp" {
      fmt.Printf("%x\n",bytes)
    } else if testName == "ld" {
      if bytes[0] != 0x40 + byte(i) {
        t.Errorf(v+" wrote %x want %x", bytes, 0x40 + i)
      }
    }
    }
  }
}
