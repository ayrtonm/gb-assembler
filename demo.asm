//this demo draws a ball on the screen, checks which buttons are pressed and
//updates its position accordingly. when the ball hits an edge, it shows up on
//the other side. this is just a proof-of-concept since a a more realistic
//program would define some calling convention and use functions to not repeat
//code as much

.title
  pong

start:
  jp setup

main:
  call draw_ball
  call draw_bar
  call update_ball_position
  call check_collision
  call check_keypad
  call wait
  jp main

//I should add some kind of scoping capability to the assembler
//I often have to make sure that label names don't collide leading to long names
//if I could use indentation to define some kind of "function" scope I could reuse
//label names
check_collision:
  check_left_wall:
    ld $b 8
    ld $a [ball_px]
    cp $b
    jpnz check_right_wall
    call negate_vx
    jp check_top_wall
  check_right_wall:
    ld $b 160
    cp $b
    jpnz check_top_wall
    call negate_vx
  check_top_wall:
    ld $b 16
    ld $a [ball_py]
    cp $b
    jpnz check_bar_collision
    call negate_vy
    ret
  //double check this section
  check_bar_collision:
    ld $b 145
    cp $b
    jpnz check_bottom_wall
    ld $a [ball_px]
    ld $b $a
    ld $a [bar_px]
    cp $b
    jpc l1
    sub $b
    jp l2
    l1:
    ld $c $a
    ld $a $b
    sub $a
    l2:
    cp 12
    jpnc check_bottom_wall
    call negate_vy
    ret
  check_bottom_wall:
    ld $b 152
    cp $b
    retnz
    jp reset_game
    //this is unreachable if we call setup
    call negate_vy
    ret

negate_vx:
  ld $a [ball_vx]
  call negate
  ld [ball_vx] $a
  ret

negate_vy:
  ld $a [ball_vy]
  call negate
  ld [ball_vy] $a
  ret

check_keypad:
  check_left:
    ld $a [keypad]
    and 0x02
    jpnz check_right
    ld $b 16
    ld $a [bar_px]
    cp $b
    jpz check_right
    sub 1
    ld [bar_px] $a
  check_right:
    ld $a [keypad]
    and 0x01
    retnz
    ld $b 152
    ld $a [bar_px]
    cp $b
    retz
    add 1
    ld [bar_px] $a
    ret

update_ball_position:
  update_ball_px:
    ld $a [ball_vx]
    cp 0x80
    jpnc move_ball_left
    move_ball_right:
      ld $b $a
      ld $a [ball_px]
      add $b
      ld [ball_px] $a
      jp update_ball_py
    move_ball_left:
      call negate
      ld $b $a
      ld $a [ball_px]
      sub $b
      ld [ball_px] $a
  update_ball_py:
    ld $a [ball_vy]
    cp 0x80
    jpnc move_ball_up
    move_ball_down:
      ld $b $a
      ld $a [ball_py]
      add $b
      ld [ball_py] $a
      ret
    move_ball_up:
      call negate
      ld $b $a
      ld $a [ball_py]
      sub $b
      ld [ball_py] $a
      ret

draw_ball:
  ld $hl ball_py
  ld $de ball_y
  ld $b 2
  call move_data
  ret

draw_bar:
  ld $a [bar_px]
  ld [bar_x] $a
  ld $a [bar_px]
  add 8
  ld [bar_x2] $a
  ld $a [bar_px]
  sub 8
  ld [bar_x3] $a
  ret

setup:
  //load ball tile data to the first tile pattern
  ld $hl ball_tile_data
  ld $de vram
  ld $b 16
  call move_data

  //load bar tile data to the second tile pattern
  ld $hl bar_tile_data
  ld $de vram
  ld $a $e
  add 16
  ld $e $a
  ld $b 16
  call move_data

  //load palette (obp0)
  ld $a 0x1b
  ld [palette] $a
  //disable bkg, win/enable sprites, set 8x8 sprites
  ld $a 0x82
  ld [lcdcontrol] $a
  //select direction keys
  ld $a 0x20
  ld [keypad] $a

  //initialize oam sprite table
  ld $a 152
  ld [bar_y] $a
  ld $a 152
  ld [bar_y2] $a
  ld $a 152
  ld [bar_y3] $a
  ld $a 1
  ld [bar_tile] $a
  ld $a 1
  ld [bar_tile2] $a
  ld $a 1
  ld [bar_tile3] $a
  ld $a 0x00
  ld [bar_attr] $a
  ld $a 0x20
  ld [bar_attr2] $a
  ld $a 0x20
  ld [bar_attr3] $a

  //initialize variables in work ram
  ld $a 80
  ld [bar_px] $a
  ld $a 0
  ld [bar_vx] $a
  reset_game:
    ld $a 80
    ld [ball_py] $a
    ld $a 80
    ld [ball_px] $a
    call random
    ld [ball_vy] $a
    call random
    ld [ball_vx] $a
    jp main

//these are aliases to variables in work RAM
.var ball_py byte
.var ball_px byte
.var ball_vy byte
.var ball_vx byte
.var bar_px byte
.var bar_vx byte

.alias vram 0x8000
.alias keypad 0xff00
.alias lcdcontrol 0xff40
.alias lcdstatus 0xff41
.alias palette 0xff48
.alias oam 0xfe00
//these are aliases to the coordinates of sprites 1 and 2
//functions should update the position and velocity variables shown above
//then the drawing functions copy those variables to these locations
//although this is slightly more inefficient than directly updating the
//following addresses, it adds more flexibility since the ball and bar aren't
//so closely coupled with specific sprites
.alias ball_y 0xfe00
.alias ball_x 0xfe01
.alias bar_y 0xfe04
.alias bar_x 0xfe05
.alias bar_tile 0xfe06
.alias bar_attr 0xfe07
.alias bar_y2 0xfe08
.alias bar_x2 0xfe09
.alias bar_tile2 0xfe0a
.alias bar_attr2 0xfe0b
.alias bar_y3 0xfe0c
.alias bar_x3 0xfe0d
.alias bar_tile3 0xfe0e
.alias bar_attr3 0xfe0f

.data
ball_tile_data:
  0x003c0000
  0x3c4a1866
  0x18663c42
  0x0000003c

bar_tile_data:
  0xff00ff00
  0xff00ff00
  0xff00ff00
  0xff00ff00

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
