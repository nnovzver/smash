#include <stdio.h>
#include <string.h>
#include "Simple.gen.h"

int main() {
  int err;
  int i;
  First c = {0, 0, 3, 4};
  First cc;
  First cc_etalon = {2, 4, 3, 4};
  unsigned char etalon[sizeof(First)] = {0x82, 0, 1, 0x80, 0, 0, 0x2, 0};
  unsigned char buf[sizeof(First)];

  // marshal test
  err = Marshal_First(&c, buf, sizeof(First));
  if (err == -1 ) {
    printf("FAIL! marshal\n");
    return 1;
  }
  if (memcmp(etalon, buf, sizeof(First)) != 0) {
    printf("FAIL! marshal check\n");
    for (i = 0; i < 8; ++i) printf("0x%hhX\n", buf[i]);
    return 1;
  }

  // unmarshal test
  err = Unmarshal_First(&cc, buf, sizeof(First));
  if (err == -1 ) {
    printf("FAIL! unmarshal\n");
    return 1;
  }
  if (memcmp(&cc, &cc_etalon, sizeof(First))) {
    printf("FAIL! unmarshal check\n");
    printf("cc.i = %d\n", cc.i);
    printf("cc.j = %d\n", cc.j);
    printf("cc.k = %d\n", cc.k);
    printf("cc.l = %d\n", cc.l);
    printf("cc_etalon.i = %d\n", cc_etalon.i);
    printf("cc_etalon.j = %d\n", cc_etalon.j);
    printf("cc_etalon.k = %d\n", cc_etalon.k);
    printf("cc_etalon.l = %d\n", cc_etalon.l);
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
