#include <iostream>
#include <cstring>
#include <cstdio>
#include <QByteArray>

#include "simple_proto.gen.hpp"

int cpptest() {
  Codograms::First tomarshal;
  tomarshal.m.i = 0;
  tomarshal.m.j = 0;
  tomarshal.m.k = 3;
  tomarshal.m.l = 4;
  tomarshal.m.b[0] = 0xAA;
  tomarshal.m.b[1] = 0xBB;
  tomarshal.m.b[2] = 0xCC;
  tomarshal.m.b[3] = 0xDD;

  Codograms::First tounmarshal;
  Codograms::First tounmarshal_etalon;
  tounmarshal_etalon.m.i = 2;
  tounmarshal_etalon.m.j = 4;
  tounmarshal_etalon.m.k = 3;
  tounmarshal_etalon.m.l = 4;
  tounmarshal_etalon.m.b[0] = 0xAA;
  tounmarshal_etalon.m.b[1] = 0xBB;
  tounmarshal_etalon.m.b[2] = 0xCC;
  tounmarshal_etalon.m.b[3] = 0xDD;
  unsigned char buf_etalon[First__BUFSIZE] = {0x82, 0, 1, 0x80, 0, 0, 0x2, 0, 0xAA, 0xBB, 0xCC, 0xDD};

  // marshal test
  if (!tomarshal.marshal()) {
    printf("FAIL! marshal\n");
    return 1;
  }
  if (memcmp(buf_etalon, tomarshal.buf.constData(), tomarshal.buf.size()) != 0) {
    printf("FAIL! marshal check\n");
    printf("buf buf_etalon\n");
    for (int i = 0; i < tomarshal.buf.size(); ++i)
      printf("0x%hhX 0x%hhX\n", tomarshal.buf.constData()[i], buf_etalon[i]);
    return 1;
  }

  // unmarshal test
  tounmarshal.buf = tomarshal.buf;
  if (!tounmarshal.unmarshal()) {
    printf("FAIL! unmarshal\n");
    return 1;
  }
  if (memcmp(&tounmarshal.m, &tounmarshal_etalon.m, tounmarshal_etalon.msize()) != 0) {
    printf("FAIL! unmarshal check\n");
    printf("tounmarshal tounmarshal_etalon\n");
    for (int i = 0; i < tounmarshal_etalon.msize(); ++i)
      printf("0x%hhX 0x%hhX\n", ((uint8_t*)&tounmarshal.m)[i], ((uint8_t*)&tounmarshal_etalon.m)[i]);

    return 1;
  }

  // Is test
  if (!tomarshal.checkBuf()) {
    printf("FAIL! IS test\n");
    return 1;
  }

  printf("OK - CPP TEST\n");
  return 0;
}
