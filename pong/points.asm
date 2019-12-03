lost_point:
  ld $hl wins
  push $hl
  ld $hl losses
  ld $b 160
  ld $e -8
  jp point
won_point:
  ld $hl losses
  push $hl
  ld $hl wins
  ld $b 8
  ld $e 8
  //implicitly jumps to point
point:
  //increase losses variable
  //c stores the number of losses
  ld $c [hl]
  inc [hl]
  //d stores the total number of points scored
  pop $hl
  ld $a [hl]
  add $c
  ld $d $a
  //draw point on the screen
  //hl stores the oam pointer
  ld $hl point1_x
  //b stores the x coordinate
  //ld $b 160
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
    //sub 8
    add $e
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
