# gb-assembler
![demo](https://github.com/ayrtonm/gb-assembler/blob/master/demo.gif)
This is a basic assembler written in Go for creating binaries compatible with gameboy binaries. While it is still under development, it's capable of making [games](https://github.com/ayrtonm/gb-assembler/blob/master/demo.asm) like the one shown above.

## Building
Run `./compile.sh` (or equivalently `go build -o main src/*.go`).

## Usage
To assemble `input.asm` (and any files it links with `.include`) into a binary `output.gb`, run `./main input.asm output.gb`. A third argument can optionally be used to output a list of the opcodes written to the binary. This can be useful for debugging emulators before all opcodes are implemented. See [demo.ops](https://github.com/ayrtonm/gb-assembler/blob/master/demo.ops) for more details on the output format.

## Instruction documentation
[This page](http://www.pastraiser.com/cpu/gameboy/gameboy_opcodes.html) includes almost everything you need to know to start writing programs. The syntax used in this program is mostly copied verbatim from that page with the following notable exceptions:

- Registers are prefixed by `$`
- Values used as indices into memory are enclosed by `[ ]` instead of `( )`
- No commas between operands in instructions that take two arguments (e.g. `ld`, `set`, ...)
- Opcodes that store a result in the accumulator register (`$a`), don't explicitly include it as an operand (e.g. use `add $b` instead of `add $a $b`)

## Labels
To abstract away hardware details from code, labels to 16-bit addresses can be defined in the following ways:

- A line with a 16-bit number suffixed by `:` directs the assembler to start writing from that address
- A line with a single word suffixed by `:` defines a named label to the current address
- A line with `.var` followed by a name and `byte` or `word` defines a named label to an arbitrary location in work RAM (0xC000-0xDFFF). Using `word` ensures that two bytes can be stored at that address without overlapping with other `.var` labels.
- `.save` can be used in the same way as `.var` to define named labels to external RAM (0xA000-0xBFFF)

Labels can be used as arguments to jump and call instructions (without the `:`) anywhere inside their scope where they were defined (even before they're defined). `start` and `main` are included by default as a labels to 0x0100 and 0x0150, respectively, in the top level scope. The hardware/emulator starts execution from 0x0100, but most games jump straight to 0x0150 since the binary has [header info](bgb.bircd.org/pandocs.htm#thecartridgeheader) between these locations. By using named labels and the two default labels, programs can be written independent of the code's location in the binary making them far easier to modify.

**TODO** explain indentation-based label scopes used here and how it's implemented

## Other quirks
Multiple files can be linked together by adding `.include` followed by a filename to the main input file. Unlike the include preprocessor directive in C/C++, included files are parsed after the current file is finished (for now). This means that included files that only have code and named labels would be appended at the end of the binary. See `demo.asm` for other features I have yet to document well.
