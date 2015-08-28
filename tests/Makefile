test: runtest
	./runtest

../smash: ../*.go
	cd ..; go build; cd -

runtest: test.c simple_proto.gen.c
	clang -o runtest -Wall test.c simple_proto.gen.c

simple_proto.gen.c: simple_proto.json ../smash
	../smash simple_proto.json

clean:
	rm -f ../smash runtest simple_proto.gen{c,h}
