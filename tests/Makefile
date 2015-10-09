test: runtest
	cd ..; gofmt -d -s -w *.go; cd -
	cd ..; go test; cd -
	./runtest

../smash: ../*.go
	cd ..; go build; cd -

runtest: main.o testcpp.o testc.o simple_proto.gen.c.o simple_proto.gen.cpp.o
	clang++ -Wall -lQt5Core -o runtest main.o testcpp.o testc.o \
		simple_proto.gen.c.o simple_proto.gen.cpp.o

simple_proto.gen.h: simple_proto.json ../smash
	../smash -h simple_proto.json

simple_proto.gen.c: simple_proto.json ../smash
	../smash -c simple_proto.json

simple_proto.gen.hpp: simple_proto.json ../smash simple_proto.gen.h
	../smash -H simple_proto.json

simple_proto.gen.cpp: simple_proto.json ../smash
	../smash -C simple_proto.json

main.o: main.cpp testc.h
	clang++ -c -Wall -o main.o main.cpp

testcpp.o: testcpp.cpp simple_proto.gen.hpp
	clang++ -c -Wall -isystem /usr/include/qt/QtCore -isystem /usr/include/qt\
		-o testcpp.o testcpp.cpp

testc.o: testc.c simple_proto.gen.h
	clang -c -Wall -o testc.o testc.c

simple_proto.gen.c.o: simple_proto.gen.h simple_proto.gen.c
	clang -c -Wall -o simple_proto.gen.c.o simple_proto.gen.c

simple_proto.gen.cpp.o: simple_proto.gen.hpp simple_proto.gen.cpp
	clang++ -c -Wall -isystem /usr/include/qt/QtCore -isystem /usr/include/qt \
		-o simple_proto.gen.cpp.o simple_proto.gen.cpp

clean:
	rm -f *.o ../smash runtest simple_proto.gen.{c,h,cpp,hpp}
