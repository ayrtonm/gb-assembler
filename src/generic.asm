title:
test rom

start:
jp 0x0150

0x0150:
inc $bc
inc $sp
inc $af
ld $b, $b
ld $b, $a
rst 0x08
add $b
add [hl]
add $c
sub $c
sbc [hl]
xor $d
cp $e
and $b
