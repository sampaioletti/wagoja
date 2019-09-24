//compile with emcc main.c EXPORTED_FUNCTIONS="['_mult','main']"
#include <stdio.h>
int hello(int x) {
  puts("Hello From WaSM");
  return x;
}
int main(){
  puts("Main");
}
