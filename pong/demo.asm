.include registers.asm
.include util.asm
.include points.asm
.include setup.asm
//this is a simple pong demo
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
    //check the ball and top bar's x coordinates
    ld $a [opponent_px]
    sub 12
    ld $b $a
    ld $a [opponent_px]
    add 16
    ld $c $a
    ld $a [ball_px]
    call in_range
    and 1
    //if the ball's x coordinate is not within reach of the bar
    //check the ball's y coordinate to see if we scored a point
    jpz check_won_point
    //otherwise check if the ball bounced off the bar
    ld $a [ball_py]
    ld $b 22
    cp $b
    jpc bounce
    jp bottom
    bounce:
      //the ball starts moving down after bouncing off the top bar
      call negate_vy
      //check if the ball and the bar are moving in the same direction
      ld $a [opponent_vx]
      ld $b $a
      ld $a [ball_vx]
      rlc $a
      xor $b
      and 1
      //if they are moving in opposite directions reverse the ball's direction
      callz negate_vx
      jp bottom
    check_won_point:
      ld $a [ball_py]
      ld $b 22
      cp $b
      callc won_point
  bottom:
    //check the ball and the bottom bar's x coordinates
    ld $a [bar_px]
    sub 12
    ld $b $a
    ld $a [bar_px]
    add 16
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
