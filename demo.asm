.include registers.asm
.include util.asm
//this demo is basically a single-player version of pong
//for now the game is only reset when the ball hits the bottom edge
//whenever it hits the other edges or the bar, the component of it's velocity
//perpendicular to the edge is negated. there's also a bar_vx variable defined
//though it's not currently used. I should fix the check_collision function
//and make a better pseudo rng to initialize the ball before implementing
//acceleration and friction with the bar_vx variable
//an include directive would also be nice to avoid having these huge files
.title
  pong

start:
  jp setup

main:
  call draw_ball
  call draw_bar
  call draw_opponent
  call move_ball
  call move_opponent
  call check_collision
  call check_keypad
  call wait
  jp main

//I should add some kind of scoping capability to the assembler
//I often have to make sure that label names don't collide leading to long names
//if I could use indentation to define some kind of "function" scope I could reuse
//label names
check_collision:
  left:
    ld $b 8
    ld $a [ball_px]
    cp $b
    jpnz right
    call negate_vx
    jp top
  right:
    ld $b 160
    cp $b
    jpnz top
    call negate_vx
  top:
    ld $b 24
    ld $a [ball_py]
    cp $b
    jpnz bar
    call negate_vy
    ret
  //double check this section
  bar:
    ld $b 137
    cp $b
    jpnz bottom
    ld $a [ball_px]
    ld $b $a
    ld $a [bar_px]
    cp $b
    jpc less_than
    sub $b
    jp greater_than
    less_than:
      ld $c $a
      ld $a $b
      sub $a
    greater_than:
      cp 12
      jpnc bottom
      call negate_vy
      ret
  bottom:
    ld $b 152
    cp $b
    retnz
    jp lost_point

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
  left:
    ld $a [keypad]
    and 0x02
    jpnz right
    ld $b 16
    ld $a [bar_px]
    cp $b
    jpz right
    sub 1
    ld [bar_px] $a
  right:
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

move_ball:
  horizontal:
    ld $a [ball_vx]
    cp 0x80
    jpnc left
    right:
      ld $b $a
      ld $a [ball_px]
      add $b
      ld [ball_px] $a
      jp vertical
    left:
      call negate
      ld $b $a
      ld $a [ball_px]
      sub $b
      ld [ball_px] $a
  vertical:
    ld $a [ball_vy]
    cp 0x80
    jpnc up
    down:
      ld $b $a
      ld $a [ball_py]
      add $b
      ld [ball_py] $a
      ret
    up:
      call negate
      ld $b $a
      ld $a [ball_py]
      sub $b
      ld [ball_py] $a
      ret

move_opponent:
  ld $a [ball_px]
  ld $b $a
  ld $a [opponent_px]
  sub $b
  jpnc move_left
  move_right:
    call random
    and 1
    ld $b $a
    ld $a [opponent_px]
    inc $a
    add $b
    ld [opponent_px] $a
    ret
  move_left:
    call random
    and 1
    ld $b $a
    ld $a [opponent_px]
    dec $a
    sub $b
    ld [opponent_px] $a
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

draw_opponent:
  ld $a [opponent_px]
  ld [opponent_x] $a
  ld $a [opponent_px]
  add 8
  ld [opponent_x2] $a
  ld $a [opponent_px]
  sub 8
  ld [opponent_x3] $a
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

  //load bar edge tile data to the third tile pattern
  ld $hl bar_edge_data
  ld $de vram
  ld $a $e
  add 32
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
  ld $a 144
  ld [bar_y] $a
  ld [bar_y2] $a
  ld [bar_y3] $a
  ld $a 1
  ld [bar_tile] $a
  ld $a 2
  ld [bar_tile2] $a
  ld [bar_tile3] $a
  ld $a 0x00
  ld [bar_attr] $a
  ld [bar_attr2] $a
  ld $a 0x20
  ld [bar_attr3] $a

  ld $a 16
  ld [opponent_y] $a
  ld [opponent_y2] $a
  ld [opponent_y3] $a
  ld $a 1
  ld [opponent_tile] $a
  ld $a 2
  ld [opponent_tile2] $a
  ld [opponent_tile3] $a
  ld $a 0x00
  ld [opponent_attr] $a
  ld [opponent_attr2] $a
  ld $a 0x20
  ld [opponent_attr3] $a

  ld $a 152
  ld [point1_y] $a
  ld [point2_y] $a
  ld [point3_y] $a
  ld [point4_y] $a
  ld [point5_y] $a

  //initialize variables in work ram
  ld $a 80
  ld [bar_px] $a
  ld [opponent_px] $a
  ld $a 0
  ld [bar_vx] $a
  ld $a 0
  ld [wins] $a
  ld [losses] $a
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

lost_point:
  //increase losses variable
  ld $a [losses]
  inc $a
  ld [losses] $a
  //draw point on the screen
  ld $hl point1_y
  ld $b 160
  ld $a [losses]
  dec $a
  jpz update_oam
  //get pointer to the right place in the OAM
  move_pointer:
    inc $hl
    inc $hl
    inc $hl
    inc $hl
    dec $b
    dec $b
    dec $b
    dec $b
    dec $b
    dec $b
    dec $b
    dec $a
    jpnz move_pointer
  update_oam:
    inc $hl
    ld $a $b
    ld [hl] $a
    jp reset_game

won_point:
  jp reset_game

//these are aliases to variables in work RAM
.var ball_py byte
.var ball_px byte
.var ball_vy byte
.var ball_vx byte
.var bar_px byte
//this is currently unused
.var bar_vx byte
.var opponent_px byte
.var wins byte
.var losses byte

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

.alias opponent_y 0xfe10
.alias opponent_x 0xfe11
.alias opponent_tile 0xfe12
.alias opponent_attr 0xfe13
.alias opponent_y2 0xfe14
.alias opponent_x2 0xfe15
.alias opponent_tile2 0xfe16
.alias opponent_attr2 0xfe17
.alias opponent_y3 0xfe18
.alias opponent_x3 0xfe19
.alias opponent_tile3 0xfe1a
.alias opponent_attr3 0xfe1b

.alias point1_y 0xfe1c
.alias point1_x 0xfe1d
.alias point1_tile 0xfe1e
.alias point1_attr 0xfe1f

.alias point2_y 0xfe20
.alias point2_x 0xfe21
.alias point2_tile 0xfe22
.alias point2_attr 0xfe23

.alias point3_y 0xfe24
.alias point3_x 0xfe25
.alias point3_tile 0xfe26
.alias point3_attr 0xfe27

.alias point4_y 0xfe28
.alias point4_x 0xfe29
.alias point4_tile 0xfe2a
.alias point4_attr 0xfe2b

.alias point5_y 0xfe2c
.alias point5_x 0xfe2d
.alias point5_tile 0xfe2e
.alias point5_attr 0xfe2f

ball_tile_data:
  0x0000
  0x003c
  0x1866
  0x3c4a
  0x3c42
  0x1866
  0x003c
  0x0000

bar_tile_data:
  0x0000
  0x00ff
  0xff00
  0xff00
  0xff00
  0xff00
  0x00ff
  0x0000

bar_edge_data:
  0x0000
  0x00fc
  0xf806
  0xfc0a
  0xfc02
  0xf806
  0x00fc
  0x0000
