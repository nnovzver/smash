test: runtest
	./runtest

../smash: ../*.go
	cd ..; go build; cd -

runtest: test.c Simple.gen.c
	clang -o runtest -Wall test.c Simple.gen.c

Simple.gen.c: simple_proto.json ../smash
	../smash simple_proto.json

clean:
	rm -f ../smash runtest Simple.gen{c,h}
