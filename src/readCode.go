package main

import(
  "os"
  "strings"
)

//take a line of assembly and turn it into a sequence of bytes
func readCode(line string) (byteCode []byte) {
  output := make([]byte,0)
  cmd := strings.Fields(line)
  instruction := cmd[0]
  //handle instructions with no arguments here
  switch instruction {
    case "cpl":
      output = append(output, 0x2f)
    case "nop":
      output = append(output, 0x00)
    case "stop":
      output = append(output, 0x10)
    case "halt":
      output = append(output, 0x76)
    case "ei":
      output = append(output, 0xfb)
    case "di":
      output = append(output, 0xf3)
    case "ret":
      output = append(output, 0xc9)
    case "retz":
      output = append(output, 0xc8)
    case "retc":
      output = append(output, 0xd8)
    case "retnz":
      output = append(output, 0xc0)
    case "retnc":
      output = append(output, 0xd0)
    //if instruction not found, read one argument then switch case with one-argument instructions
    default:
      if len(cmd) < 2 {
        //instructions missing arguments
        bailout(21)
      }
      dest := cmd[1]
      opFunc := getOneArgOpFunc(instruction)
      if opFunc != nil {
        output = append(output, opFunc(dest, instruction)...)
      } else {
        if instruction == "rst" {
          newAddress := getUint8(dest)
          if !isValidRst(newAddress) {
            //reset vector is not valid
            os.Exit(6)
          }
          output = append(output, 0xc7 + newAddress)
        } else {
          if len(cmd) < 2 {
            bailout(21)
          }
          data := cmd[2]
          opFunc := getTwoArgOpFunc(instruction)
          if opFunc != nil {
            output = append(output, opFunc(dest, data)...)
          } else {
            bailout(21)
          }
        }
      }
  }
  return output
}
