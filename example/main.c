//compile with emcc mult.c -Os -s WASM=1 -s SIDE_MODULE=1 -o mult.wasm -s EXPORTED_FUNCTIONS="['_mult']"
#include <stdio.h>
int hello(int x) {
  puts("Hello From WaSM");
  return x;
}
int main(){
  puts("Main");
}