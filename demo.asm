.include registers.asm
.include util.asm
.include points.asm
.include setup.asm
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
  call short_wait
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
    ret
  right:
    ld $b 160
    ld $a [ball_px]
    cp $b
    jpnz top
    retnz
    call negate_vx
  top:
    //check the ball and bar's x coordinates
    ld $a [opponent_px]
    sub 12
    ld $b $a
    ld $a [opponent_px]
    add 20
    ld $c $a
    ld $a [ball_px]
    call in_range
    and 1
    jpz check_won_point
    ld $a [ball_py]
    ld $b 22
    cp $b
    jpc bounce
    jp bottom
    bounce:
      call negate_vy
      ld $a [opponent_vx]
      ld $b $a
      ld $a [ball_vx]
      rlc $a
      xor $b
      and 1
      callz negate_vx
      jp bottom
    check_won_point:
      ld $a [ball_py]
      ld $b 22
      cp $b
      callc won_point
  bottom:
    ld $a [bar_px]
    sub 12
    ld $b $a
    ld $a [bar_px]
    add 20
    ld $c $a
    ld $a [ball_px]
    call in_range
    and 1
    jpz check_lost_point
    ld $a [ball_py]
    ld $b 138
    cp $b
    jpnc bounce
    ret
    bounce:
      call negate_vy
      ld $a [bar_vx]
      ld $b $a
      ld $a [ball_vx]
      rlc $a
      xor $b
      and 1
      callz negate_vx
      ret
    check_lost_point:
      ld $a [ball_py]
      ld $b 138
      cp $b
      callnc lost_point
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
    ld $a 0
    ld [bar_vx] $a
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
    ld $a 1
    ld [bar_vx] $a
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
    ld $a 1
    ld [opponent_vx] $a
    ret
  move_left:
    call random
    and 1
    ld $b $a
    ld $a [opponent_px]
    dec $a
    sub $b
    ld [opponent_px] $a
    ld $a 0
    ld [opponent_vx] $a
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
