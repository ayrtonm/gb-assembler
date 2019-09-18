.title
  demo game

start:
  jp 0x0150

0x0150:
  jp main

0x2000:
main:
  add $c
  add $h
  ld $b $f
  jp 336
