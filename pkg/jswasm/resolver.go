package jswasm

import (
	"log"

	"github.com/dop251/goja"
	"github.com/go-interpreter/wagon/wasm"
)

func (jw *JSWasmVM) resolver(vm *goja.Runtime, m *wasm.Module, i map[string]interface{}) wasm.ResolveFunc {
	wasm.NewModule()
	builder := wasm.NewBuilder().NewModule()
	builder.SetLinearMemory(jw.buf).SetTableMemory(jw.table)
	builder.M.Elements = m.Elements
	builder.M.Export.Entries = map[string]wasm.ExportEntry{}
	builder.M.Types = m.Types

	env := i["env"].(map[string]interface{})
	for _, v := range m.Import.Entries {
		v := v
		switch t := v.Type.(type) {
		case wasm.FuncImport:
			cb := env[v.FieldName].(func(goja.FunctionCall) goja.Value)
			sig := m.Types.Entries[t.Type]
			wrap := func(args ...uint64) uint64 {
				params := goja.FunctionCall{
					This: goja.Undefined(),
				}
				for _, arg := range args {
					params.Arguments = append(params.Arguments, vm.ToValue(arg))
				}
				ret := cb(params)
				return uint64(ret.ToInteger())
			}
			builder.AddFunction(v.FieldName, true, sig, wrap)

		case wasm.GlobalVarImport:
			builder.AddGlobalVal(v.FieldName, true, env[v.FieldName])
		case wasm.MemoryImport:
			builder.AddMemory(v.FieldName, true, t.Type.Limits.Initial, t.Type.Limits.Maximum)
		case wasm.TableImport:
			builder.AddTable(v.FieldName, true, t.Type.Limits.Initial, t.Type.Limits.Maximum)
		default:
			log.Fatalf("unknown import type %+v", v.Type)
		}
	}
	err := builder.Populate().Err
	if err != nil {
		log.Fatal(err)
	}
	return func(name string) (*wasm.Module, error) {
		return builder.M, nil
	}
}
