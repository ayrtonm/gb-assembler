go test -run ''
for i in src/*;
  do go run assembler.go $i bin/${i#src/}.gb;
  hexdump -C $i > dumps/${i#src/}.hex;
done
