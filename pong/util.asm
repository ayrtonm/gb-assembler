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

long_wait:
  push $hl
  ld $hl 0xa0ff
  jp wait
short_wait:
  push $hl
  //the product of the values in $h and $l determine the wait time
  ld $hl 0x10ff
  //there's an implicit jp wait here
//dec $hl does not modify any flags... meaning we have to do decrease them
//individually and do this silly relative jump thing
wait:
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

//a slightly less shitty pseudo rng that uses the status of the lcd (vblank, hblank,...)
//to get two random bits. Leaves the result in $a. Since I'm using this to initialize
//the ball's velocity, I'm filtering out zeros and making the result signed so
//this returns -1, 1 or 2
.var rand_seed byte
random:
  //setup aux registers
  push $bc
  ld $b 0x00
  ld $c 8
  loop:
    //lowest 2 bits of [lcdstatus] depend on lcd state
    ld $a [lcdstatus]
    and 0x03
    rl $b
    xor $b
    ld $b $a
    dec $c
    jpnz loop
  ld $a [rand_seed]
  add $b
  ld [rand_seed] $a
  pop $bc
  and 0x03
  sub 1
  jpz random
  ret

//return true (set $a to 1) if $b < $a < $c
in_range:
  cp $b
  jpc false
  cp $c
  jpnc false
  true:
    ld $a 1
    ret
  false:
    ld $a 0
    ret
