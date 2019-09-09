#go test -run ''
for i in src_z80/*;
  do inputfile=${i#src_z80/};
  go run src/* $i bin_z80/${inputfile%.asm}.gb;
  hexdump -C bin_z80/${inputfile%.asm}.gb > dumps_z80/${inputfile%.asm}.hex;
done
