//this demo draws a ball on the screen, checks which buttons are pressed and
//updates its position accordingly. when the ball hits an edge, it shows up on
//the other side. this is just a proof-of-concept since a a more realistic
//program would define some calling convention and use functions to not repeat
//code as much (ex: one function to check any button) and loading the sprite
//data which occurs at the setup label would use the .data directive instead

.title
  demo game

start:
  jp setup

0x0150:
main:
  call wait
  call wait
  call update_ball
  call check_right
  call check_left
  call check_down
  call check_up
  jp main

check_up:
  //check if up is pressed
  ld $hl 0xff00
  ld [hl] 0x20
  ld $a 0x04
  and [hl]
  //return if it's not pressed
  retnz
  //if it's pressed, decrement ball position
  ld $a $d
  sub $e
  ld $d $a
  retnz
  ld $d 0x99
  ret

check_down:
  //check if down is pressed
  ld $hl 0xff00
  ld [hl] 0x20
  ld $a 0x08
  and [hl]
  //return if it's not pressed
  retnz
  //if it's pressed, increment ball position
  ld $a $d
  add $e
  ld $d $a
  ld $h 0x99
  sbc $h
  retnz
  ld $d 0x00
  ret

check_right:
  ld $hl 0xff00
  ld [hl] 0x20
  ld $a 0x01
  and [hl]
  //return if it's not pressed
  retnz
  //if it's pressed, increment ball position
  ld $a $b
  add $c
  ld $b $a
  ld $h 0xa9
  sbc $h
  retnz
  ld $b 0x00
  ret

check_left:
  ld $hl 0xff00
  ld [hl] 0x20
  ld $a 0x02
  and [hl]
  //return if it's not pressed
  retnz
  //if it's pressed, increment ball position
  ld $a $b
  sub $c
  ld $b $a
  retnz
  ld $b 0xa9
  ret

update_ball:
  ld $hl 0xfe00
  //set ball's y-coordinate
  ld [hl] $d
  inc $hl
  //set ball's x-coordinate
  ld [hl] $b
  ret

setup:
  //load sprite tile data
  ld $hl 0x8002
  ld [hl] 0x3c
  inc $hl
  inc $hl
  ld [hl] 0x66
  inc $hl
  ld [hl] 0x18
  inc $hl
  ld [hl] 0x4a
  inc $hl
  ld [hl] 0x3c
  inc $hl
  ld [hl] 0x42
  inc $hl
  ld [hl] 0x3c
  inc $hl
  ld [hl] 0x66
  inc $hl
  ld [hl] 0x18
  inc $hl
  ld [hl] 0x3c
  //load sprite tile data
  ld $hl 0x8012
  ld [hl] 0x3c
  inc $hl
  inc $hl
  ld [hl] 0x66
  inc $hl
  ld [hl] 0x18
  inc $hl
  ld [hl] 0x4a
  inc $hl
  ld [hl] 0x3c
  inc $hl
  ld [hl] 0x42
  inc $hl
  ld [hl] 0x3c
  inc $hl
  ld [hl] 0x66
  inc $hl
  ld [hl] 0x18
  inc $hl
  ld [hl] 0x3c

  //load palette (obp0)
  ld $hl 0xff48
  ld [hl] 0x1b
  
  //disable bkg, win/enable sprites
  ld $hl 0xff40
  ld [hl] 0x86
  //select direction keys
  ld $hl 0xff00
  ld [hl] 0x20
  
  //set 8x8 sprites
  ld $hl 0xff40
  ld $a [hl]
  //used to clear 3rd bit of [0xff40]
  ld $b 0xfb
  and $b
  //used to set 3rd bit of [0xff40]
  //ld $b 0x04
  //or $b
  ld [hl] $a
  
  //initial sprite x coordinate
  ld $b 0x10
  //increment x coordinate by $c if right button is not pressed
  ld $c 0x01
  //initial sprite y coordinate
  ld $d 0x10
  //increment y coordinate by $e if B button is not pressed
  ld $e 0x01
  jp main

wait:
  ld $h 0xff
  ld $l 0xff
  wait_loop_1:
    dec $h
    nop
    jpnz wait_loop_1
    wait_loop_2:
      ld $h 0xff
      dec $l
      nop
      jpnz wait_loop_2
      ret
