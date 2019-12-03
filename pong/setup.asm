//these are aliases to variables in work RAM
.var ball_py byte
.var ball_px byte
.var ball_vy byte
.var ball_vx byte
.var bar_px byte
.var bar_vx byte
.var opponent_px byte
.var opponent_vx byte
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

data_section:
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
  initial_oam:
    //bar
    0x00010090
    //bar_2
    0x00020090
    //bar_3
    0x20020090
    //opponent
    0x00010010
    //opponent_2
    0x00020010
    //opponent_3
    0x20020010

setup:
  //load ball tile data to the first tile pattern
  ld $hl data_section
  ld $de vram
  ld $b 48
  call move_data
  //load data into OAM table
  ld $de bar_y
  ld $b 24
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


  //this hides the points after a game over
  ld $a 0
  ld [point1_x] $a
  ld [point2_x] $a
  ld [point3_x] $a
  ld [point4_x] $a
  ld [point5_x] $a
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
  ld [wins] $a
  ld [losses] $a
reset_game:
    ld $a 80
    ld [ball_py] $a
    ld [ball_px] $a
    call random
    ld [ball_vy] $a
    call random
    ld [ball_vx] $a
    jp main
