lost_point:
  //increase losses variable
  ld $a [losses]
  inc $a
  ld [losses] $a
  //draw point on the screen
  //hl stores the oam pointer
  ld $hl point1_x
  //b stores the x coordinate
  ld $b 160
  //c stores the number of losses
  ld $a [losses]
  ld $c $a
  dec $c
  //d stores the total number of points scored
  ld $a [wins]
  add $c
  ld $d $a
  jpz update_oam
  //get pointer to the right place in the OAM
  move_pointer:
    ld $a $l
    add 4
    ld $l $a
    dec $d
    jpnz move_pointer
  move_x:
    ld $a $b
    sub 8
    ld $b $a
    dec $c
    jpnz move_x
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
  //increase wins variable
  ld $a [wins]
  inc $a
  ld [wins] $a
  //draw point on the screen
  //hl stores the oam pointer
  ld $hl point1_x
  //b stores the x coordinate
  ld $b 8
  //c stores the number of wins
  ld $a [wins]
  ld $c $a
  dec $c
  //d stores the total number of points scored
  ld $a [losses]
  add $c
  ld $d $a
  jpz update_oam
  //get pointer to the right place in the OAM
  move_pointer:
    ld $a $l
    add 4
    ld $l $a
    dec $d
    jpnz move_pointer
  move_x:
    ld $a $b
    add 8
    ld $b $a
    dec $c
    jpnz move_x
  update_oam:
    ld $a $b
    ld [hl] $a
    call long_wait
  //check for game over
  ld $a [wins]
  sub 3
  jpz setup
  jp reset_game
