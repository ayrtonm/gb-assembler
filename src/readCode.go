package main

import(
  "os"
  "strings"
)

func readCode(line string) (byteCode []byte) {
  output := make([]byte,0)
  cmd := strings.Fields(line)
  instruction := cmd[0]
  //handle instructions with no arguments here
  switch instruction {
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
    //if instruction not found, read one argument then switch case with one-argument instructions
    default:
      if len(cmd) < 1 {
        //instructions missing arguments
        os.Exit(5)
      }
      dest := cmd[1]
      switch instruction {
        case "jp":
          output = append(output, jumpCall(dest, "jp")...)
        case "jpz":
          output = append(output, jumpCall(dest, "jpz")...)
        case "jpnz":
          output = append(output, jumpCall(dest, "jpnz")...)
        case "jpc":
          output = append(output, jumpCall(dest, "jpc")...)
        case "jpnc":
          output = append(output, jumpCall(dest, "jpnc")...)
        case "call":
          output = append(output, jumpCall(dest, "call")...)
        case "callz":
          output = append(output, jumpCall(dest, "callz")...)
        case "callnz":
          output = append(output, jumpCall(dest, "callnz")...)
        case "callc":
          output = append(output, jumpCall(dest, "callc")...)
        case "callnc":
          output = append(output, jumpCall(dest, "callnc")...)
        case "rst":
          newAddress := parseByte(dest)
          if !isValidRst(newAddress) {
            //reset vector is not valid
            os.Exit(6)
          }
          output = append(output, 0xc7 + newAddress)
        case "push":
          output = append(output, pushPop(dest, "push"))
        case "pop":
          output = append(output, pushPop(dest, "pop"))
        case "inc":
          output = append(output, incDec(dest, "inc"))
        case "dec":
          output = append(output, incDec(dest, "dec"))
        case "add":
          output = append(output, arithmetic(dest, "add"))
        case "adc":
          output = append(output, arithmetic(dest, "adc"))
        case "sub":
          output = append(output, arithmetic(dest, "sub"))
        case "sbc":
          output = append(output, arithmetic(dest, "sbc"))
        case "and":
          output = append(output, arithmetic(dest, "and"))
        case "xor":
          output = append(output, arithmetic(dest, "xor"))
        case "or":
          output = append(output, arithmetic(dest, "or"))
        case "cp":
          output = append(output, arithmetic(dest, "cp"))
        //instruction not found, read second argument the switch case with two-argument instructions
        default:
          if len(cmd) < 2 {
            os.Exit(5)
          }
          data := cmd[2]
          switch instruction {
            case "ld":
            output = append(output, load(dest, data)...)
            default:
              output = append(output , 0xff)
          }
      }
  }
  return output
}

