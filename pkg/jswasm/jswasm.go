// Package jswasm creates a goja vm with webassembly support from wagon
package jswasm

import (
	"bytes"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/go-interpreter/wagon/exec"

	"github.com/go-interpreter/wagon/wasm"

	"github.com/dop251/goja"

	"github.com/dop251/goja_nodejs/eventloop"
)

const (
	WasmPage = 65536
)

func New() *JSWasmVM {
	jw := &JSWasmVM{
		js: NewJS(),
	}
	jw.js.Run(func(vm *goja.Runtime) {
		vm.Set("__mem", func(cfg map[string]int) []byte {
			// initial, ok := cfg["initial"]
			// if !ok {
			// 	initial = 1
			// }
			//need to look into this...large memory is very slow
			//was faster when was setting it directly to global?
			initial := 1
			jw.buf = make([]byte, initial*WasmPage)
			return jw.buf
		})
		vm.Set("__table", func(cfg map[string]int) {
			initial, ok := cfg["initial"]
			if !ok {
				initial = 0
			}
			jw.table = make([]uint32, initial)
		})
		vm.Set("read", func(filename string) []byte {
			data, err := ioutil.ReadFile(filepath.Join(jw.dir, filename))
			if err != nil {
				log.Fatal(err)
			}
			return data
		})
		vm.Set("__inst", func(args goja.FunctionCall) goja.Value {
			// Get Args Binary and Info
			b := args.Argument(0).Export().([]uint8)
			im := args.Argument(1).Export().(map[string]interface{})

			//Decode Module
			builder := wasm.NewBuilder().Decode(bytes.NewReader(b))
			builder.ResolveImports(jw.resolver(vm, builder.M, im))
			builder.Populate()
			if builder.Err != nil {
				log.Fatal(builder.Err)
			}
			m := builder.M
			wvm, err := exec.NewVM(m, exec.EnableAOT(false))
			if err != nil {
				log.Fatal(err)
			}
			exports := map[string]func(args ...uint64) goja.Value{}
			for k, v := range m.Export.Entries {
				//capture local for closure
				v := v
				k := k
				fn := m.GetFunction(int(v.Index))
				params := len(fn.Sig.ParamTypes)
				f := func(args ...uint64) goja.Value {
					if params < len(args) {
						args = args[:params]
					}
					val, err := wvm.ExecCode(int64(v.Index), args...)
					if err != nil {
						log.Fatal(err)
					}
					return vm.ToValue(val)
				}
				exports[k] = f
			}
			res := map[string]interface{}{
				"module":   nil,
				"instance": map[string]interface{}{"exports": exports},
			}
			return vm.ToValue(res)
		})
		script := `
	require('index.js')
	var WebAssembly={
		Memory:function(config){
			var buf=__mem(config)
			this.buffer=new ArrayBuffer(buf)
		},
		Table:function(desc){
			var t=__table(desc)
			var length=0
			this.get=function(arg1){
				console.log("table get",arg1)
				return 0
			}
			this.grow=function(){
				console.log("table grow")
			}
			this.set=function(){
				length++
				console.log("table set")
			}

		},
		instantiate:function(){
			var args=arguments
			return new Promise(function(res,rej){
				var val=__inst.apply(this,args)
				res(val)
			})
			
		}

	}
	`
		_, err := vm.RunString(script)
		if err != nil {
			log.Fatal(err)
		}
	})
	return jw
}

type JSWasmVM struct {
	js    *eventloop.EventLoop
	buf   []byte
	table []uint32
	ws    *exec.VM
	dir   string
}

func (jw *JSWasmVM) RunFile(file string) {
	jw.dir = filepath.Dir(file)
	js, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	jw.js.Run(func(vm *goja.Runtime) {
		_, err := vm.RunString(string(js))
		if err != nil {
			log.Fatal(err)
		}

	})

}
