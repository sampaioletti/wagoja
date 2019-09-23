# Wagoja

TLDR;

```console
go get github.com/sampaioletti/wagoja
cd cmd/wagoja
go run main.go ./path/to/emscripten/js/file.js
```

## What is it

wagoja is a dumpster fire mashup of github.com/go-interpreter/wagon and github.com/dop251/goja (hence the name WAGon gOJA)

Both are excellent projects, so no shade on them (: this was just an experiment

While trying to make an emscripten support module for wagon, I decided to try this to see if it would help explain how some of the emscripten internals work

At the time it sounded like an easier way to learn how emscripten worked...in hindsight...tough to say.

## What can it do

Not much has been tested at the moment

see [Example](./example/)

## Why would this be helpful

- Existing environments that support emscripten have limited platform support, most claim they are planning on it, but nothing on the radar.
- Wagon is a good wasm runtime and supports go platforms, but currently doesn't have a libc implemenation
- Goja is a great but currently not maintained javscript runtime in Go, I chose it over the more popular Otto, becasue I'm familiar with it and last i checked the performance was slightly better (should be tested)

## Is this going to turn into a project

Not likely, unless there is enough interest. Currently the plan is to use this as a basis and start implementing the emscripten api (thats what is generated in a.out.js by emscripten) inside wagon in go.

The downside to that is it will take a lot of maintenance

- The emscripten project is only supporting js environments at the moment.
- Other libc implementations are all written in C so would require making wagon use CGO.

As WASM matures, maybe another solution will present itself.

Turning this into an actual project wouldn't take a lot of work, but would take a lot of testing and I question the performance and worthwileness....prove me wrong..
