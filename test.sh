go test -run ''
for i in src/*;
  do inputfile=${i#src/};
  go run assembler.go $i bin/${inputfile%.asm}.gb;
  hexdump -C bin/${inputfile%.asm}.gb > dumps/${inputfile%.asm}.hex;
done
