//move $b bytes from [hl] to [de]
move_data:
  ld $a [hl++]
  push $hl
  push $de
  pop $hl
  ld [hl++] $a
  push $hl
  pop $de
  pop $hl
  dec $b
  jpnz move_data
  ret

//move 1 byte from [hl] to [de]
move_byte:
  push $hl
  push $de
  ld $d [hl]
  pop $hl
  ld [hl] $d
  push $hl
  pop $de
  pop $hl
  ret

//dec $hl does not modify any flags... meaning we have to do decrease them
//individually and do this silly relative jump thing
wait:
  push $hl
  //the product of the values in $h and $l determine the wait time
  ld $hl 0x10ff
  //decrease $h until it's zero
  dec $l
  jrnz -3
  //reload $l and decrease $h by 1
  //go back to decreasing $l
  dec $h
  ld $l 0xff
  jrnz -8
  //return when $h is zero
  pop $hl
  ret

//two's complement of $a
negate:
  cpl
  add 1
  ret

//a really shitty pseudo rng that uses the status of the lcd (vblank, hblank,...)
//to get two random bits. Leaves the result in $a. Since I'm using this to initialize
//the ball's velocity, I'm filtering out zeros and making the result signed so
//this returns -1, 1 or 2. since this depends on the lcd's status, repeatedly calling
//it will tend to give sequences of the same value
random:
  ld $a [lcdstatus]
  and 0x03
  sub 1
  jpz random
  ret
