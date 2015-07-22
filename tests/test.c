#include <stdio.h>
#include <string.h>
#include "Simple.gen.h"

int main() {
  int err;
  int i;
  First tomarshal = {0, 0, 3, 4};
  First tounmarshal;
  First tounmarshal_etalon = {2, 4, 3, 4};
  unsigned char etalon_buf[sizeof(First)] = {0x82, 0, 1, 0x80, 0, 0, 0x2, 0};
  unsigned char buf[sizeof(First)];

  // marshal test
  err = Marshal_First(&tomarshal, buf, sizeof(First));
  if (err == -1 ) {
    printf("FAIL! marshal\n");
    return 1;
  }
  if (memcmp(etalon_buf, buf, sizeof(First)) != 0) {
    printf("FAIL! marshal check\n");
    for (i = 0; i < 8; ++i) printf("0x%hhX\n", buf[i]);
    return 1;
  }

  // unmarshal test
  err = Unmarshal_First(&tounmarshal, buf, sizeof(First));
  if (err == -1 ) {
    printf("FAIL! unmarshal\n");
    return 1;
  }
  if (memcmp(&tounmarshal, &tounmarshal_etalon, sizeof(First))) {
    printf("FAIL! unmarshal check\n");
    printf("tounmarshal.i = %d\n", tounmarshal.i);
    printf("tounmarshal.j = %d\n", tounmarshal.j);
    printf("tounmarshal.k = %d\n", tounmarshal.k);
    printf("tounmarshal.l = %d\n", tounmarshal.l);
    printf("tounmarshal_etalon.i = %d\n", tounmarshal_etalon.i);
    printf("tounmarshal_etalon.j = %d\n", tounmarshal_etalon.j);
    printf("tounmarshal_etalon.k = %d\n", tounmarshal_etalon.k);
    printf("tounmarshal_etalon.l = %d\n", tounmarshal_etalon.l);
    return 1;
  }

  // test test =)
  if (isFirst(buf, sizeof(First)) != 1) {
    printf("FAIL! test\n");
    return 1;
  }

  printf("OK\n");
  return 0;
}
