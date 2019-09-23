# Example

get wagoja

```console
go get github.com/sampaioletti/wagoja
```

if you take a file main.c

```C
#include <stdio.h>

int main(){
  puts("Main");
}
```

and build it with the emscripten compiler https://emscripten.org/docs/getting_started/downloads.html

```console
emcc main.c -s EXPORTED_FUNCTIONS="['_main']"

```

it will generate a.out.js and a.out.wasm

you can generate a wat (human readable wasm) with the wasm2wat tool https://github.com/WebAssembly/wabt

```console
wasm2wat a.out.wasm -o a.out.wat
```

run it from ../cmd/wagoja with

```console
go run main.go ../../example/a.out.js
```

and with a little luck you will see it print "Main" to your screen
