package dom

import (
	"encoding/json"
	"fmt"
	"reflect"
	"syscall/js"
)

func ConsoleLog(data ...any) {
	consoleArgs := processConsoleData(data...)
	js.Global().Get("console").Call("log", consoleArgs...)
}

func ConsoleError(data ...any) {
	consoleArgs := processConsoleData(data...)
	js.Global().Get("console").Call("log", consoleArgs...)
}

func ConsoleWarn(data ...any) {
	consoleArgs := processConsoleData(data...)
	js.Global().Get("console").Call("log", consoleArgs...)
}

func processConsoleData(data ...any) []any {
	var processedArgs []any

	for _, item := range data {
		if item == nil {
			processedArgs = append(processedArgs, "nil")
			continue
		}

		switch v := item.(type) {
		case string, bool, int, int32, int64, uint, uint32, uint64, float32, float64:
			processedArgs = append(processedArgs, v)
			continue

		case fmt.Stringer:
			processedArgs = append(processedArgs, v.String())
			continue

		case js.Value:
			if v.IsUndefined() || v.IsNull() {
				processedArgs = append(processedArgs, "undefined or null")
			} else {
				processedArgs = append(processedArgs, v)
			}
			continue
		}

		rVal := reflect.ValueOf(item)

		switch rVal.Kind() {
		case reflect.String:
			processedArgs = append(processedArgs, rVal.String())
		case reflect.Struct, reflect.Ptr, reflect.Map, reflect.Slice:
			b, err := json.Marshal(item)
			if err != nil {
				processedArgs = append(processedArgs, "Error marshalling to JSON:", err.Error())
			} else {
				processedArgs = append(processedArgs, string(b))
			}
		default:
			processedArgs = append(processedArgs, item)
		}
	}

	return processedArgs
}
