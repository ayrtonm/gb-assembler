title:
dynarec test

start:
jp 0x0150

0x0150:
push $bc
jp 0x0250

0x0250:
pop $bc
jp 0x0150