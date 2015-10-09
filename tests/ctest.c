#include <stdio.h>
#include <string.h>
#include "simple_proto.gen.h"

int ctest() {
  int err;
  int i;
  First tomarshal = {0, 0, 3, 4};
  tomarshal.b[0] = 0xAA;
  tomarshal.b[1] = 0xBB;
  tomarshal.b[2] = 0xCC;
  tomarshal.b[3] = 0xDD;

  First tounmarshal;
  First tounmarshal_etalon = {2, 4, 3, 4};
  tounmarshal_etalon.b[0] = 0xAA;
  tounmarshal_etalon.b[1] = 0xBB;
  tounmarshal_etalon.b[2] = 0xCC;
  tounmarshal_etalon.b[3] = 0xDD;
  unsigned char buf_etalon[First__BUFSIZE] = {0x82, 0, 1, 0x80, 0, 0, 0x2, 0, 0xAA, 0xBB, 0xCC, 0xDD};
  unsigned char buf[First__BUFSIZE];


  // marshal test
  err = Marshal_First(&tomarshal, buf, First__BUFSIZE);
  if (err == -1 ) {
    printf("FAIL! marshal\n");
    return 1;
  }
  if (memcmp(buf_etalon, buf, First__BUFSIZE) != 0) {
    printf("FAIL! marshal check\n");
    printf("buf buf_etalon\n");
    for (i = 0; i < First__BUFSIZE; ++i)
      printf("0x%hhX 0x%hhX\n", buf[i], buf_etalon[i]);
    return 1;
  }

  // unmarshal test
  err = Unmarshal_First(&tounmarshal, buf, First__BUFSIZE);
  if (err == -1 ) {
    printf("FAIL! unmarshal\n");
    return 1;
  }
  if (memcmp(&tounmarshal, &tounmarshal_etalon, sizeof(First))) {
    printf("FAIL! unmarshal check\n");
    printf("tounmarshal tounmarshal_etalon\n");
    for (i = 0; i < sizeof(First); ++i)
      printf("0x%hhX 0x%hhX\n", ((uint8_t*)&tounmarshal)[i], ((uint8_t*)&tounmarshal_etalon)[i]);

    return 1;
  }

  // test test =)
  if (Is_First(buf, First__BUFSIZE) != 1) {
    printf("FAIL! IS test\n");
    return 1;
  }

  printf("OK - C TEST\n");
  return 0;
}
