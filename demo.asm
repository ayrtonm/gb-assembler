//this demo draws a ball on the screen, checks which buttons are pressed and
//updates its position accordingly. when the ball hits an edge, it shows up on
//the other side. this is just a proof-of-concept since a a more realistic
//program would define some calling convention and use functions to not repeat
//code as much and loading the sprite data which occurs at the setup label would
//use the .data directive instead

.title
  demo game

//this label is always set to 0x0100 since this is where the program counter is
//initialized
start:
  jp setup

//this label is always set to 0x0150 since this is right after the cartridge
//header. games typically start by jumping to this point since there is very
//little unused space between 0x0100 and where the cartridge header starts
main:
  call wait
  call update_ball
  call check_right
  call check_left
  call check_down
  call check_up
  jp main

check_button:
  ld $hl 0xff00
  ld [hl] 0x20
  and [hl]
  ret

check_up:
  //check if up is pressed
  ld $a 0x04
  call check_button
  //return if it's not pressed
  retnz
  //if it's pressed, decrement ball position
  dec $c
  //if it didn't hit the top edge we're done
  retnz
  //wrap around if we hit the top edge
  ld $c 153
  ret

check_down:
  //check if down is pressed
  ld $a 0x08
  call check_button
  //return if it's not pressed
  retnz
  //if it's pressed, increment ball position
  inc $c
  //check if we hit the bottom edge
  ld $a $c
  sbc 153
  //if we didn't hit it, we're done
  retnz
  //wrap around if we hit the bottom edge
  ld $c 0x00
  ret

check_right:
  //check if right is pressed
  ld $a 0x01
  call check_button
  //return if it's not pressed
  retnz
  //if it's pressed, increment ball position
  inc $b
  //check if we hit the right edge
  ld $a $b
  sbc 169
  //if we didn't hit it, we're done
  retnz
  //wrap around if we hit the right edge
  ld $b 0x00
  ret

check_left:
  //check if left is pressed
  ld $a 0x02
  call check_button
  //return if it's not pressed
  retnz
  //if it's pressed, increment ball position
  dec $b
  //if we didn't the left edge we're done
  retnz
  //wrap around if we hit the left edge
  ld $b 0xa9
  ret

update_ball:
  ld $hl 0xfe00
  //set ball's y-coordinate
  ld $a $c
  ld [hl++] $a
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
  ld $a [0xff40]
  //ld $hl 0xff40
  //ld $a [hl]
  //used to clear 3rd bit of [0xff40]
  ld $b 0xfb
  and $b
  //used to set 3rd bit of [0xff40]
  //ld $b 0x04
  //or $b
  ld [0xff40] $a
  
  //initialize sprite coordinates
  //$b holds the x-coordinate and $c holds the y-coordinate
  ld $bc 0x1010
  jp main

wait:
//dec $hl does not modify any flags... meaning we have to do decrease them
//individually and do this silly double loop thing 
  ld $hl 0xffff
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
