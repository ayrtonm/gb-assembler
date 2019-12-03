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
  ld [opponent_tile] $a
  ld $a 2
  ld [bar_tile2] $a
  ld [bar_tile3] $a
  ld [opponent_tile2] $a
  ld [opponent_tile3] $a
  ld $a 0x00
  ld [bar_attr] $a
  ld [bar_attr2] $a
  ld [opponent_attr] $a
  ld [opponent_attr2] $a
  ld $a 0x20
  ld [bar_attr3] $a
  ld [opponent_attr3] $a

  ld $a 16
  ld [opponent_y] $a
  ld [opponent_y2] $a
  ld [opponent_y3] $a

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


