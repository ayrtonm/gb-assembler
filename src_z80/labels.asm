title:
testing labels

start:
jp 0x0150

0x0150:
call 0x2000

testb:
call 0x0150

0x2000:
call testb
