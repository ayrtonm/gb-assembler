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
    ld $a $l
    add 4
    ld $l $a
    ld $a $b
    sub 8
    ld $b $a
    dec $c
    jpnz move_pointer
  update_oam:
    ld $a $b
    ld [hl] $a
    call long_wait
  //check for game over
  ld $a [losses]
  sub 3
  jpz setup
  jp reset_game

won_point:
  //increase losses variable
  ld $a [wins]
  inc $a
  ld [wins] $a
  //draw point on the screen
  ld $hl point1_x
  ld $b 8
  ld $a [wins]
  ld $c $a
  dec $c
  jpz update_oam
  //get pointer to the right place in the OAM
  move_pointer:
    ld $a $l
    add 4
    ld $l $a
    ld $a $b
    add 8
    ld $b $a
    dec $c
    jpnz move_pointer
  update_oam:
    ld $a $b
    ld [hl] $a
    call long_wait
  //check for game over
  ld $a [wins]
  sub 3
  jpz setup
  jp reset_game
