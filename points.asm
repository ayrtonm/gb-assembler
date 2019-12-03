lost_point:
  //increase losses variable
  ld $a [losses]
  inc $a
  ld [losses] $a
  //draw point on the screen
  ld $hl point1_x
  ld $b 160
  ld $a [losses]
  ld $c $a
  dec $c
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
    dec $b
    dec $c
    jpnz move_pointer
  update_oam:
    ld $a $b
    ld [hl] $a
  jp reset_game

won_point:
  jp reset_game

