test: runtest
	cd ..; go test; cd -
	./runtest

../smash: ../*.go
	cd ..; go build; cd -

runtest: main.o cpptest.o ctest.o simple_proto.gen.c.o simple_proto.gen.cpp.o
	clang++ -Wall -lQt5Core -o runtest main.o cpptest.o ctest.o \
		simple_proto.gen.c.o simple_proto.gen.cpp.o

simple_proto.gen.h: simple_proto.json ../smash
	../smash -h simple_proto.json

simple_proto.gen.c: simple_proto.json ../smash
	../smash -c simple_proto.json

simple_proto.gen.hpp: simple_proto.json ../smash
	../smash -H simple_proto.json

simple_proto.gen.cpp: simple_proto.json ../smash
	../smash -C simple_proto.json

main.o: main.cpp ctest.h
	clang++ -c -Wall -o main.o main.cpp

cpptest.o: cpptest.cpp simple_proto.gen.hpp
	clang++ -c -Wall -isystem /usr/include/qt/QtCore -isystem /usr/include/qt\
		-o cpptest.o cpptest.cpp

ctest.o: ctest.c simple_proto.gen.h
	clang -c -Wall -o ctest.o ctest.c

simple_proto.gen.c.o: simple_proto.gen.h simple_proto.gen.c
	clang -c -Wall -o simple_proto.gen.c.o simple_proto.gen.c

simple_proto.gen.cpp.o: simple_proto.gen.hpp simple_proto.gen.cpp
	clang++ -c -Wall -isystem /usr/include/qt/QtCore -isystem /usr/include/qt \
		-o simple_proto.gen.cpp.o simple_proto.gen.cpp

clean:
	rm -f ../smash runtest simple_proto.gen{c,h,cpp,hpp}
